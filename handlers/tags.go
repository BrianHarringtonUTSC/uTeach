package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/models"
)

// GetSubjects renders all subjects.
func getTags(a *application.App, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	subjectName := strings.ToLower(vars["subject"])
	sm := models.NewSubjectModel(a.DB)
	subject, err := sm.GetSubjectByName(subjectName)
	if err != nil {
		return err
	}

	tm := models.NewTagModel(a.DB)
	tags, err := tm.GetTagsBySubject(subject)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Tags": tags}

	return renderTemplate(a, w, r, "tags.html", data)
}

func getThreadsByTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	subjectName := strings.ToLower(vars["subject"])
	sm := models.NewSubjectModel(a.DB)
	subject, err := sm.GetSubjectByName(subjectName)
	if err != nil {
		return err
	}

	tagName := strings.ToLower(vars["tag"])
	tm := models.NewTagModel(a.DB)
	tag, err := tm.GetTagByNameAndSubject(tagName, subject)
	if err != nil {
		return err
	}

	threads, err := tm.GetThreadsByTag(tag)
	if err != nil {
		return err
	}
	data := map[string]interface{}{"Threads": threads}
	return renderTemplate(a, w, r, "threads_by_tag.html", data)
}
