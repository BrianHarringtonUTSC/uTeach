package models

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
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

func (tm *ThreadModel) getThreads(eq squirrel.Eq) ([]*Thread, error) {
	threads := []*Thread{}

	builder := squirrel.
		Select("threads.*, count(thread_votes.thread_id) as score").
		From("threads").
		LeftJoin("thread_votes ON threads.id=thread_votes.thread_id").
		Where(eq).
		GroupBy("threads.id").
		OrderBy("count(thread_votes.thread_id) DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	err = tm.db.Select(&threads, query, args...)
	return threads, err
}

func (tm *ThreadModel) getOneThread(eq squirrel.Eq) (*Thread, error) {
	threads, err := tm.getThreads(eq)
	if err != nil {
		return nil, err
	}
	if len(threads) != 1 {
		return nil, fmt.Errorf("Expected: 1, got: %d.", len(threads))
	}
	return threads[0], err
}

func (tm *ThreadModel) GetThreadByID(id int64) (*Thread, error) {
	return tm.getOneThread(squirrel.Eq{"threads.id": id})
}

// GetThreadsBySubject gets all threads with the given subject.
func (tm *ThreadModel) GetThreadsBySubjectAndIsPinned(subject string, isPinned bool) ([]*Thread, error) {
	return tm.getThreads(squirrel.Eq{"threads.subject_name": subject, "threads.is_pinned": isPinned})
}

// GetThreadsByEmail gets all threads created by the user with the email.
func (tm *ThreadModel) GetThreadsByEmail(email string) ([]*Thread, error) {
	return tm.getThreads(squirrel.Eq{"threads.created_by_email": email})
}

func (tm *ThreadModel) GetThreadIdsUpvotedByEmail(email string) (map[int64]bool, error) {
	rows, err := tm.db.Query("SELECT thread_id FROM thread_votes WHERE email=?", email)
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

func (tm *ThreadModel) AddThread(title, content, subject_name, created_by_email string) (*Thread, error) {
	if title == "" || content == "" || subject_name == "" || created_by_email == "" {
		return nil, errors.New("Empty values not allowed.")
	}

	query := "INSERT INTO threads(title, content, subject_name, created_by_email) VALUES(?, ?, ?, ?)"
	result, err := tm.exec(query, title, content, subject_name, created_by_email)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return tm.GetThreadByID(id)
}

func (tm *ThreadModel) AddThreadVoteForUser(threadID int64, email string) error {
	_, err := tm.exec("INSERT INTO thread_votes(email, thread_id) VALUES(?, ?)", email, threadID)
	return err
}

func (tm *ThreadModel) RemoveTheadVoteForUser(threadID int64, email string) error {
	_, err := tm.exec("DELETE FROM thread_votes where email=? AND thread_id=?", email, threadID)
	return err
}

func (tm *ThreadModel) HideThread(id int64) error {
	_, err := tm.exec("UPDATE threads SET is_visible=? WHERE id=?", false, id)
	return err
}

func (tm *ThreadModel) UnhideThread(id int64) error {
	_, err := tm.exec("UPDATE threads SET is_visible=? WHERE id=?", true, id)
	return err
}

func (tm *ThreadModel) PinThread(id int64) error {
	_, err := tm.exec("UPDATE threads SET is_pinned=? WHERE id=?", true, id)
	return err
}

func (tm *ThreadModel) UnpinThread(id int64) error {
	_, err := tm.exec("UPDATE threads SET is_pinned=? WHERE id=?", false, id)
	return err
}
