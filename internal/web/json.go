package web

import (
	"net/http"
	"strings"

	"github.com/NickDiPreta1/toolhub/internal/tools/jsonutil"
)

type JSONFormatterData struct {
	Error  string
	Input  string
	Output string
	Mode   string
}

// jsonFormatter formats or minifies JSON submitted by the user.
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
