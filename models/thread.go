package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// Thread represents a thread in the app.
type Thread struct {
	ID          int64
	Title       string
	Content     string
	TimeCreated time.Time
	IsPinned    bool
	IsVisible   bool
	Score       int
	Subject     *Subject
	Creator     *User
}

// URL returns the unique URL for a thread.
func (t *Thread) URL() string {
	return fmt.Sprintf("/t/%d", t.ID)
}

// ThreadModel handles getting and creating threads.
type ThreadModel struct {
	Base
}

// NewThreadModel returns a new thread model.
func NewThreadModel(db *sqlx.DB) *ThreadModel {
	return &ThreadModel{Base{db}}
}

var threadsSqlizer = squirrel.
	Select(`threads.id AS thread_id,
			threads.title AS thread_title,
			threads.content,
			threads.time_created,
			threads.is_pinned,
			threads.is_visible,
			count(thread_votes.thread_id),
			subjects.id AS subject_id,
			subjects.name AS subject_name,
			subjects.title AS subject_title,
			users.id AS user_id,
			users.email,
			users.name AS user_name,
			users.is_admin`).
	From("threads").
	Join("subjects ON subjects.id=threads.subject_id").
	Join("users ON users.id=threads.creator_user_id").
	LeftJoin("thread_votes ON thread_votes.thread_id=threads.id").
	GroupBy("threads.id").
	OrderBy("count(thread_votes.thread_id) DESC")

func (tm *ThreadModel) findAll(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) ([]*Thread, error) {
	threads := []*Thread{}

	query, args, err := sqlizer.ToSql()
	if err != nil {
		return threads, err
	}

	rows, err := tm.Query(tx, query, args...)
	if err != nil {
		return threads, err
	}
	defer rows.Close()

	for rows.Next() {
		thread := &Thread{}
		subject := &Subject{}
		creator := &User{}

		err = rows.Scan(&thread.ID, &thread.Title, &thread.Content, &thread.TimeCreated, &thread.IsPinned, &thread.IsVisible, &thread.Score,
			&subject.ID, &subject.Name, &subject.Title,
			&creator.ID, &creator.Email, &creator.Name, &creator.IsAdmin)
		if err != nil {
			return threads, err
		}

		thread.Subject = subject
		thread.Creator = creator
		threads = append(threads, thread)
	}

	return threads, err
}

func (tm *ThreadModel) findOne(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) (*Thread, error) {
	threads, err := tm.findAll(tx, sqlizer)
	if err != nil {
		return nil, err
	}
	if len(threads) != 1 {
		return nil, fmt.Errorf("Expected: 1, got: %d.", len(threads))
	}
	return threads[0], err
}

// GetThreadByID gets a thread by the id.
func (tm *ThreadModel) GetThreadByID(tx *sqlx.Tx, id int64) (*Thread, error) {
	return tm.findOne(tx, threadsSqlizer.Where(squirrel.Eq{"threads.id": id}))
}

// GetThreadsBySubjectAndIsPinned gets all threads by subject and whether they are pinned or not pinned.
func (tm *ThreadModel) GetThreadsBySubjectAndIsPinned(tx *sqlx.Tx, subject *Subject, isPinned bool) ([]*Thread, error) {
	threads, err := tm.findAll(tx, threadsSqlizer.Where(squirrel.Eq{"threads.subject_id": subject.ID, "threads.is_pinned": isPinned}))
	if err == sql.ErrNoRows {
		return []*Thread{}, nil
	}
	return threads, err
}

// GetThreadsByUser gets all threads by the user.
func (tm *ThreadModel) GetThreadsByUser(tx *sqlx.Tx, user *User) ([]*Thread, error) {
	return tm.findAll(tx, threadsSqlizer.Where(squirrel.Eq{"threads.creator_user_id": user.ID}))
}

// GetThreadIdsUpvotedByUser gets the ids of all threads upvoted by the user. It returns a map which can be used to
// check if a thread was upvoted by a user in constant time.
// TODO: this method may need to be made more precise. For example, finding all upvoted threads for a subject, etc.
func (tm *ThreadModel) GetThreadIdsUpvotedByUser(tx *sqlx.Tx, user *User) (map[int64]bool, error) {
	rows, err := tm.Query(tx, "SELECT thread_id FROM thread_votes WHERE user_id=?", user.ID)
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

// AddThread adds a new thread.
func (tm *ThreadModel) AddThread(tx *sqlx.Tx, title, content string, subject *Subject, creator *User) (*Thread, error) {
	if title == "" || content == "" {
		return nil, errors.New("Empty values not allowed.")
	}

	query := "INSERT INTO threads(title, content, subject_id, creator_user_id) VALUES(?, ?, ?, ?)"
	result, err := tm.Exec(tx, query, title, content, subject.ID, creator.ID)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()

	return tm.GetThreadByID(tx, id)
}

// AddThreadVoteForUser adds a vote for the thread for the user.
func (tm *ThreadModel) AddThreadVoteForUser(tx *sqlx.Tx, thread *Thread, user *User) error {
	_, err := tm.Exec(tx, "INSERT INTO thread_votes(user_id, thread_id) VALUES(?, ?)", user.ID, thread.ID)
	return err
}

// RemoveTheadVoteForUser removes a vote for the thread for the user.
func (tm *ThreadModel) RemoveTheadVoteForUser(tx *sqlx.Tx, thread *Thread, user *User) error {
	_, err := tm.Exec(tx, "DELETE FROM thread_votes where user_id=? AND thread_id=?", user.ID, thread.ID)
	return err
}

// HideThread hides the thread.
func (tm *ThreadModel) HideThread(tx *sqlx.Tx, thread *Thread) error {
	_, err := tm.Exec(tx, "UPDATE threads SET is_visible=? WHERE id=?", false, thread.ID)
	return err
}

// UnhideThread unhides the thread.
func (tm *ThreadModel) UnhideThread(tx *sqlx.Tx, thread *Thread) error {
	_, err := tm.Exec(tx, "UPDATE threads SET is_visible=? WHERE id=?", true, thread.ID)
	return err
}

// PinThread pins a thread.
func (tm *ThreadModel) PinThread(tx *sqlx.Tx, thread *Thread) error {
	_, err := tm.Exec(tx, "UPDATE threads SET is_pinned=? WHERE id=?", true, thread.ID)
	return err
}

// UnpinThread unpins a thread.
func (tm *ThreadModel) UnpinThread(tx *sqlx.Tx, thread *Thread) error {
	_, err := tm.Exec(tx, "UPDATE threads SET is_pinned=? WHERE id=?", false, thread.ID)
	return err
}

// GetThreadsByTag gets all threads with tag.
func (tm *ThreadModel) GetThreadsByTag(tx *sqlx.Tx, tag *Tag) ([]*Thread, error) {
	threads, err := tm.findAll(tx,
		threadsSqlizer.Join("thread_tags ON thread_tags.thread_id=threads.id").Where(squirrel.Eq{"thread_tags.tag_id": tag.ID}))
	return threads, err
}
