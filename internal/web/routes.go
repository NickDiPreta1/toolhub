package web

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", app.Ping)
	mux.HandleFunc("/", app.home)
	return mux
}
