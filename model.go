package main

type User struct {
	Username string
}

type Subject struct {
	Name  string
	Title string
}

type Topic struct {
	Name        string
	Title       string
	SubjectName string
}

type Thread struct {
	ID                int
	Title             string
	Content           string
	Score             int
	SubjectName       string
	TopicName         string
	CreatedByUsername string
}
