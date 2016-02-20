// Package models provides models for the uTeach app.
package models

import (
	"fmt"
)

// User is a user in the uTeach system.
type User struct {
	Username string
}

// Subject represents a subject, the base database object.
type Subject struct {
	Name  string
	Title string
}

// Thread is a post inside of a topic.
type Thread struct {
	ID                int64
	Title             string
	Content           string
	Score             int
	SubjectName       string `db:"subject_name"`
	CreatedByUsername string `db:"created_by_username"`
}

// URL returns the unique URL for a thread.
func (t *Thread) URL() string {
	return fmt.Sprintf("/s/%s/%d", t.SubjectName, t.ID)
}
