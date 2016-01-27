package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var templates map[string]*template.Template

func LoadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	baseTemplate := template.Must(template.ParseFiles("tmpl/base.html"))

	layouts, _ := filepath.Glob("tmpl/layout/*.html")
	for _, layoutFile := range layouts {
		baseTemplateCopy, _ := baseTemplate.Clone()
		templates[filepath.Base(layoutFile)] = template.Must(baseTemplateCopy.ParseFiles(layoutFile))
	}
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "The template %s does not exist.", http.StatusInternalServerError)
		return
	}

	user, ok := getSessionUser(r)
	if !ok {
		// if failed to get user, make sure user is nil so templates don't render a user
		user = nil
	}

	data["User"] = user

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
}
