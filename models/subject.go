package models

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Subject represents a subject in the app.
type Subject struct {
	ID    int64
	Name  string
	Title string
}

// URL returns the unique URL for a subject.
func (s *Subject) URL() string {
	return "/s/" + s.Name
}

// SubjectModel handles getting and creating subjects.
type SubjectModel struct {
	Base
}

var subjectNameRegex = regexp.MustCompile(`^[[:alnum:]]+(_[[:alnum:]]+)*$`)

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

// GetSubjectByID gets a subject by id.
func (sm *SubjectModel) GetSubjectByID(tx *sqlx.Tx, id int64) (*Subject, error) {
	subject := new(Subject)
	err := sm.Get(tx, subject, "SELECT * FROM subjects WHERE id=?", id)
	return subject, err
}

// GetSubjectByName gets a subject by name.
func (sm *SubjectModel) GetSubjectByName(tx *sqlx.Tx, name string) (*Subject, error) {
	name = strings.ToLower(name)

	subject := new(Subject)
	err := sm.Get(tx, subject, "SELECT * FROM subjects WHERE name=?", name)
	return subject, err
}

func (sm *SubjectModel) AddSubject(tx *sqlx.Tx, name, title string) (*Subject, error) {
	if title == "" || !subjectNameRegex.MatchString(name) {
		fmt.Println(name)
		fmt.Println(title)
		return nil, InputError{"Invalid name and/or title."}
	}

	name = strings.ToLower(name)

	query := "INSERT INTO subjects(name, title) VALUES(?, ?)"
	result, err := sm.Exec(tx, query, name, title)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return sm.GetSubjectByID(tx, id)
}
