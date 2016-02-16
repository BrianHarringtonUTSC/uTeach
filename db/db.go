// Package db provides functionaltiy database functionaltiy for the uTeach models using an sqlite database.
package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/umairidris/uTeach/models"
)

// DB wraps a generic sql DB to provide db functionality for models.
type DB struct {
	*sqlx.DB
}

// NewDB opens a connection to an sqlite database at path and creates all necessary tables that the app uses.
func New(path string) *DB {
	sqlDb := sqlx.MustOpen("sqlite3", path)

	db := &DB{sqlDb}
	db.createTables()

	return db
}

func (db *DB) createTables() {
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS users(
			username TEXT PRIMARY KEY
		)`)

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS subjects(
			name TEXT PRIMARY KEY,
			title TEXT NOT NULL
		)`)

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS topics(
			name TEXT NOT NULL,
			title TEXT NOT NULL,
			subject_name TEXT NOT NULL,
			PRIMARY KEY(name, subject_name)
			FOREIGN KEY(subject_name) REFERENCES subjects(name)
		)`)

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS threads(
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			subject_name TEXT NOT NULL,
			topic_name TEXT NOT NULL,
			created_by_username TEXT NOT NULL,
			FOREIGN KEY(subject_name, topic_name) REFERENCES topics(subject_name, name)
			FOREIGN KEY(created_by_username) REFERENCES users(username)
		)`)

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS upvotes(
			username TEXT NOT NULL,
			thread_id INTEGER NOT NULL,
			FOREIGN KEY(username) REFERENCES users(username)
			FOREIGN KEY(thread_id) REFERENCES threads(rowid)
		)`)
}

// User gets the user at the given username.
func (db *DB) User(username string) (user *models.User, err error) {
	user = &models.User{}
	err = db.Get(user, "SELECT * FROM users WHERE username=?", username)
	return
}

// Subjects gets all subjects.
func (db *DB) Subjects() (subjects []*models.Subject, err error) {
	err = db.Select(&subjects, "SELECT * FROM subjects")
	return
}

// Topics gets all topics with the given subject name.
func (db *DB) Topics(subjectName string) (topics []*models.Topic, err error) {
	err = db.Select(&topics, "SELECT * FROM topics WHERE subject_name=?", subjectName)
	return
}

// Threads gets all threads with the given subject and topic names.
func (db *DB) Threads(subjectName string, topicName string) (threads []*models.Thread, err error) {
	query := `SELECT threads.rowid, threads.*, count(upvotes.thread_id) as score
		FROM threads LEFT OUTER JOIN upvotes ON threads.rowid=upvotes.thread_id
		WHERE threads.subject_name=? AND threads.topic_name=?
		GROUP BY threads.rowid
		ORDER BY count(upvotes.thread_id) DESC`
	err = db.Select(&threads, query, subjectName, topicName)
	return
}

// UserCreatedThreads gets all threads created by the user.
func (db *DB) UserCreatedThreads(username string) (threads []*models.Thread, err error) {
	query := `SELECT threads.rowid, threads.*, count(upvotes.thread_id) as score
			  FROM threads LEFT OUTER JOIN upvotes ON threads.rowid=upvotes.thread_id
			  WHERE threads.created_by_username=?
			  GROUP BY threads.rowid
			  ORDER BY count(upvotes.thread_id) DESC`
	err = db.Select(&threads, query, username)
	return
}

// Thread gets the thread with the given id.
func (db *DB) Thread(id int) (thread *models.Thread, err error) {
	query := `SELECT threads.rowid, threads.*, count(upvotes.thread_id) as score
			  FROM threads LEFT OUTER JOIN upvotes ON threads.rowid=upvotes.thread_id
			  WHERE threads.rowid=?
			  GROUP BY threads.rowid`
	thread = &models.Thread{}
	err = db.Get(thread, query, id)
	return
}

// UserUpvotedThreadIDs returns the IDs of the threads that the user has upvoted.
// All threadIDs are mapped to "true". The purpose of the map is to act as a set.
func (db *DB) UserUpvotedThreadIDs(username string) (threadIDs map[int]bool, err error) {
	rows, err := db.Query("SELECT thread_id FROM upvotes WHERE username=?", username)
	if err != nil {
		return
	}
	defer rows.Close()

	threadIDs = make(map[int]bool)
	var threadID int
	for rows.Next() {
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

// AddUpVote adds upvote for user on the thread.
func (db *DB) AddUpVote(username string, threadID int) error {
	return db.runUpvoteQuery("INSERT INTO upvotes(username, thread_id) VALUES(?, ?)", username, threadID)
}

// RemoveUpvote removes the vote for user on the thread.
func (db *DB) RemoveUpvote(username string, threadID int) error {
	return db.runUpvoteQuery("DELETE FROM upvotes where username=? AND thread_id=?", username, threadID)
}
