package models

import (
	"github.com/jmoiron/sqlx"
)

func NewSubjectModel(db *sqlx.DB) *SubjectModel {
	return &SubjectModel{Base{db}}
}

type Subject struct {
	ID    int64
	Name  string
	Title string
}

type SubjectModel struct {
	Base
}

func (sm *SubjectModel) GetAllSubjects(tx *sqlx.Tx) ([]*Subject, error) {
	subjects := []*Subject{}
	err := sm.Select(tx, &subjects, "SELECT * FROM subjects")
	return subjects, err
}

func (sm *SubjectModel) GetSubjectByName(tx *sqlx.Tx, name string) (*Subject, error) {
	subject := &Subject{}
	err := sm.Get(tx, subject, "SELECT * FROM subjects WHERE name=?", name)
	return subject, err
}
