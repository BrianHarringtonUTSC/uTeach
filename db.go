package main

// TODO: replace with real db
var subjects = []Subject{Subject{1, "Python"}, Subject{2, "Java"}}

var topics = []Topic{Topic{1, "For loops", 1}, Topic{1, "While loops", 1}}

func GetSubjects() []Subject {
	return subjects
}

func GetTopics(subjectID int64) []Topic {
	return topics
}
