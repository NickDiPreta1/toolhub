package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/NickDiPreta1/toolhub/internal/tools/fileconvert"
)

func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.tmpl.html")
}

func (app *Application) fileConverterForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	app.render(w, http.StatusOK, "fileconvert.tmpl.html")
}

func (app *Application) fileConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const maxUploadSize = 2 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// this is resource management - tells how much can be stored in ram
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "unable to parse form (file too large or invalid)", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file field is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	mode := r.FormValue("mode")
	if mode == "" {
		mode = "uppercase"
	}

	converted, err := fileconvert.ToUpperText(file)
	if err != nil {
		app.serverError(w, err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", `attachment; filename="converted.txt"`)
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, converted)
	if err != nil {
		app.errorLog.Printf("error sending converted file: %v", err)
		return
	}
}
