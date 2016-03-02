package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
)

// GetSubjects renders all subjects.
func GetSubjects(w http.ResponseWriter, r *http.Request) {
	app := application.Get(r)
	subjects, err := app.DB.Subjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	renderTemplate(w, r, "subjects.html", data)
}
