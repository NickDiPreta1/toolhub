package web

import (
	"fmt"
	"net/http"
)

func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	panic("boom")
	app.render(w, http.StatusOK, "home.tmpl.html")
}
