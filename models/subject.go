package models

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// Subject represents a subject in the app.
type Subject struct {
	ID    int64
	Name  string
	Title string
}

// SubjectModel handles getting and creating subjects.
type SubjectModel struct {
	Base
}

// NewSubjectModel returns a new subject model.
func NewSubjectModel(db *sqlx.DB) *SubjectModel {
	return &SubjectModel{Base{db}}
}

// GetAllSubjects gets all subjects.
func (sm *SubjectModel) GetAllSubjects(tx *sqlx.Tx) ([]*Subject, error) {
	subjects := []*Subject{}
	err := sm.Select(tx, &subjects, "SELECT * FROM subjects")
	return subjects, err
}

// GetSubjectByName gets a subject by name.
func (sm *SubjectModel) GetSubjectByName(tx *sqlx.Tx, name string) (*Subject, error) {
	name = strings.ToLower(name)
	subject := &Subject{}
	err := sm.Get(tx, subject, "SELECT * FROM subjects WHERE name=?", name)
	return subject, err
}
