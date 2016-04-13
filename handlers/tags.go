package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
)

func getTags(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)

	tm := models.NewTagModel(a.DB)
	tags, err := tm.GetTagsByTopic(nil, topic)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Tags": tags, "Topic": topic}

	return renderTemplate(a, w, r, "tags.html", data)
}

func getNewTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(a, w, r, "new_tag.html", nil)
}

func postNewTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)

	name := r.FormValue("name")

	tm := models.NewTagModel(a.DB)
	tag, err := tm.AddTag(nil, name, topic)
	if err != nil {
		return err
	}

	http.Redirect(w, r, tag.Topic.URL(), http.StatusFound)
	return nil
}
