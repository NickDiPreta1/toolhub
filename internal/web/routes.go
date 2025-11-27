package web

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", app.Ping)
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/tools/fileconvert", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.fileConverterForm(w, r)
		case http.MethodPost:
			app.fileConvert(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

	})

	return app.PanicRecover(app.LogRequest(mux))
}
