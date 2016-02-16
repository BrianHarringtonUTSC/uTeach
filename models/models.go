// Package models provides models for the uTeach app.
package models

// User is a user in the uTeach system.
type User struct {
	Username string
}

// Subject represents a subject, the base database object.
type Subject struct {
	Name  string
	Title string
}

// Topic is a sub category within subject.
type Topic struct {
	Name        string
	Title       string
	SubjectName string `db:"subject_name"`
}

// Thread is a post inside of a topic.
type Thread struct {
	ID                int `db:"rowid"`
	Title             string
	Content           string
	Score             int
	SubjectName       string `db:"subject_name"`
	TopicName         string `db:"topic_name"`
	CreatedByUsername string `db:"created_by_username"`
}
