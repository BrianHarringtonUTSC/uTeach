// Package application provides context to bind the uTeach app together.
package application

import (
	"github.com/gorilla/context"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/umairidris/uTeach/config"
	"github.com/umairidris/uTeach/db"
	"github.com/umairidris/uTeach/session"
)

const (
	contextKey = "application" // key for storing/retrieving Application from gorilla context
)

// Application is the application context which contains application-wide configuration and components.
type Application struct {
	Config    *config.Config
	DB        *db.DB
	Store     *session.Store
	Templates map[string]*template.Template
}

// New initializes a new Application.
func New(configPath string) *Application {
	config := config.Load(configPath)
	db := db.New(config.DBPath)
	store := session.NewStore("todo-proper-secret")
	templates := loadTemplates(config.TemplatesPath)

	app := &Application{config, db, store, templates}

	return app
}

// Get gets the global application from the request. Application MUST be set in the context before, else the application
// will be nil. See SetContext helper function.
func Get(r *http.Request) *Application {
	return context.Get(r, contextKey).(*Application)
}

// SetContext sets the global application in the request.
func SetContext(r *http.Request, app *Application) {
	context.Set(r, contextKey, app)
}

// markdownToHTML converts the markdown into HTML.
func markdownToHTML(markdown string) template.HTML {
	unsafe := blackfriday.MarkdownBasic([]byte(markdown))
	safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return template.HTML(safe)
}

// loadTemplates gets all templates at path into a mapping of the template name to its template object.
// The path should contain a file "base.html" which is the base template.
// It should also contain a "layouts" subfolder which contains child templates to join with the base.
func loadTemplates(path string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	funcMap := template.FuncMap{"markdownToHTML": markdownToHTML}
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
