package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// LoadTemplates gets all templates at path into a mapping of the template name to its template object.
// The path should contain a file "base.html" which is the base template.
// It should also contain a "layout" subfolder which contains child templates to join with the base.
func LoadTemplates(path string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	baseTemplate := template.Must(template.ParseFiles(filepath.Join(path, "base.html")))

	layoutFiles, _ := filepath.Glob(filepath.Join(path, "layout/*.html"))
	for _, layoutFile := range layoutFiles {
		baseTemplateCopy, err := baseTemplate.Clone()
		if err != nil {
			panic(err)
		}
		templates[filepath.Base(layoutFile)] = template.Must(baseTemplateCopy.ParseFiles(layoutFile))
	}

	return templates
}

// RenderTemplate renders the template at name with data.
// It also adds the session user to the data for templates to access.
func (a *App) RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	tmpl, ok := a.templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", name), http.StatusInternalServerError)
		return
	}

	user, ok := a.store.SessionUser(r)
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
