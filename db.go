package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (err error) {
	DB, err = sql.Open("sqlite3", "./uteach.db")
	if err != nil {
		return
	}

	// ensure all required tables exist
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS user(
			username TEXT PRIMARY KEY
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS subject(
			name TEXT PRIMARY KEY,
			title TEXT
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS topic(
			name TEXT,
			title TEXT,
			subjectName TEXT,
			PRIMARY KEY(name, subjectName)
			FOREIGN KEY(subjectName) REFERENCES subject(name)
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS thread(
			id INTEGER PRIMARY KEY,
			title TEXT,
			content TEXT,
			score INTEGER,
			subjectName TEXT,
			topicName TEXT,
			postedByUsername TEXT,
			FOREIGN KEY(subjectName, topicName) REFERENCES topic(subjectName, name)
			FOREIGN KEY(postedByUsername) REFERENCES user(username)
		)`)
	if err != nil {
		return
	}

	return
}

func GetUser(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow("SELECT * FROM user WHERE username=?", username).Scan(&user.Username)
	return user, err
}

func GetSubjects() (subjects []*Subject, err error) {
	rows, err := DB.Query("SELECT * FROM subject")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		subject := &Subject{}
		rows.Scan(&subject.Name, &subject.Title)
		subjects = append(subjects, subject)
	}
	return
}

func GetTopics(subjectName string) (topics []*Topic, err error) {
	rows, err := DB.Query("SELECT * FROM topic WHERE subjectName=?", subjectName)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		topic := &Topic{}
		rows.Scan(&topic.Name, &topic.Title, &topic.SubjectName)
		topics = append(topics, topic)
	}
	return
}

func GetThreads(subjectName string, topicName string) (threads []*Thread, err error) {
	rows, err := DB.Query("SELECT * FROM thread WHERE subjectName=? AND topicName=?", subjectName, topicName)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		thread := &Thread{}
		rows.Scan(&thread.ID, &thread.Title, &thread.Content, &thread.Score, &thread.SubjectName, &thread.TopicName, &thread.PostedByUsername)
		threads = append(threads, thread)
	}
	return
}

func GetThread(threadID int) (*Thread, error) {
	thread := &Thread{}
	err := DB.QueryRow("SELECT * FROM thread WHERE id=?", threadID).Scan(&thread.ID, &thread.Title, &thread.Content, &thread.Score, &thread.SubjectName, &thread.TopicName, &thread.PostedByUsername)
	return thread, err // &Thread{1, "Some random tutorial on for loops", "Here's a random link to said random tutorial on for loops", "python", "for_loops", "umair"}
}
