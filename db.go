package main

// TODO: replace with real db
var subjects = []*Subject{&Subject{"python", "Python"}, &Subject{"java", "Java"}}

var topics = []*Topic{&Topic{"for_loops", "For loops", "python"}, &Topic{"while_loops", "While loops", "python"}}

var threads = []*Thread{&Thread{1, "Some random tutorial on for loops", "Here's a random link to said random tutorial on for loops", "python", "for_loops", "umair"}}

func GetUser(utorid string) (*User, bool) {
	return &User{utorid}, true
}

func GetSubjects() []*Subject {
	return subjects
}

func GetTopics(subjectName string) []*Topic {
	return topics
}

func GetThreads(subjectName string, topicName string) []*Thread {
	return threads
}

func GetThread(subjectName string, topicName string, threadID int) *Thread {
	return threads[0]
}
