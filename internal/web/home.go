package web

import (
	"net/http"
)

// home renders the landing page.
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	data := &templateData{}
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}
