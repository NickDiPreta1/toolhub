package web

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/NickDiPreta1/toolhub/internal/tools/fileconvert"
	"github.com/NickDiPreta1/toolhub/internal/tools/jsonutil"
	"github.com/NickDiPreta1/toolhub/internal/tools/textutil"
)

func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	data := &templateData{}
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

type FileConvertData struct {
	Error string
}

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

		// this is resource management - tells how much can be stored in ram
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

type SlugifyData struct {
	Error  string
	Input  string
	Output string
}

func (app *Application) slugify(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := &templateData{
			ToolData: &SlugifyData{},
		}
		app.render(w, http.StatusOK, "slugify.tmpl.html", data)

	case http.MethodPost:
		input := r.FormValue("input")

		if strings.TrimSpace(input) == "" {
			data := &templateData{
				ToolData: &SlugifyData{
					Error: "Please enter some text to slugify.",
				},
			}
			app.render(w, http.StatusBadRequest, "slugify.tmpl.html", data)
			return
		}

		slug := textutil.Slugify(input)
		data := &templateData{
			ToolData: &SlugifyData{
				Input:  input,
				Output: slug,
			},
		}
		app.render(w, http.StatusOK, "slugify.tmpl.html", data)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

type JSONFormatterData struct {
	Error  string
	Input  string
	Output string
	Mode   string
}

func (app *Application) jsonFormatter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := &templateData{
			ToolData: &JSONFormatterData{},
		}
		app.render(w, http.StatusOK, "json.tmpl.html", data)
		return

	case http.MethodPost:
		input := r.FormValue("input")
		mode := r.FormValue("mode")
		if mode == "" {
			mode = "pretty"
		}

		if strings.TrimSpace(input) == "" {
			data := &templateData{
				ToolData: &JSONFormatterData{
					Input: input,
					Error: "Input cannot be empty.",
				},
			}
			app.render(w, http.StatusBadRequest, "json.tmpl.html", data)
			return
		}

		if mode == "minify" {
			minified, err := jsonutil.Minify(input)
			if err != nil {
				data := &templateData{
					ToolData: &JSONFormatterData{
						Input: input,
						Error: err.Error(),
					},
				}
				app.render(w, http.StatusBadRequest, "json.tmpl.html", data)
				return
			}

			data := &templateData{
				ToolData: &JSONFormatterData{
					Input:  input,
					Output: minified,
				},
			}
			app.render(w, http.StatusOK, "json.tmpl.html", data)
			return
		}

		pretty, err := jsonutil.PrettyPrint(input)
		if err != nil {
			data := &templateData{
				ToolData: &JSONFormatterData{
					Input: input,
					Error: err.Error(),
				},
			}
			app.render(w, http.StatusBadRequest, "json.tmpl.html", data)
			return
		}

		data := &templateData{
			ToolData: &JSONFormatterData{
				Input:  input,
				Output: pretty,
			},
		}
		app.render(w, http.StatusOK, "json.tmpl.html", data)
		return

	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
