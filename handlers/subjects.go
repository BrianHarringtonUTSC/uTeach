package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/models"
)

// GetSubjects renders all subjects.
func GetSubjects(w http.ResponseWriter, r *http.Request) {
	app := application.GetFromContext(r)

	s := models.NewSubjectModel(app.DB)
	subjects, err := s.GetAllSubjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	renderTemplate(w, r, "subjects.html", data)
}
