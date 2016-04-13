package models

import (
	"strings"

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

// GetAllTopics gets all topics.
func (tm *TopicModel) GetAllTopics(tx *sqlx.Tx) ([]*Topic, error) {
	topics := []*Topic{}
	err := tm.Select(tx, &topics, "SELECT * FROM topics")
	return topics, err
}

// GetTopicByID gets a topic by id.
func (tm *TopicModel) GetTopicByID(tx *sqlx.Tx, id int64) (*Topic, error) {
	topic := new(Topic)
	err := tm.Get(tx, topic, "SELECT * FROM topics WHERE id=?", id)
	return topic, err
}

// GetTopicByName gets a topic by name.
func (tm *TopicModel) GetTopicByName(tx *sqlx.Tx, name string) (*Topic, error) {
	name = strings.ToLower(name)

	topic := new(Topic)
	err := tm.Get(tx, topic, "SELECT * FROM topics WHERE name=?", name)
	return topic, err
}

// AddTopic adds a new topic.
func (tm *TopicModel) AddTopic(tx *sqlx.Tx, name, title, description string) (*Topic, error) {
	if title == "" || description == "" || !singleWordAlphaNumRegex.MatchString(name) {
		return nil, InputError{"Invalid name and/or title."}
	}

	name = strings.ToLower(name)

	query := "INSERT INTO topics(name, title, description) VALUES(?, ?, ?)"
	result, err := tm.Exec(tx, query, name, title, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return tm.GetTopicByID(tx, id)
}
