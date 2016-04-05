package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
)

// GetSubjects renders all subjects.
func getTags(a *application.App, w http.ResponseWriter, r *http.Request) error {
	subject := context.Subject(r)

	tm := models.NewTagModel(a.DB)
	tags, err := tm.GetTagsBySubject(nil, subject)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Tags": tags}

	return renderTemplate(a, w, r, "tags.html", data)
}

func getThreadsByTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tag := context.Tag(r)

	tm := models.NewTagModel(a.DB)
	threads, err := tm.GetThreadsByTag(nil, tag)
	if err != nil {
		return err
	}
	data := map[string]interface{}{"Threads": threads}
	return renderTemplate(a, w, r, "threads_by_tag.html", data)
}
