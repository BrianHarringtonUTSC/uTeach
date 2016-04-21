package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// Tag represents a tag in the app.
type Tag struct {
	ID    int64
	Name  string
	Topic *Topic
}

// URL returns the unique URL for a topic.
func (t *Tag) URL() string {
	return t.Topic.TagsURL() + "/" + t.Name
}

// TagModel handles getting and creating tags.
type TagModel struct {
	Base
}

// NewTagModel returns a new tag model.
func NewTagModel(db *sqlx.DB) *TagModel {
	return &TagModel{Base{db}}
}

var tagsBuilder = squirrel.
	Select("tags.id, tags.name, topics.id AS topic_id, topics.name AS topic_name, topics.title").
	From("tags").
	Join("topics ON topics.id=tags.topic_id").
	OrderBy("tags.name")

// Find gets all tags filtered by wheres.
func (tm *TagModel) Find(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) ([]*Tag, error) {
	rows, err := tm.queryWhere(tx, tagsBuilder, wheres...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		tag := new(Tag)
		topic := new(Topic)
		err := rows.Scan(&tag.ID, &tag.Name, &topic.ID, &topic.Name, &topic.Title)
		if err != nil {
			return nil, err
		}
		tag.Topic = topic
		tags = append(tags, tag)
	}
	return tags, err
}

// FindOne gets the user filtered by wheres.
func (tm *TagModel) FindOne(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) (*Tag, error) {
	tags, err := tm.Find(tx, wheres...)
	if err != nil {
		return nil, err
	}

	switch len(tags) {
	case 0:
		return nil, sql.ErrNoRows
	case 1:
		return tags[0], err
	default:
		return nil, fmt.Errorf("Expected: 1, got: %d.", len(tags))
	}
}

// AddTag adds a new tag for the topic.
func (tm *TagModel) AddTag(tx *sqlx.Tx, name string, topic *Topic) (*Tag, error) {
	if !singleWordAlphaNumRegex.MatchString(name) {
		return nil, InputError{"Invalid name."}
	}

	name = strings.ToLower(name)
	result, err := tm.exec(tx, "INSERT INTO tags(name, topic_id) VALUES(?, ?)", name, topic.ID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return tm.FindOne(tx, squirrel.Eq{"tags.id": id})
}

// AddPostTag adds a tag for the post.
func (tm *TagModel) AddPostTag(tx *sqlx.Tx, post *Post, tag *Tag) error {
	_, err := tm.exec(tx, "INSERT INTO post_tags(post_id, tag_id, topic_id) VALUES(?, ?, ?)",
		post.ID, tag.ID, post.Topic.ID)
	return err
}
