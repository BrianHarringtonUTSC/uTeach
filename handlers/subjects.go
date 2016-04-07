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
