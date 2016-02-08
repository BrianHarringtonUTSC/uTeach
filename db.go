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
			title TEXT NOT NULL,
			content TEXT NOT NULL,
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
			username TEXT NOT NULL,
			thread_id INTEGER NOT NULL,
			FOREIGN KEY(username) REFERENCES users(username)
			FOREIGN KEY(thread_id) REFERENCES threads(rowid)
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
	query := `SELECT threads.rowid, threads.*, count(upvotes.thread_id)
			  FROM threads LEFT OUTER JOIN upvotes ON threads.rowid=upvotes.thread_id
			  WHERE threads.subject_name=? AND threads.topic_name=?
			  GROUP BY threads.rowid
			  ORDER BY count(upvotes.thread_id) DESC`
	rows, err := DB.Query(query, subjectName, topicName)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		thread := &Thread{}
		rows.Scan(&thread.ID, &thread.Title, &thread.Content, &thread.SubjectName, &thread.TopicName,
			&thread.CreatedByUsername, &thread.Score)
		threads = append(threads, thread)
	}
	return
}

func GetThreadScore(threadID int) (score int, err error) {
 	err = DB.QueryRow("SELECT COUNT(*) FROM upvotes WHERE thread_id=?", threadID).Scan(&score)
 	return
 }

func GetThread(threadID int) (*Thread, error) {
	thread := &Thread{}
	err := DB.QueryRow("SELECT rowid, * FROM threads WHERE rowid=?", threadID).Scan(&thread.ID, &thread.Title,
		&thread.Content, &thread.SubjectName, &thread.TopicName, &thread.CreatedByUsername)
	if err != nil {
		return nil, err
	}

	thread.Score, err = GetThreadScore(thread.ID)
	return thread, err
}

func GetUserUpvotedThreadIDs(username string) (threadIDs map[int]bool, err error) {
	rows, err := DB.Query("SELECT thread_id FROM upvotes WHERE username=?", username)
	if err != nil {
		return
	}
	defer rows.Close()

	threadIDs = make(map[int]bool)
	for rows.Next() {
		var threadID int
		rows.Scan(&threadID)
		threadIDs[threadID] = true
	}
	return
}

func runUpvoteQuery(query string, username string, threadID int) (err error) {
	stmt, err := DB.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, threadID)
	return
}

func AddUpVote(username string, threadID int) error {
	return runUpvoteQuery("INSERT INTO upvotes(username, thread_id) VALUES(?, ?)", username, threadID)
}

func RemoveUpvote(username string, threadID int) error {
	return runUpvoteQuery("DELETE FROM upvotes where username=? AND thread_id=?", username, threadID)
}
