// Package libtemplate provides template related functions.
package libtemplate

import (
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

// LoadTemplates gets all templates at path into a mapping of the template name to its template object.
// The path should contain a file "base.html" which is the base template.
// It should also contain a "layouts" subfolder which contains child templates to join with the base.
func LoadTemplates(path string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	funcMap := template.FuncMap{
		"markdownToHTML":        MarkdownToHTML,
		"formatAndLocalizeTime": FormatAndLocalizeTime}

	baseTemplate := template.Must(template.New("base").Funcs(funcMap).ParseFiles(filepath.Join(path, "base.html")))

	layoutFiles, _ := filepath.Glob(filepath.Join(path, "layouts/*.html"))
	for _, layoutFile := range layoutFiles {
		baseTemplateCopy, err := baseTemplate.Clone()
		if err != nil {
			panic(err)
		}
		templates[filepath.Base(layoutFile)] = template.Must(baseTemplateCopy.ParseFiles(layoutFile))
	}

	return templates
}
