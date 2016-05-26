package models

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

// IsValid returns true if the tag is valid else false.
func (t *Tag) IsValid() bool {
	return singleWordAlphaNumRegex.MatchString(t.Name)
}

// TagModel handles getting and creating tags.
type TagModel struct {
	Base
}

// NewTagModel returns a new tag model.
func NewTagModel(db *sqlx.DB) *TagModel {
	return &TagModel{Base{db}}
}

var (
	// ErrInvalidTag is returned when adding or updating an invalid tag
	ErrInvalidTag = InputError{"Invalid name"}

	tagsBuilder = squirrel.
			Select("tags.id, tags.name, topics.id, topics.name, topics.title").
			From("tags").
			Join("topics ON topics.id=tags.topic_id").
			OrderBy("tags.name")
)

// Find gets all tags filtered by wheres.
func (tm *TagModel) Find(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) ([]*Tag, error) {
	rows, err := tm.queryWhere(tx, tagsBuilder, wheres...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		tag := new(Tag)
		topic := new(Topic)
		err = rows.Scan(&tag.ID, &tag.Name, &topic.ID, &topic.Name, &topic.Title)
		if err != nil {
			return nil, errors.Wrap(err, "scan error")
		}
		tag.Topic = topic
		tags = append(tags, tag)
	}
	return tags, nil
}

// FindOne gets the user filtered by wheres.
func (tm *TagModel) FindOne(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) (*Tag, error) {
	tags, err := tm.Find(tx, wheres...)
	if err != nil {
		return nil, errors.Wrap(err, "find error")
	}

	switch len(tags) {
	case 0:
		return nil, sql.ErrNoRows
	case 1:
		return tags[0], nil
	default:
		return nil, errors.Errorf("expected 1, got %d", len(tags))
	}
}

// Add adds a new tag.
func (tm *TagModel) Add(tx *sqlx.Tx, tag *Tag) error {
	if !tag.IsValid() {
		return ErrInvalidTag
	}

	tag.Name = strings.ToLower(tag.Name)
	result, err := tm.exec(tx, "INSERT INTO tags(name, topic_id) VALUES(?, ?)", tag.Name, tag.Topic.ID)
	if err != nil {
		return errors.Wrap(err, "exec error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "last inserted id error")
	}

	t, err := tm.FindOne(tx, squirrel.Eq{"tags.id": id})
	if err != nil {
		return errors.Wrap(err, "find one error")
	}

	*tag = *t
	return nil
}

// AddPostTag adds a tag for the post.
func (tm *TagModel) AddPostTag(tx *sqlx.Tx, post *Post, tag *Tag) error {
	_, err := tm.exec(tx, "INSERT INTO post_tags(post_id, tag_id, topic_id) VALUES(?, ?, ?)",
		post.ID, tag.ID, post.Topic.ID)
	return errors.Wrap(err, "exec error")
}
