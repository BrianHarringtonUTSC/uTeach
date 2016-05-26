package handlers

import (
	"net/http"

	"github.com/BrianHarringtonUTSC/uTeach/application"
	"github.com/BrianHarringtonUTSC/uTeach/context"
	"github.com/BrianHarringtonUTSC/uTeach/libtemplate"
	"github.com/BrianHarringtonUTSC/uTeach/models"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func getTags(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)

	tm := models.NewTagModel(a.DB)
	tags, err := tm.Find(nil, squirrel.Eq{"tags.topic_id": topic.ID})
	if err != nil {
		return errors.Wrap(err, "find error")
	}

	data := context.TemplateData(r)
	data["Tags"] = tags
	data["Topic"] = topic

	err = libtemplate.Render(w, a.Templates, "tags.html", data)
	return errors.Wrap(err, "render template error")
}

func getNewTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	err := libtemplate.Render(w, a.Templates, "new_tag.html", context.TemplateData(r))
	return errors.Wrap(err, "render template error")
}

func postNewTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)
	name := r.FormValue("name")

	tm := models.NewTagModel(a.DB)
	tag := &models.Tag{Name: name, Topic: topic}
	if err := tm.Add(nil, tag); err != nil {
		return err
	}

	http.Redirect(w, r, tag.Topic.URL(), http.StatusFound)
	return nil
}
