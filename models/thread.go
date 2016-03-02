package models

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewThreadModel(db *sqlx.DB) *ThreadModel {
	return &ThreadModel{Base{db}}
}

type Thread struct {
	ID                int64
	Title             string
	Content           string
	Score             int
	SubjectName       string `db:"subject_name"`
	CreatedByUsername string `db:"created_by_username"`
}

type ThreadModel struct {
	Base
}

// URL returns the unique URL for a thread.
func (t *Thread) URL() string {
	return fmt.Sprintf("/s/%s/%d", t.SubjectName, t.ID)
}

// GetAllThreads gets all threads with the given subject.
func (t *ThreadModel) GetThreadsBySubject(subject string) ([]*Thread, error) {
	threads := []*Thread{}
	query := `SELECT threads.*, count(upvotes.thread_id) as score
			  FROM threads LEFT OUTER JOIN upvotes ON threads.id=upvotes.thread_id
			  WHERE threads.subject_name=?
			  GROUP BY threads.id
			  ORDER BY count(upvotes.thread_id) DESC`
	err := t.db.Select(&threads, query, subject)
	return threads, err
}

func (t *ThreadModel) GetThreadByID(id int64) (*Thread, error) {
	thread := &Thread{}
	query := `SELECT threads.*, count(upvotes.thread_id) as score
			  FROM threads LEFT OUTER JOIN upvotes ON threads.id=upvotes.thread_id
			  WHERE threads.id=?
			  GROUP BY threads.id`
	err := t.db.Get(thread, query, id)
	return thread, err
}

// GetThreadsByUsername gets all threads created by the user.
func (t *ThreadModel) GetThreadsByUsername(username string) ([]*Thread, error) {
	threads := []*Thread{}
	query := `SELECT threads.*, count(upvotes.thread_id) as score
			  FROM threads LEFT OUTER JOIN upvotes ON threads.id=upvotes.thread_id
			  WHERE threads.created_by_username=?
			  GROUP BY threads.id
			  ORDER BY count(upvotes.thread_id) DESC`
	err := t.db.Select(&threads, query, username)
	return threads, err
}

func (t *ThreadModel) GetThreadIdsUpvotedByUsername(username string) (map[int64]bool, error) {
	rows, err := t.db.Query("SELECT thread_id FROM upvotes WHERE username=?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threadIDs := map[int64]bool{}
	var threadID int64
	for rows.Next() {
		rows.Scan(&threadID)
		threadIDs[threadID] = true
	}
	return threadIDs, err
}

func (t *ThreadModel) AddThread(title, content, subject_name, created_by_username string) (*Thread, error) {
	if title == "" || content == "" || subject_name == "" || created_by_username == "" {
		return nil, errors.New("Empty values not allowed.")
	}

	query := "INSERT INTO threads(title, content, subject_name, created_by_username) VALUES(?, ?, ?, ?)"
	result, err := t.exec(query, title, content, subject_name, created_by_username)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return t.GetThreadByID(id)
}

func (t *ThreadModel) AddThreadVoteForUser(threadID int64, username string) error {
	_, err := t.exec("INSERT INTO upvotes(username, thread_id) VALUES(?, ?)", username, threadID)
	return err
}

func (t *ThreadModel) RemoveTheadVoteForUser(threadID int64, username string) error {
	_, err := t.exec("DELETE FROM upvotes where username=? AND thread_id=?", username, threadID)
	return err
}
