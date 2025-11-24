package web

import (
	"fmt"
	"net/http"
)

func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}
