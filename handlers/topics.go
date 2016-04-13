package handlers

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/models"
)

func getTopics(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewTopicModel(a.DB)
	topics, err := tm.GetAllTopics(nil)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Topics": topics}
	return renderTemplate(a, w, r, "topics.html", data)
}

func getNewTopic(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(a, w, r, "new_topic.html", nil)
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
