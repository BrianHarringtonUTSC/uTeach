package main

type Subject struct {
	ID   int64
	Name string
}

type Topic struct {
	ID        int64
	Name      string
	SubjectId int64
}

type Thread struct {
	ID      int64
	Title   string
	Content string
	TopicID int64
	// PostedByUserID int64
}
