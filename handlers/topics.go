package handlers

import (
	"net/http"

	"github.com/BrianHarringtonUTSC/uTeach/application"
	"github.com/BrianHarringtonUTSC/uTeach/context"
	"github.com/BrianHarringtonUTSC/uTeach/libtemplate"
	"github.com/BrianHarringtonUTSC/uTeach/models"
	"github.com/pkg/errors"
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

	topic := &models.Topic{Name: name, Title: title, Description: description}
	if err := tm.Add(nil, topic); err != nil {
		return err
	}

	http.Redirect(w, r, topic.URL(), http.StatusFound)
	return nil
}
