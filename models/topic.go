package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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

	topics := make([]*Topic, 0)
	err = tm.sel(tx, &topics, query, args...)
	return topics, err
}

// FindOne gets the topic filtered by wheres.
func (tm *TopicModel) FindOne(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) (*Topic, error) {
	topics, err := tm.Find(tx, wheres...)
	if err != nil {
		return nil, err
	}

	switch len(topics) {
	case 0:
		return nil, sql.ErrNoRows
	case 1:
		return topics[0], err
	default:
		return nil, fmt.Errorf("topic: Expected: 1, got: %d.", len(topics))
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
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return tm.FindOne(tx, squirrel.Eq{"topics.id": id})
}
