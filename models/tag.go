package models

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func NewTagModel(db *sqlx.DB) *TagModel {
	return &TagModel{Base{db}}
}

type TagModel struct {
	Base
}

type Tag struct {
	ID      int64
	Name    string
	Subject *Subject
}

var tagsSqlizer = squirrel.
	Select("tags.id, tags.name, subjects.id AS subject_id, subjects.name AS subject_name, subjects.title").
	From("tags").
	Join("subjects ON subjects.id=tags.subject_id")

func (tm *TagModel) getAll(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) ([]*Tag, error) {
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

func (tm *TagModel) getOne(tx *sqlx.Tx, sqlizer squirrel.Sqlizer) (*Tag, error) {
	tags, err := tm.getAll(tx, sqlizer)
	if err != nil {
		return nil, err
	}

	if len(tags) != 1 {
		return nil, fmt.Errorf("Expected: 1, got: %d.", len(tags))
	}

	return tags[0], err
}

func (tm *TagModel) GetTagByID(tx *sqlx.Tx, id int64) (*Tag, error) {
	return tm.getOne(tx, tagsSqlizer.Where(squirrel.Eq{"tags.id": id}))
}

func (tm *TagModel) GetTagByNameAndSubject(tx *sqlx.Tx, name string, subject *Subject) (*Tag, error) {
	return tm.getOne(tx, tagsSqlizer.Where(squirrel.Eq{"tags.name": name, "tags.subject_id": subject.ID}))
}

func (tm *TagModel) GetTagsBySubject(tx *sqlx.Tx, subject *Subject) ([]*Tag, error) {
	return tm.getAll(tx, tagsSqlizer.Where(squirrel.Eq{"subject_id": subject.ID}))
}

func (tm *TagModel) GetThreadsByTag(tx *sqlx.Tx, tag *Tag) ([]*Thread, error) {
	threadModel := NewThreadModel(tm.db)
	threads, err := threadModel.getAll(tx,
		threadsSqlizer.Join("thread_tags ON thread_tags.thread_id=threads.id").Where(squirrel.Eq{"thread_tags.tag_id": tag.ID}))
	return threads, err
}

func (tm *TagModel) AddThreadTag(tx *sqlx.Tx, thread *Thread, tag *Tag) error {
	_, err := tm.Exec(tx, "INSERT INTO thread_tags(thread_id, tag_id, subject_id) VALUES(?, ?, ?)",
		thread.ID, tag.ID, thread.Subject.ID)
	return err
}
