package main

// TODO: replace with real db
var subjects = []*Subject{&Subject{1, "Python"}, &Subject{2, "Java"}}

var topics = []*Topic{&Topic{1, "For loops", 1}, &Topic{2, "While loops", 1}}

var threads = []*Thread{&Thread{1, "Some random tutorial on for loops", "Here's a random link to said random tutorial on for loops", 1}}

func GetSubjects() []*Subject {
	return subjects
}

func GetTopics(subjectID int64) []*Topic {
	return topics
}

func GetThreads(topicID int64) []*Thread {
	return threads
}

func GetThread(threadID int64) *Thread {
	return threads[0]
}
