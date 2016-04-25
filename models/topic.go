package models

import (
	"database/sql"
	"fmt"
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
func (s *Topic) URL() string {
	return "/topics/" + s.Name
}

// NewPostURL returns the URL of the page to create a new post under the topic.
func (s *Topic) NewPostURL() string {
	return s.URL() + "/new"
}

// TagsURL returns the URL of the page listing the tags under the topic.
func (s *Topic) TagsURL() string {
	return s.URL() + "/tags"
}

// NewTagURL returns the URL of the page to create a new tag under the topic.
func (s *Topic) NewTagURL() string {
	return s.TagsURL() + "/new"
}

// TopicModel handles getting and creating topics.
type TopicModel struct {
	Base
}

// NewTopicModel returns a new topic model.
func NewTopicModel(db *sqlx.DB) *TopicModel {
	return &TopicModel{Base{db}}
}

var topicsBuilder = squirrel.Select("* FROM topics")

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
		msg := fmt.Sprintf("expected 1, got %d", len(topics))
		return nil, errors.New(msg)
	}
}

// AddTopic adds a new topic.
func (tm *TopicModel) AddTopic(tx *sqlx.Tx, name, title, description string) (*Topic, error) {
	if title == "" || description == "" || !singleWordAlphaNumRegex.MatchString(name) {
		return nil, InputError{"Invalid name and/or title."}
	}

	name = strings.ToLower(name)

	query := "INSERT INTO topics(name, title, description) VALUES(?, ?, ?)"
	result, err := tm.exec(tx, query, name, title, description)
	if err != nil {
		return nil, errors.Wrap(err, "exec error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "last inserted id error")
	}

	topic, err := tm.FindOne(tx, squirrel.Eq{"topics.id": id})
	return topic, errors.Wrap(err, "topic error")
}
