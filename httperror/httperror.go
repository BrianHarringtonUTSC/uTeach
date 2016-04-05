// Package httperror provides functionality to handler errors for http requests.
package httperror

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	if se.Err == nil {
		return fmt.Sprintf("%d %s", se.Code, http.StatusText(se.Code))
	}
	return fmt.Sprintf("%d %s", se.Code, se.Err.Error())
}

func HandleError(w http.ResponseWriter, err error) {
	if err == sql.ErrNoRows {
		err = StatusError{http.StatusNotFound, nil}
	}

	if err != nil {
		switch e := err.(type) {
		case StatusError:
			http.Error(w, e.Error(), e.Code)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println(err)
		}
	}
}
