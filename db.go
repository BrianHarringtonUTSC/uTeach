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
		CREATE TABLE IF NOT EXISTS users(
			username TEXT PRIMARY KEY
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS subjects(
			name TEXT PRIMARY KEY,
			title TEXT NOT NULL
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS topics(
			name TEXT NOT NULL,
			title TEXT NOT NULL,
			subject_name TEXT NOT NULL,
			PRIMARY KEY(name, subject_name)
			FOREIGN KEY(subject_name) REFERENCES subjects(name)
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS threads(
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			score INTEGER NOT NULL,
			subject_name TEXT NOT NULL,
			topic_name TEXT NOT NULL,
			created_by_username TEXT NOT NULL,
			FOREIGN KEY(subject_name, topic_name) REFERENCES topics(subject_name, name)
			FOREIGN KEY(created_by_username) REFERENCES users(username)
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS upvotes(
			username TEXT NOT NULL
			thread_id INTEGER NOT NULL,
			FOREIGN KEY(username) REFERENCES users(username)
			FOREIGN KEY(thread_id) REFERENCES threads(id)
		)`)
	if err != nil {
		return
	}

	return
}

func GetUser(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow("SELECT * FROM users WHERE username=?", username).Scan(&user.Username)
	return user, err
}

func GetSubjects() (subjects []*Subject, err error) {
	rows, err := DB.Query("SELECT * FROM subjects")
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
	rows, err := DB.Query("SELECT * FROM topics WHERE subject_name=?", subjectName)
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
	rows, err := DB.Query("SELECT * FROM threads WHERE subject_name=? AND topic_name=?", subjectName, topicName)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		thread := &Thread{}
		rows.Scan(&thread.ID, &thread.Title, &thread.Content, &thread.Score, &thread.SubjectName, &thread.TopicName,
			&thread.CreatedByUsername)
		threads = append(threads, thread)
	}
	return
}

func GetThread(threadID int) (*Thread, error) {
	thread := &Thread{}
	err := DB.QueryRow("SELECT * FROM threads WHERE id=?", threadID).Scan(&thread.ID, &thread.Title, &thread.Content,
		&thread.Score, &thread.SubjectName, &thread.TopicName, &thread.CreatedByUsername)
	return thread, err
}
