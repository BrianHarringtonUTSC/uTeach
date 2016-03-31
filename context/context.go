// Package context provides convenience setters and getters for request-scoped context values.
package context

import (
	"net/http"

	"github.com/gorilla/context"
)

// TODO: replace this with context which will be added in go 1.7
const (
	threadIDKey = "threadID"
)

func SetThreadID(r *http.Request, threadID int64) {
	context.Set(r, threadIDKey, threadID)
}

func ThreadID(r *http.Request) int64 {
	return context.Get(r, threadIDKey).(int64)
}
