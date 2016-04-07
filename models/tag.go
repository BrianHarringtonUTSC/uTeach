package models

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// Thread represents a tag in the app.
type Tag struct {
	ID      int64
	Name    string
	Subject *Subject
}

// TagModel handles getting and creating tags.
type TagModel struct {
	Base
}

// NewThreadModel returns a new tag model.
func NewTagModel(db *sqlx.DB) *TagModel {
	return &TagModel{Base{db}}
}

var tagsSqlizer = squirrel.
	Select("tags.id, tags.name, subjects.id AS subject_id, subjects.name AS subject_name, subjects.title").
	From("tags").
	Join("subjects ON subjects.id=tags.subject_id")

func (tm *TagModel) findAll(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) ([]*Tag, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tm.Query(tx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []*Tag{}
	for rows.Next() {
		tag := &Tag{}
		subject := &Subject{}
		err := rows.Scan(&tag.ID, &tag.Name, &subject.ID, &subject.Name, &subject.Title)
		if err != nil {
			return nil, err
		}
		tag.Subject = subject
		tags = append(tags, tag)
	}
	return tags, err
}

func (tm *TagModel) findOne(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) (*Tag, error) {
	tags, err := tm.findAll(tx, sqlizer)
	if err != nil {
		return nil, err
	}

	if len(tags) != 1 {
		return nil, fmt.Errorf("Expected: 1, got: %d.", len(tags))
	}

	return tags[0], err
}

// GetTagByID gets a tag by the id.
func (tm *TagModel) GetTagByID(tx *sqlx.Tx, id int64) (*Tag, error) {
	return tm.findOne(tx, tagsSqlizer.Where(squirrel.Eq{"tags.id": id}))
}

// GetTagByNameAndSubject gets a tag by the name and subject.
func (tm *TagModel) GetTagByNameAndSubject(tx *sqlx.Tx, name string, subject *Subject) (*Tag, error) {
	return tm.findOne(tx, tagsSqlizer.Where(squirrel.Eq{"tags.name": name, "tags.subject_id": subject.ID}))
}

// GetTagsBySubject gets all tags by the subject.
func (tm *TagModel) GetTagsBySubject(tx *sqlx.Tx, subject *Subject) ([]*Tag, error) {
	return tm.findAll(tx, tagsSqlizer.Where(squirrel.Eq{"subject_id": subject.ID}))
}

// AddThreadTag adds a tag for the thread.
func (tm *TagModel) AddThreadTag(tx *sqlx.Tx, thread *Thread, tag *Tag) error {
	_, err := tm.Exec(tx, "INSERT INTO thread_tags(thread_id, tag_id, subject_id) VALUES(?, ?, ?)",
		thread.ID, tag.ID, thread.Subject.ID)
	return err
}
