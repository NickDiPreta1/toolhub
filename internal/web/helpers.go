package web

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError logs the error with a stack trace and returns a 500 response.
func (app *Application) serverError(w http.ResponseWriter, err error) {
	// Include a stack trace to make server-side logs actionable.
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
