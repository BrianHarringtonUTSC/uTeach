package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	return sanitizeString(p.Content)
}

// PostModel handles getting and creating posts.
type PostModel struct {
	Base
}

// NewPostModel returns a new post model.
func NewPostModel(db *sqlx.DB) *PostModel {
	return &PostModel{Base{db}}
}

var postsBuilder = squirrel.
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
		msg := fmt.Sprintf("expected 1, got %d", len(posts))
		return nil, errors.New(msg)
	}
}

// GetVotedPostIds gets the ids of upvoted posts filtered by wheres. It returns a map that acts as a set (all values
// are true) which can be used for quick lookup.
func (pm *PostModel) GetVotedPostIds(tx *sqlx.Tx, where squirrel.Sqlizer) (map[int64]bool, error) {
	rows, err := pm.queryWhere(tx, squirrel.Select("post_id FROM post_votes"), where)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()

	postIDs := map[int64]bool{}
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

// AddPost adds a new post.
func (pm *PostModel) AddPost(tx *sqlx.Tx, title, content string, topic *Topic, creator *User) (*Post, error) {
	if title == "" || sanitizeString(content) == "" {
		return nil, InputError{"Empty title or body not allowed"}
	}

	query := "INSERT INTO posts(title, content, topic_id, creator_user_id) VALUES(?, ?, ?, ?)"
	result, err := pm.exec(tx, query, title, content, topic.ID, creator.ID)
	if err != nil {
		return nil, errors.Wrap(err, "exec error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "last inserted id error")
	}
	post, err := pm.FindOne(tx, squirrel.Eq{"posts.id": id})
	return post, errors.Wrap(err, "find one error")
}

// AddPostVoteForUser adds a vote for the post for the user.
func (pm *PostModel) AddPostVoteForUser(tx *sqlx.Tx, post *Post, user *User) error {
	_, err := pm.exec(tx, "INSERT INTO post_votes(user_id, post_id) VALUES(?, ?)", user.ID, post.ID)
	return errors.Wrap(err, "exec error")
}

// RemovePostVoteForUser removes a vote for the post for the user.
func (pm *PostModel) RemovePostVoteForUser(tx *sqlx.Tx, post *Post, user *User) error {
	_, err := pm.exec(tx, "DELETE FROM post_votes where user_id=? AND post_id=?", user.ID, post.ID)
	return errors.Wrap(err, "exec error")
}

// HidePost hides the post.
func (pm *PostModel) HidePost(tx *sqlx.Tx, post *Post) error {
	_, err := pm.exec(tx, "UPDATE posts SET is_visible=? WHERE id=?", false, post.ID)
	return errors.Wrap(err, "exec error")
}

// UnhidePost unhides the post.
func (pm *PostModel) UnhidePost(tx *sqlx.Tx, post *Post) error {
	_, err := pm.exec(tx, "UPDATE posts SET is_visible=? WHERE id=?", true, post.ID)
	return errors.Wrap(err, "exec error")
}

// PinPost pins a post.
func (pm *PostModel) PinPost(tx *sqlx.Tx, post *Post) error {
	_, err := pm.exec(tx, "UPDATE posts SET is_pinned=? WHERE id=?", true, post.ID)
	return errors.Wrap(err, "exec error")
}

// UnpinPost unpins a post.
func (pm *PostModel) UnpinPost(tx *sqlx.Tx, post *Post) error {
	_, err := pm.exec(tx, "UPDATE posts SET is_pinned=? WHERE id=?", false, post.ID)
	return errors.Wrap(err, "exec error")
}
