package models

import (
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

var tagsSqlizer = squirrel.
	Select("tags.id, tags.name, topics.id AS topic_id, topics.name AS topic_name, topics.title").
	From("tags").
	Join("topics ON topics.id=tags.topic_id")

func (tm *TagModel) findAll(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) ([]*Tag, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tm.Query(tx, query, args...)
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

func (tm *TagModel) findOne(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) (*Tag, error) {
	tags, err := tm.findAll(tx, sqlizer)
	if err != nil {
		return nil, err
	}

	if len(tags) != 1 {
		return nil, fmt.Errorf("Expected: 1, got: %d.", len(tags))
	}

	return tags[0], err
}

// GetTagByID gets a tag by the id.
func (tm *TagModel) GetTagByID(tx *sqlx.Tx, id int64) (*Tag, error) {
	return tm.findOne(tx, tagsSqlizer.Where(squirrel.Eq{"tags.id": id}))
}

// GetTagByNameAndTopic gets a tag by the name and topic.
func (tm *TagModel) GetTagByNameAndTopic(tx *sqlx.Tx, name string, topic *Topic) (*Tag, error) {
	return tm.findOne(tx, tagsSqlizer.Where(squirrel.Eq{"tags.name": name, "tags.topic_id": topic.ID}))
}

// GetTagsByTopic gets all tags by the topic.
func (tm *TagModel) GetTagsByTopic(tx *sqlx.Tx, topic *Topic) ([]*Tag, error) {
	return tm.findAll(tx, tagsSqlizer.Where(squirrel.Eq{"topic_id": topic.ID}))
}

// AddTag adds a new tag for the topic.
func (tm *TagModel) AddTag(tx *sqlx.Tx, name string, topic *Topic) (*Tag, error) {
	if !singleWordAlphaNumRegex.MatchString(name) {
		return nil, InputError{"Invalid name."}
	}

	name = strings.ToLower(name)
	result, err := tm.Exec(tx, "INSERT INTO tags(name, topic_id) VALUES(?, ?)", name, topic.ID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return tm.GetTagByID(tx, id)
}

// AddPostTag adds a tag for the post.
func (tm *TagModel) AddPostTag(tx *sqlx.Tx, post *Post, tag *Tag) error {
	_, err := tm.Exec(tx, "INSERT INTO post_tags(post_id, tag_id, topic_id) VALUES(?, ?, ?)",
		post.ID, tag.ID, post.Topic.ID)
	return err
}
