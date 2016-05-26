package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
)

// Post represents a post in the app.
type Post struct {
	ID        int64
	Title     string
	Content   string
	CreatedAt time.Time
	IsPinned  bool
	IsVisible bool
	Score     int
	Topic     *Topic
	Creator   *User
}

// URL returns the unique URL for a post.
func (p *Post) URL() string {
	return p.Topic.URL() + fmt.Sprintf("/posts/%d", p.ID)
}

// SanitizedContent returns the post's content with markdown converted to HTML and sanitized.
func (p *Post) SanitizedContent() string {
	unsafe := blackfriday.MarkdownBasic([]byte(p.Content))
	safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	trimmed := strings.TrimSpace(string(safe))
	return trimmed
}

// IsValid returns true if the post is valid else false.
func (p *Post) IsValid() bool {
	return p.Title != "" && p.SanitizedContent() != ""
}

// PostModel handles getting and creating posts.
type PostModel struct {
	Base
}

// NewPostModel returns a new post model.
func NewPostModel(db *sqlx.DB) *PostModel {
	return &PostModel{Base{db}}
}

var (
	// ErrInvalidPost is returned when adding or updating an invalid post
	ErrInvalidPost = InputError{"Invalid post id or empty title or empty body"}

	postsBuilder = squirrel.
			Select(`posts.id, posts.title, posts.content, posts.created_at, posts.is_pinned, posts.is_visible,
			count(post_votes.post_id),
			topics.id, topics.name, topics.title, topics.description,
			users.id, users.email, users.name, users.is_admin`).
			From("posts").
			Join("topics ON topics.id=posts.topic_id").
			Join("users ON users.id=posts.creator_user_id").
			LeftJoin("post_votes ON post_votes.post_id=posts.id").
			LeftJoin("post_tags ON post_tags.post_id=posts.id").
			GroupBy("posts.id, post_tags.tag_id").
			OrderBy("count(post_votes.post_id) DESC, posts.created_at DESC").
			Distinct()
)

// Find gets all posts filtered by wheres.
func (pm *PostModel) Find(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) ([]*Post, error) {
	rows, err := pm.queryWhere(tx, postsBuilder, wheres...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := new(Post)
		topic := new(Topic)
		creator := new(User)

		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.IsPinned, &post.IsVisible, &post.Score,
			&topic.ID, &topic.Name, &topic.Title, &topic.Description,
			&creator.ID, &creator.Email, &creator.Name, &creator.IsAdmin)
		if err != nil {
			return nil, errors.Wrap(err, "scan error")
		}

		post.Topic = topic
		post.Creator = creator
		posts = append(posts, post)
	}

	return posts, nil
}

// FindOne gets the post filtered by wheres.
func (pm *PostModel) FindOne(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) (*Post, error) {
	posts, err := pm.Find(tx, wheres...)
	if err != nil {
		return nil, errors.Wrap(err, "find error")
	}

	switch len(posts) {
	case 0:
		return nil, sql.ErrNoRows
	case 1:
		return posts[0], nil
	default:
		return nil, errors.Errorf("expected 1, got %d", len(posts))
	}
}

// Add adds a new post.
func (pm *PostModel) Add(tx *sqlx.Tx, post *Post) error {
	if !post.IsValid() || post.ID > 0 {
		return ErrInvalidPost
	}

	result, err := pm.exec(tx, "INSERT INTO posts(title, content, topic_id, creator_user_id) VALUES(?, ?, ?, ?)",
		post.Title, post.Content, post.Topic.ID, post.Creator.ID)
	if err != nil {
		return errors.Wrap(err, "exec error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "last inserted id error")
	}

	p, err := pm.FindOne(tx, squirrel.Eq{"posts.id": id})
	if err != nil {
		return errors.Wrap(err, "find one error")
	}

	*post = *p
	return nil
}

// Update updates a post.
func (pm *PostModel) Update(tx *sqlx.Tx, post *Post) error {
	if post.ID < 1 || !post.IsValid() {
		return ErrInvalidPost
	}

	_, err := pm.exec(tx, "UPDATE posts SET title=?, content=?, is_pinned=?, is_visible=? WHERE id=?",
		post.Title, post.Content, post.IsPinned, post.IsVisible, post.ID)

	if err != nil {
		return errors.Wrap(err, "exec error")
	}

	p, err := pm.FindOne(tx, squirrel.Eq{"posts.id": post.ID})
	if err != nil {
		return errors.Wrap(err, "find one error")
	}

	*post = *p
	return nil
}

// GetVotedPostIds gets the ids of upvoted posts filtered by wheres. It returns a map that acts as a set (all values
// are true) which can be used for quick lookup.
func (pm *PostModel) GetVotedPostIds(tx *sqlx.Tx, where squirrel.Sqlizer) (map[int64]bool, error) {
	rows, err := pm.queryWhere(tx, squirrel.Select("post_id FROM post_votes"), where)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()

	postIDs := make(map[int64]bool)
	var postID int64
	for rows.Next() {
		err = rows.Scan(&postID)
		if err != nil {
			return nil, errors.Wrap(err, "scan error")
		}
		postIDs[postID] = true
	}
	return postIDs, nil
}

// UpdatePostVoteForUser addsa vote for the user for the post if voted is true else it removes the vote.
func (pm *PostModel) UpdatePostVoteForUser(tx *sqlx.Tx, post *Post, user *User, voted bool) error {
	var err error
	if voted {
		_, err = pm.exec(tx, "INSERT INTO post_votes(user_id, post_id) VALUES(?, ?)", user.ID, post.ID)
	} else {
		_, err = pm.exec(tx, "DELETE FROM post_votes where user_id=? AND post_id=?", user.ID, post.ID)
	}
	return errors.Wrap(err, "exec error")
}
