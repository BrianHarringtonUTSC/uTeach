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

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}
