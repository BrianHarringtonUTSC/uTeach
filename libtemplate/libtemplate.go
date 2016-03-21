// Package libtemplate provides template related functions.
package libtemplate

import (
	"log"
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"path/filepath"
	"time"
)

// MarkdownToHTML converts a markdown string into HTML.
func MarkdownToHTML(markdown string) template.HTML {
	unsafe := blackfriday.MarkdownBasic([]byte(markdown))
	safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return template.HTML(safe)
}

func FormatAndLocalizeTime(t time.Time) string {
	return t.Local().Format("Jan 2 2006 3:04PM")
}

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

// Load gets all templates at path into a mapping of the template name to its template object.
// https://elithrar.github.io/article/approximating-html-template-inheritance/
func Load(path string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	funcMap := template.FuncMap{
		"dict":                  Dict,
		"markdownToHTML":        MarkdownToHTML,
		"formatAndLocalizeTime": FormatAndLocalizeTime}

	layouts, err := filepath.Glob(filepath.Join(path, "layouts/*.html"))
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(filepath.Join(path, "includes/*.html"))
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.New(layout).Funcs(funcMap).ParseFiles(files...))
	}

	return templates
}
