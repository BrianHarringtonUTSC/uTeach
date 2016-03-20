package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
)

// GetSubjects renders all subjects.
func GetSubjects(w http.ResponseWriter, r *http.Request) {
	sm := models.NewSubjectModel(context.DB(r))
	subjects, err := sm.GetAllSubjects()
	if err != nil {
		handleError(w, err)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	renderTemplate(w, r, "subjects.html", data)
}
