package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// DB wraps a generic sql DB to provide db functionality for models.
type DB struct {
	*sql.DB
}

// panicOnErr will panic if the err is not nil.
func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

// NewDB opens a connection to an sqlite database at path and creates all necessary tables that the app uses.
func NewDB(path string) *DB {

	sqlDb, err := sql.Open("sqlite3", path)
	panicOnErr(err)

	db := &DB{sqlDb}

	// ensure all required tables exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			username TEXT PRIMARY KEY
		)`)
	panicOnErr(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS subjects(
			name TEXT PRIMARY KEY,
			title TEXT NOT NULL
		)`)
	panicOnErr(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS topics(
			name TEXT NOT NULL,
			title TEXT NOT NULL,
			subject_name TEXT NOT NULL,
			PRIMARY KEY(name, subject_name)
			FOREIGN KEY(subject_name) REFERENCES subjects(name)
		)`)
	panicOnErr(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS threads(
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			subject_name TEXT NOT NULL,
			topic_name TEXT NOT NULL,
			created_by_username TEXT NOT NULL,
			FOREIGN KEY(subject_name, topic_name) REFERENCES topics(subject_name, name)
			FOREIGN KEY(created_by_username) REFERENCES users(username)
		)`)
	panicOnErr(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS upvotes(
			username TEXT NOT NULL,
			thread_id INTEGER NOT NULL,
			FOREIGN KEY(username) REFERENCES users(username)
			FOREIGN KEY(thread_id) REFERENCES threads(rowid)
		)`)
	panicOnErr(err)

	return db
}

// User gets the User at username.
func (db *DB) User(username string) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT * FROM users WHERE username=?", username).Scan(&user.Username)
	return user, err
}

// Subjects gets all subjects.
func (db *DB) Subjects() (subjects []*Subject, err error) {
	rows, err := db.Query("SELECT * FROM subjects")
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

// Topics gets all topics with the subjectName.
func (db *DB) Topics(subjectName string) (topics []*Topic, err error) {
	rows, err := db.Query("SELECT * FROM topics WHERE subject_name=?", subjectName)
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

// Threads gets all threads with the subjectName and topicName.
func (db *DB) Threads(subjectName string, topicName string) (threads []*Thread, err error) {
	query := `SELECT threads.rowid, threads.*, count(upvotes.thread_id)
			  FROM threads LEFT OUTER JOIN upvotes ON threads.rowid=upvotes.thread_id
			  WHERE threads.subject_name=? AND threads.topic_name=?
			  GROUP BY threads.rowid
			  ORDER BY count(upvotes.thread_id) DESC`
	rows, err := db.Query(query, subjectName, topicName)
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

// ThreadScore gets the total upvotes for a thread at the given threadID.
func (db *DB) ThreadScore(threadID int) (score int, err error) {
	err = db.QueryRow("SELECT COUNT(*) FROM upvotes WHERE thread_id=?", threadID).Scan(&score)
	return
}

// Thread gets the thread with the given id.
func (db *DB) Thread(id int) (*Thread, error) {
	thread := &Thread{}
	err := db.QueryRow("SELECT rowid, * FROM threads WHERE rowid=?", id).Scan(&thread.ID, &thread.Title,
		&thread.Content, &thread.SubjectName, &thread.TopicName, &thread.CreatedByUsername)
	if err != nil {
		return nil, err
	}

	thread.Score, err = db.ThreadScore(thread.ID)
	return thread, err
}

// UserUpvotedThreadIDs returns the IDs of the threads the user with username has upvoted.
func (db *DB) UserUpvotedThreadIDs(username string) (threadIDs map[int]bool, err error) {
	rows, err := db.Query("SELECT thread_id FROM upvotes WHERE username=?", username)
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

func (db *DB) runUpvoteQuery(query string, username string, threadID int) (err error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, threadID)
	return
}

// AddUpVote adds upvote for user with username for the thread with threadID.
func (db *DB) AddUpVote(username string, threadID int) error {
	return db.runUpvoteQuery("INSERT INTO upvotes(username, thread_id) VALUES(?, ?)", username, threadID)
}

// RemoveUpvote removes the vote for user with username for the thread with threadID.
func (db *DB) RemoveUpvote(username string, threadID int) error {
	return db.runUpvoteQuery("DELETE FROM upvotes where username=? AND thread_id=?", username, threadID)
}
