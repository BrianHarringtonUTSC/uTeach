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

func (sm *SubjectModel) GetAllSubjects() ([]*Subject, error) {
	subjects := []*Subject{}
	err := sm.db.Select(&subjects, "SELECT * FROM subjects")
	return subjects, err
}

func (sm *SubjectModel) GetSubjectByName(name string) (*Subject, error) {
	subject := &Subject{}
	err := sm.db.Get(subject, "SELECT * FROM subjects WHERE name=?", name)
	return subject, err
}
