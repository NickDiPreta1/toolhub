package web

import (
	"fmt"
	"net/http"
)

// Ping provides a simple health check endpoint.
func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}
