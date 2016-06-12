// Package httperror provides functionality to handler errors for http requests.
package httperror

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/BrianHarringtonUTSC/uTeach/models"
	"github.com/pkg/errors"
)

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Stacktrace interface to get err stack info
type Stacktrace interface {
	Stacktrace() []errors.Frame
}

// Error allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	if se.Err == nil {
		return fmt.Sprintf("%d %s", se.Code, http.StatusText(se.Code))
	}
	return fmt.Sprintf("%d %s", se.Code, se.Err.Error())
}

// HandleError handles error messaging for the client and server. Internal server errors are logged and not written to
// client to not expose sensitive information.
func HandleError(w http.ResponseWriter, err error) {
	cause := errors.Cause(err)

	if cause == sql.ErrNoRows {
		cause = StatusError{http.StatusNotFound, nil}
	}

	if err != nil {
		switch e := cause.(type) {
		case StatusError:
			http.Error(w, e.Error(), e.Code)
		case models.InputError:
			http.Error(w, e.Error(), http.StatusBadRequest)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Printf("%+v\n", err)

		}
	}
}
