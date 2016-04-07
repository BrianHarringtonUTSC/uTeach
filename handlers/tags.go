package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
)

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
