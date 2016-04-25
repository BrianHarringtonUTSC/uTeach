package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/libtemplate"
	"github.com/umairidris/uTeach/models"
)

func getTopics(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewTopicModel(a.DB)
	topics, err := tm.Find(nil)
	if err != nil {
		return errors.Wrap(err, "find error")
	}

	data := context.TemplateData(r)
	data["Topics"] = topics
	err = libtemplate.Render(w, a.Templates, "topics.html", data)
	return errors.Wrap(err, "render template error")
}

func getNewTopic(a *application.App, w http.ResponseWriter, r *http.Request) error {
	err := libtemplate.Render(w, a.Templates, "new_topic.html", context.TemplateData(r))
	return errors.Wrap(err, "render template error")
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
