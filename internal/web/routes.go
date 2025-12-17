package web

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", app.Ping)
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/tools/fileconvert", app.fileConvert)
	mux.HandleFunc("/tools/slugify", app.slugify)
	mux.HandleFunc("/tools/json", app.jsonFormatter)

	return app.PanicRecover(app.LogRequest(mux))
}
