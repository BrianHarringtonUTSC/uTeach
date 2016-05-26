package models

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Topic represents a topic in the app.
type Topic struct {
	ID          int64
	Name        string
	Title       string
	Description string
}

// URL returns the unique URL for a topic.
func (t *Topic) URL() string {
	return "/topics/" + t.Name
}

// NewPostURL returns the URL of the page to create a new post under the topic.
func (t *Topic) NewPostURL() string {
	return t.URL() + "/new"
}

// IsValid returns true if the topic is valid else false.
func (t *Topic) IsValid() bool {
	return t.Title != "" && t.Description != "" && singleWordAlphaNumRegex.MatchString(t.Name)
}

// TagsURL returns the URL of the page listing the tags under the topic.
func (t *Topic) TagsURL() string {
	return t.URL() + "/tags"
}

// NewTagURL returns the URL of the page to create a new tag under the topic.
func (t *Topic) NewTagURL() string {
	return t.TagsURL() + "/new"
}

// TopicModel handles getting and creating topics.
type TopicModel struct {
	Base
}

// NewTopicModel returns a new topic model.
func NewTopicModel(db *sqlx.DB) *TopicModel {
	return &TopicModel{Base{db}}
}

var (
	// ErrInvalidTopic is returned when adding or updating an invalid topic
	ErrInvalidTopic = InputError{"Cannot have empty name and/or title"}

	topicsBuilder = squirrel.Select("* FROM topics")
)

// Find gets all topics filtered by wheres.
func (tm *TopicModel) Find(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) ([]*Topic, error) {
	selectBuilder := tm.addWheresToBuilder(topicsBuilder, wheres...)
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	var topics []*Topic
	err = tm.sel(tx, &topics, query, args...)
	return topics, errors.Wrap(err, "select error")
}

// FindOne gets the topic filtered by wheres.
func (tm *TopicModel) FindOne(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) (*Topic, error) {
	topics, err := tm.Find(tx, wheres...)
	if err != nil {
		return nil, errors.Wrap(err, "find error")
	}

	switch len(topics) {
	case 0:
		return nil, sql.ErrNoRows
	case 1:
		return topics[0], nil
	default:
		return nil, errors.Errorf("expected 1, got %d", len(topics))
	}
}

// Add adds a new topic.
func (tm *TopicModel) Add(tx *sqlx.Tx, topic *Topic) error {
	if !topic.IsValid() {
		return ErrInvalidTopic
	}

	topic.Name = strings.ToLower(topic.Name)

	query := "INSERT INTO topics(name, title, description) VALUES(?, ?, ?)"
	result, err := tm.exec(tx, query, topic.Name, topic.Title, topic.Description)
	if err != nil {
		return errors.Wrap(err, "exec error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "last inserted id error")
	}

	t, err := tm.FindOne(tx, squirrel.Eq{"topics.id": id})
	if err != nil {
		return errors.Wrap(err, "find one error")
	}

	*topic = *t
	return nil
}
