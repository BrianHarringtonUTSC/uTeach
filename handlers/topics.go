package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/libtemplate"
	"github.com/umairidris/uTeach/models"
)

func getTopics(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewTopicModel(a.DB)
	topics, err := tm.Find(nil)
	if err != nil {
		return err
	}

	data := context.TemplateData(r)
	data["Topics"] = topics
	return libtemplate.Render(w, a.Templates, "topics.html", data)
}

func getNewTopic(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return libtemplate.Render(w, a.Templates, "new_topic.html", context.TemplateData(r))
}

func postNewTopic(a *application.App, w http.ResponseWriter, r *http.Request) error {
	name := r.FormValue("name")
	title := r.FormValue("title")
	description := r.FormValue("description")

	tm := models.NewTopicModel(a.DB)
	topic, err := tm.AddTopic(nil, name, title, description)
	if err != nil {
		return err
	}

	http.Redirect(w, r, topic.URL(), http.StatusFound)
	return nil
}
