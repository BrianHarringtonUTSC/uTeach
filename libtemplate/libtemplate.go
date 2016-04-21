// Package libtemplate provides template related functions.
package libtemplate

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

// FormatAndLocalizeTime a time to remove excess information and make it local to the server.
func FormatAndLocalizeTime(t time.Time) string {
	return t.Local().Format("Jan 2 2006 3:04PM")
}

// Dict creates map given a sequence of values. Each even index must be a string and is mapped to the next value.
// Useful for calling a template with parameters.
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}

	dict := map[string]interface{}{}

	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// HTML unescapes HTML in the string. It should be sanitized prior to this function.
func HTML(s string) template.HTML {
	return template.HTML(s)
}

// Load gets all templates at path into a mapping of the template name to its template object.
// The path should contain a layouts/ subdirectory with all the templates. The path should also contain a includes/
// subdirectory which contains parent and reusable templates, they will be parsed with each template in the layouts/
// directory. See: https://elithrar.github.io/article/approximating-html-template-inheritance/ for implementation
// details.
func Load(path string) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	funcMap := template.FuncMap{
		"dict":                  Dict,
		"formatAndLocalizeTime": FormatAndLocalizeTime,
		"html":                  HTML,
	}

	layouts, err := filepath.Glob(filepath.Join(path, "layouts/*.html"))
	if err != nil {
		return nil, err
	}

	includes, err := filepath.Glob(filepath.Join(path, "includes/*.html"))
	if err != nil {
		return nil, err
	}

	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.New(layout).Funcs(funcMap).ParseFiles(files...))
	}

	return templates, err
}

// Render renders the template at name with data and writes out the result.
func Render(w http.ResponseWriter, templates map[string]*template.Template, name string, data map[string]interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	// TODO: to speed this up use a buffer pool (https://elithrar.github.io/article/using-buffer-pools-with-go/)
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}
	buf.WriteTo(w)
	return nil
}
