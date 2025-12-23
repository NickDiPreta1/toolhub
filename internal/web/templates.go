package web

import (
	"html/template"
	"path/filepath"
)

// templateData is the shared view model for templates.
type templateData struct {
	Error     string
	Flash     string
	PageTitle string
	ToolData  any
}

// newTemplateCache parses page, base, and partial templates into a lookup map.
// Templates are parsed in a specific order to support template inheritance:
// 1. Parse the base layout template (defines the "base" template block)
// 2. Parse all partial templates (reusable components like nav, footer)
// 3. Parse the specific page template (extends "base" with page-specific content)
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// Step 1: Parse base layout template
		ts, err := template.New(name).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Step 2: Parse all partial templates (nav, footer, etc.)
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Step 3: Parse the specific page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
