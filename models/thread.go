package models

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewThreadModel(db *sqlx.DB) *ThreadModel {
	return &ThreadModel{Base{db}}
}

type Thread struct {
	ID             int64
	Title          string
	Content        string
	Score          int
	SubjectName    string    `db:"subject_name"`
	CreatedByEmail string    `db:"created_by_email"`
	TimeCreated    time.Time `db:"time_created"`
	IsPinned       bool      `db:"is_pinned"`
	IsVisible      bool      `db:"is_visible"`
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
	query := `SELECT threads.*, count(thread_votes.thread_id) as score
			  FROM threads LEFT OUTER JOIN thread_votes ON threads.id=thread_votes.thread_id
			  WHERE threads.subject_name=?
			  GROUP BY threads.id
			  ORDER BY count(thread_votes.thread_id) DESC`
	err := t.db.Select(&threads, query, subject)
	return threads, err
}

func (t *ThreadModel) GetThreadByID(id int64) (*Thread, error) {
	thread := &Thread{}
	query := `SELECT threads.*, count(thread_votes.thread_id) as score
			  FROM threads LEFT OUTER JOIN thread_votes ON threads.id=thread_votes.thread_id
			  WHERE threads.id=?
			  GROUP BY threads.id`
	err := t.db.Get(thread, query, id)
	return thread, err
}

// GetThreadsByEmail gets all threads created by the user with the email.
func (t *ThreadModel) GetThreadsByEmail(email string) ([]*Thread, error) {
	threads := []*Thread{}
	query := `SELECT threads.*, count(thread_votes.thread_id) as score
			  FROM threads LEFT OUTER JOIN thread_votes ON threads.id=thread_votes.thread_id
			  WHERE threads.created_by_email=?
			  GROUP BY threads.id
			  ORDER BY count(thread_votes.thread_id) DESC`
	err := t.db.Select(&threads, query, email)
	return threads, err
}

func (t *ThreadModel) GetThreadIdsUpvotedByEmail(email string) (map[int64]bool, error) {
	rows, err := t.db.Query("SELECT thread_id FROM thread_votes WHERE email=?", email)
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

func (t *ThreadModel) AddThread(title, content, subject_name, created_by_email string) (*Thread, error) {
	if title == "" || content == "" || subject_name == "" || created_by_email == "" {
		return nil, errors.New("Empty values not allowed.")
	}

	query := "INSERT INTO threads(title, content, subject_name, created_by_email) VALUES(?, ?, ?, ?)"
	result, err := t.exec(query, title, content, subject_name, created_by_email)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return t.GetThreadByID(id)
}

func (t *ThreadModel) AddThreadVoteForUser(threadID int64, email string) error {
	_, err := t.exec("INSERT INTO thread_votes(email, thread_id) VALUES(?, ?)", email, threadID)
	return err
}

func (t *ThreadModel) RemoveTheadVoteForUser(threadID int64, email string) error {
	_, err := t.exec("DELETE FROM thread_votes where email=? AND thread_id=?", email, threadID)
	return err
}
