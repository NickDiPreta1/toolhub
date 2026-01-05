package web

import (
	"net/http"
	"strings"

	"github.com/NickDiPreta1/toolhub/internal/tools/encodingutil"
)

type Base64Data struct {
	Error  string
	Input  string
	Output string
	Mode   string
}

// base64Tool encodes or decodes base64 input.
func (app *Application) base64Tool(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := &templateData{
			ToolData: &Base64Data{},
		}
		app.render(w, http.StatusOK, "base64.tmpl.html", data)
		return
	case http.MethodPost:
		input := r.FormValue("input")
		mode := r.FormValue("mode")
		if mode == "" {
			mode = "encode"
		}

		if strings.TrimSpace(input) == "" {
			data := &templateData{
				ToolData: &Base64Data{
					Input: input,
					Error: "Input cannot be empty",
				},
			}

			app.render(w, http.StatusBadRequest, "base64.tmpl.html", data)
			return
		}

		if mode == "decode" {
			decoded, err := encodingutil.Decode(input)
			if err != nil {
				data := &templateData{
					ToolData: &Base64Data{
						Input: input,
						Mode:  "decode",
						Error: "Invalid base64 input",
					},
				}
				app.render(w, http.StatusBadRequest, "base64.tmpl.html", data)
				return
			}

			data := &templateData{
				ToolData: &Base64Data{
					Input:  input,
					Output: decoded,
				},
			}
			app.render(w, http.StatusOK, "base64.tmpl.html", data)
			return
		}

		encoded := encodingutil.Encode(input)
		data := &templateData{
			ToolData: &Base64Data{
				Input:  input,
				Output: encoded,
			},
		}
		app.render(w, http.StatusOK, "base64.tmpl.html", data)
		return

	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
