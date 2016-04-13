package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/models"
)

func getSubjects(a *application.App, w http.ResponseWriter, r *http.Request) error {
	sm := models.NewSubjectModel(a.DB)
	subjects, err := sm.GetAllSubjects(nil)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Subjects": subjects}
	return renderTemplate(a, w, r, "subjects.html", data)
}

func getNewSubject(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(a, w, r, "new_subject.html", nil)
}

func postNewSubject(a *application.App, w http.ResponseWriter, r *http.Request) error {
	name := r.FormValue("name")
	title := r.FormValue("title")
	description := r.FormValue("description")

	sm := models.NewSubjectModel(a.DB)
	subject, err := sm.AddSubject(nil, name, title, description)
	if err != nil {
		return err
	}

	http.Redirect(w, r, subject.URL(), http.StatusFound)
	return nil
}
