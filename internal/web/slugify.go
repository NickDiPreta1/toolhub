package web

import (
	"net/http"
	"strings"

	"github.com/NickDiPreta1/toolhub/internal/tools/textutil"
)

type SlugifyData struct {
	Error  string
	Input  string
	Output string
}

// slugify renders the slugify tool and processes user input.
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
