package web

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
	"time"
)

type FileResult struct {
	Filename string
	Content  string
	Duration time.Duration
}

type ConcurrentUpperData struct {
	Results []FileResult
	Error   string
}

// concurrentUpper will concurrently process files
func (app *Application) concurrentUpper(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := &templateData{
			ToolData: &ConcurrentUpperData{},
		}
		app.render(w, http.StatusOK, "concurrent.tmpl.html", data)
		return
	case http.MethodPost:
		const maxUploadSize = 10 * 1024 * 1024
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			data := &templateData{
				ToolData: &ConcurrentUpperData{
					Error: "File too large or invalid upload.",
				},
			}
			app.render(w, http.StatusBadRequest, "concurrent.tmpl.html", data)
			return
		}

		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			data := &templateData{
				ToolData: &ConcurrentUpperData{
					Error: "Error: please upload at least one file.",
				},
			}
			app.render(w, http.StatusBadRequest, "concurrent.tmpl.html", data)
			return
		}

		// Process files concurrently using goroutines for improved performance.
		// WaitGroup tracks all running goroutines to ensure we wait for completion.
		// Mutex protects the shared results slice from concurrent write races.
		var wg sync.WaitGroup
		var mu sync.Mutex
		var results []FileResult
		for _, fileHeader := range files {
			wg.Add(1)
			go func(fh *multipart.FileHeader) {
				defer wg.Done()
				start := time.Now()
				file, err := fh.Open()
				if err != nil {
					app.errorLog.Printf("Failed to open %s: %v", fh.Filename, err)
					return
				}
				defer file.Close()

				data, err := io.ReadAll(file)
				if err != nil {
					app.errorLog.Printf("Failed to read %s: %v", fh.Filename, err)
					return
				}

				upper := bytes.ToUpper(data)
				duration := time.Since(start)

				// Lock protects against race conditions when multiple goroutines
				// append to the results slice simultaneously.
				mu.Lock()
				results = append(results, FileResult{
					Filename: fh.Filename,
					Content:  string(upper),
					Duration: duration,
				})
				mu.Unlock()
			}(fileHeader)
		}

		// Wait for all goroutines to complete before proceeding
		wg.Wait()

		// Check if all file processing failed. If files were uploaded but
		// no results were produced, all goroutines encountered errors.
		if len(results) == 0 && len(files) > 0 {
			data := &templateData{
				ToolData: &ConcurrentUpperData{
					Error: "Failed to process any files",
				},
			}
			app.render(w, http.StatusInternalServerError, "concurrent.tmpl.html", data)
			return
		}
		data := &templateData{
			ToolData: &ConcurrentUpperData{
				Results: results,
			},
		}

		app.render(w, http.StatusOK, "concurrent.tmpl.html", data)
		return

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
