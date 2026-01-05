package web

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/NickDiPreta1/toolhub/internal/tools/fileconvert"
)

type FileConvertData struct {
	Error string
}

// fileConvert handles file upload, conversion, and download.
func (app *Application) fileConvert(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := &templateData{
			ToolData: &FileConvertData{},
		}
		app.render(w, http.StatusOK, "fileconvert.tmpl.html", data)
	case http.MethodPost:
		const maxUploadSize = 2 * 1024 * 1024
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		// Limit total form data to keep memory usage predictable.
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			data := &templateData{
				ToolData: &FileConvertData{
					Error: "File too large or invalid upload. Maximum size is 2MB.",
				},
			}
			app.render(w, http.StatusBadRequest, "fileconvert.tmpl.html", data)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			data := &templateData{
				ToolData: &FileConvertData{
					Error: "Please choose a file to convert.",
				},
			}
			app.render(w, http.StatusBadRequest, "fileconvert.tmpl.html", data)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if ext != ".txt" {
			data := &templateData{
				ToolData: &FileConvertData{
					Error: "Only .txt files are supported right now.",
				},
			}
			app.render(w, http.StatusBadRequest, "fileconvert.tmpl.html", data)
			return
		}

		mode := r.FormValue("mode")
		if mode == "" {
			mode = "uppercase"
		}

		converted, err := fileconvert.ToUpperText(file)
		if err != nil {
			app.serverError(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", `attachment; filename="converted.txt"`)
		w.WriteHeader(http.StatusOK)

		_, err = io.Copy(w, converted)
		if err != nil {
			app.errorLog.Printf("error sending converted file: %v", err)
			return
		}
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
