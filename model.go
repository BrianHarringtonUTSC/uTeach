package main

type User struct {
	UTORid string
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
	ID             int64
	Title          string
	Content        string
	SubjectName    string
	TopicName      string
	PostedByUTORid string
}
