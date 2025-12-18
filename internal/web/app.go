package web

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Application holds shared dependencies for handlers and middleware.
type Application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
}

// NewApplication wires up dependencies and builds the initial template cache.
func NewApplication(infoLog, errorLog *log.Logger) (*Application, error) {
	tc, err := newTemplateCache()
	if err != nil {
		return nil, err
	}

	app := &Application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
	}

	return app, nil
}

// render executes a named template into a buffer, then writes it to the response.
func (app *Application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the templates %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
