// Package httperror provides functionality to handler errors for http requests.
package httperror

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/umairidris/uTeach/models"
)

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
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
			errors.Fprint(os.Stderr, err)
		}
	}
}
