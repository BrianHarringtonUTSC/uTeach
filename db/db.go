// Package db provides database functionality for the uTeach models using an sqlite database.
package db

import (
	"database/sql/driver"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // blank identifier import registers the sqlite driver

	"github.com/umairidris/uTeach/models"
)

// DB wraps a generic sql DB to provide db functionality for models.
type DB struct {
	*sqlx.DB
}

// New opens a connection to an sqlite database at path and creates all necessary tables that the app uses.
func New(path string) *DB {
	sqlDb := sqlx.MustOpen("sqlite3", path)

	db := &DB{sqlDb}
	db.createTables()
	db.MustExec("PRAGMA foreign_keys=ON;")
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
		CREATE TABLE IF NOT EXISTS threads(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			subject_name TEXT NOT NULL,
			created_by_username TEXT NOT NULL,
			FOREIGN KEY(subject_name) REFERENCES subjects(name) ON DELETE CASCADE,
			FOREIGN KEY(created_by_username) REFERENCES users(username) ON DELETE CASCADE
		)`)

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS upvotes(
			username TEXT  NOT NULL,
			thread_id INTEGER NOT NULL,
			PRIMARY KEY (username, thread_id),
			FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE,
			FOREIGN KEY(thread_id) REFERENCES threads(id) ON DELETE CASCADE
		)`)
}

func (db *DB) exec(query string, params ...interface{}) (driver.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(params...)
}

// User gets the user at the given username.
func (db *DB) User(username string) (*models.User, error) {
	user := &models.User{}
	err := db.Get(user, "SELECT * FROM users WHERE username=?", username)
	return user, err
}

func (db *DB) AddUser(username string) (*models.User, error) {
	if len(username) == 0 {
		return nil, errors.New("Empty username")
	}
	_, err := db.exec("INSERT INTO users(username) VALUES(?)", username)
	if err != nil {
		return nil, err
	}

	return db.User(username)
}

// Subjects gets all subjects.
func (db *DB) Subjects() (subjects []*models.Subject, err error) {
	err = db.Select(&subjects, "SELECT * FROM subjects")
	return
}

// Threads gets all threads with the given subject.
func (db *DB) Threads(subjectName string) (threads []*models.Thread, err error) {
	query := `SELECT threads.*, count(upvotes.thread_id) as score
			  FROM threads LEFT OUTER JOIN upvotes ON threads.id=upvotes.thread_id
			  WHERE threads.subject_name=?
			  GROUP BY threads.id
			  ORDER BY count(upvotes.thread_id) DESC`
	err = db.Select(&threads, query, subjectName)
	return
}

// UserCreatedThreads gets all threads created by the user.
func (db *DB) UserCreatedThreads(username string) (threads []*models.Thread, err error) {
	query := `SELECT threads.*, count(upvotes.thread_id) as score
			  FROM threads LEFT OUTER JOIN upvotes ON threads.id=upvotes.thread_id
			  WHERE threads.created_by_username=?
			  GROUP BY threads.id
			  ORDER BY count(upvotes.thread_id) DESC`
	err = db.Select(&threads, query, username)
	return
}

// Thread gets the thread with the given id.
func (db *DB) Thread(id int64) (thread *models.Thread, err error) {
	query := `SELECT threads.*, count(upvotes.thread_id) as score
			  FROM threads LEFT OUTER JOIN upvotes ON threads.id=upvotes.thread_id
			  WHERE threads.id=?
			  GROUP BY threads.id`
	thread = &models.Thread{}
	err = db.Get(thread, query, id)
	return
}

// UserUpvotedThreadIDs returns the IDs of the threads that the user has upvoted.
// All threadIDs are mapped to "true". The purpose of the map is to act as a set.
func (db *DB) UserUpvotedThreadIDs(username string) (threadIDs map[int64]bool, err error) {
	rows, err := db.Query("SELECT thread_id FROM upvotes WHERE username=?", username)
	if err != nil {
		return
	}
	defer rows.Close()

	threadIDs = map[int64]bool{}
	var threadID int64
	for rows.Next() {
		rows.Scan(&threadID)
		threadIDs[threadID] = true
	}
	return
}

// NewThread adds a new thread.
func (db *DB) NewThread(title string, content string, subject_name string, created_by_username string) (*models.Thread,
	error) {

	if title == "" || content == "" || subject_name == "" || created_by_username == "" {
		return errors.New("Empty values not allowed.")
	}

	query := "INSERT INTO threads(title, content, subject_name, created_by_username) VALUES(?, ?, ?, ?)"
	result, err := db.exec(query, title, content, subject_name, created_by_username)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.Thread(id)
}

// AddUpVote adds upvote for user on the thread.
func (db *DB) AddUpVote(username string, threadID int64) error {
	_, err := db.exec("INSERT INTO upvotes(username, thread_id) VALUES(?, ?)", username, threadID)
	return err
}

// RemoveUpvote removes the vote for user on the thread.
func (db *DB) RemoveUpvote(username string, threadID int64) error {
	_, err := db.exec("DELETE FROM upvotes where username=? AND thread_id=?", username, threadID)
	return err
}
