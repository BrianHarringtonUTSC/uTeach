package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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

// PostModel handles getting and creating posts.
type PostModel struct {
	Base
}

// NewPostModel returns a new post model.
func NewPostModel(db *sqlx.DB) *PostModel {
	return &PostModel{Base{db}}
}

var postsSqlizer = squirrel.
	Select(`posts.id AS post_id,
			posts.title AS post_title,
			posts.content,
			posts.created_at,
			posts.is_pinned,
			posts.is_visible,
			count(post_votes.post_id),
			topics.id AS topic_id,
			topics.name AS topic_name,
			topics.title AS topic_title,
			topics.description,
			users.id AS user_id,
			users.email,
			users.name AS user_name,
			users.is_admin`).
	From("posts").
	Join("topics ON topics.id=posts.topic_id").
	Join("users ON users.id=posts.creator_user_id").
	LeftJoin("post_votes ON post_votes.post_id=posts.id").
	GroupBy("posts.id").
	OrderBy("count(post_votes.post_id) DESC")

func (tm *PostModel) findAll(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) ([]*Post, error) {
	posts := []*Post{}

	query, args, err := sqlizer.ToSql()
	if err != nil {
		return posts, err
	}

	rows, err := tm.Query(tx, query, args...)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		post := new(Post)
		topic := new(Topic)
		creator := new(User)

		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.IsPinned, &post.IsVisible, &post.Score,
			&topic.ID, &topic.Name, &topic.Title, &topic.Description,
			&creator.ID, &creator.Email, &creator.Name, &creator.IsAdmin)
		if err != nil {
			return posts, err
		}

		post.Topic = topic
		post.Creator = creator
		posts = append(posts, post)
	}

	return posts, err
}

func (tm *PostModel) findOne(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) (*Post, error) {
	posts, err := tm.findAll(tx, sqlizer)
	if err != nil {
		return nil, err
	}
	if len(posts) != 1 {
		return nil, fmt.Errorf("Expected: 1, got: %d.", len(posts))
	}
	return posts[0], err
}

// GetPostByID gets a post by the id.
func (tm *PostModel) GetPostByID(tx *sqlx.Tx, id int64) (*Post, error) {
	return tm.findOne(tx, postsSqlizer.Where(squirrel.Eq{"posts.id": id}))
}

// GetPostByID gets a post by the id and topic.
func (tm *PostModel) GetPostByIDAndTopic(tx *sqlx.Tx, id int64, topic *Topic) (*Post, error) {
	return tm.findOne(tx, postsSqlizer.Where(squirrel.Eq{"posts.id": id, "topics.id": topic.ID}))
}

// GetPostsByTopicAndIsPinned gets all posts by topic and whether they are pinned or not pinned.
func (tm *PostModel) GetPostsByTopicAndIsPinned(tx *sqlx.Tx, topic *Topic, isPinned bool) ([]*Post, error) {
	posts, err := tm.findAll(tx, postsSqlizer.Where(squirrel.Eq{"posts.topic_id": topic.ID, "posts.is_pinned": isPinned}))
	if err == sql.ErrNoRows {
		return []*Post{}, nil
	}
	return posts, err
}

// GetPostsByUser gets all posts by the user.
func (tm *PostModel) GetPostsByUser(tx *sqlx.Tx, user *User) ([]*Post, error) {
	return tm.findAll(tx, postsSqlizer.Where(squirrel.Eq{"posts.creator_user_id": user.ID}))
}

// GetPostIdsUpvotedByUser gets the ids of all posts upvoted by the user. It returns a map which can be used to
// check if a post was upvoted by a user in constant time.
// TODO: this method may need to be made more precise. For example, finding all upvoted posts for a topic, etc.
func (tm *PostModel) GetPostIdsUpvotedByUser(tx *sqlx.Tx, user *User) (map[int64]bool, error) {
	rows, err := tm.Query(tx, "SELECT post_id FROM post_votes WHERE user_id=?", user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postIDs := map[int64]bool{}
	var postID int64
	for rows.Next() {
		rows.Scan(&postID)
		postIDs[postID] = true
	}
	return postIDs, err
}

// AddPost adds a new post.
func (tm *PostModel) AddPost(tx *sqlx.Tx, title, content string, topic *Topic, creator *User) (*Post, error) {
	if title == "" || content == "" {
		return nil, InputError{"Empty title or body not allowed"}
	}

	query := "INSERT INTO posts(title, content, topic_id, creator_user_id) VALUES(?, ?, ?, ?)"
	result, err := tm.Exec(tx, query, title, content, topic.ID, creator.ID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return tm.GetPostByID(tx, id)
}

// AddPostVoteForUser adds a vote for the post for the user.
func (tm *PostModel) AddPostVoteForUser(tx *sqlx.Tx, post *Post, user *User) error {
	_, err := tm.Exec(tx, "INSERT INTO post_votes(user_id, post_id) VALUES(?, ?)", user.ID, post.ID)
	return err
}

// RemoveTheadVoteForUser removes a vote for the post for the user.
func (tm *PostModel) RemoveTheadVoteForUser(tx *sqlx.Tx, post *Post, user *User) error {
	_, err := tm.Exec(tx, "DELETE FROM post_votes where user_id=? AND post_id=?", user.ID, post.ID)
	return err
}

// HidePost hides the post.
func (tm *PostModel) HidePost(tx *sqlx.Tx, post *Post) error {
	_, err := tm.Exec(tx, "UPDATE posts SET is_visible=? WHERE id=?", false, post.ID)
	return err
}

// UnhidePost unhides the post.
func (tm *PostModel) UnhidePost(tx *sqlx.Tx, post *Post) error {
	_, err := tm.Exec(tx, "UPDATE posts SET is_visible=? WHERE id=?", true, post.ID)
	return err
}

// PinPost pins a post.
func (tm *PostModel) PinPost(tx *sqlx.Tx, post *Post) error {
	_, err := tm.Exec(tx, "UPDATE posts SET is_pinned=? WHERE id=?", true, post.ID)
	return err
}

// UnpinPost unpins a post.
func (tm *PostModel) UnpinPost(tx *sqlx.Tx, post *Post) error {
	_, err := tm.Exec(tx, "UPDATE posts SET is_pinned=? WHERE id=?", false, post.ID)
	return err
}

// GetPostsByTag gets all posts with tag.
func (tm *PostModel) GetPostsByTag(tx *sqlx.Tx, tag *Tag) ([]*Post, error) {
	posts, err := tm.findAll(tx,
		postsSqlizer.Join("post_tags ON post_tags.post_id=posts.id").Where(squirrel.Eq{"post_tags.tag_id": tag.ID}))
	return posts, err
}
