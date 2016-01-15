package main

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
	ID          int64
	Name        string
	Title       string
	Content     string
	SubjectName string
	TopicName   string
	// PostedByUserID int64
}
