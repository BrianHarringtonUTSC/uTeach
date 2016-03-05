package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
)

// GetSubjects renders all subjects.
func GetSubjects(w http.ResponseWriter, r *http.Request) {
	app := context.GetApp(r)

	sm := models.NewSubjectModel(app.DB)
	subjects, err := sm.GetAllSubjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	renderTemplate(w, r, "subjects.html", data)
}
