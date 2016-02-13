package main

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
	SubjectName string
}

// Thread is a post inside of a topic.
type Thread struct {
	ID                int
	Title             string
	Content           string
	Score             int
	SubjectName       string
	TopicName         string
	CreatedByUsername string
}
