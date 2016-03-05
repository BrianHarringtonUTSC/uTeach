package context

import (
	"github.com/gorilla/context"
	"net/http"

	"github.com/umairidris/uTeach/application"
)

const (
	appKey      = "app"
	threadIDKey = "threadID"
)

// SetContext sets the global application in the request.
func SetApp(r *http.Request, a *application.App) {
	context.Set(r, appKey, a)
}

// Get gets the global application from the request. Application MUST be set in the context before, else the application
// will be nil. See SetContext helper function.
func GetApp(r *http.Request) *application.App {
	return context.Get(r, appKey).(*application.App)
}

func SetThreadID(r *http.Request, threadID int64) {
	context.Set(r, threadIDKey, threadID)
}

func GetThreadID(r *http.Request) int64 {
	return context.Get(r, threadIDKey).(int64)
}
