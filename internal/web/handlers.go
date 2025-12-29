package web

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/NickDiPreta1/toolhub/internal/tools/encodingutil"
	"github.com/NickDiPreta1/toolhub/internal/tools/fileconvert"
	"github.com/NickDiPreta1/toolhub/internal/tools/hashutil"
	"github.com/NickDiPreta1/toolhub/internal/tools/jsonutil"
	"github.com/NickDiPreta1/toolhub/internal/tools/textutil"
)

// Ping provides a simple health check endpoint.
func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

// home renders the landing page.
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	data := &templateData{}
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

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

type HashResult struct {
	Filename string
	Hash     string
	Error    string
}

type ConcurrentHashData struct {
	Results []HashResult
	Error   string
}

func (app *Application) concurrentHash(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := &templateData{
			ToolData: &ConcurrentHashData{},
		}
		app.render(w, http.StatusOK, "hash.tmpl.html", data)
		return
	case http.MethodPost:
		const maxUploadSize = 10 * 1024 * 1024
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			data := &templateData{
				ToolData: &ConcurrentHashData{
					Error: "File too large or invalid upload.",
				},
			}
			app.render(w, http.StatusBadRequest, "hash.tmpl.html", data)
			return
		}

		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			data := &templateData{
				ToolData: &ConcurrentHashData{
					Error: "Error: please upload at least one file",
				},
			}
			app.render(w, http.StatusBadRequest, "hash.tmpl.html", data)
			return
		}

		resChan := make(chan HashResult, len(files))

		for _, fileHeader := range files {
			go hashFile(fileHeader, resChan)
		}

		timeout := time.After(10 * time.Second)
		var successes, failures []HashResult
		for i := 0; i < len(files); i++ {
			select {
			case r := <-resChan:
				if r.Error != "" {
					failures = append(failures, r)
				} else {
					successes = append(successes, r)
				}
			case <-timeout:
				failures = append(failures, HashResult{
					Error: "Hashing took too long.",
				})
			}

		}

		data := &templateData{
			ToolData: &ConcurrentHashData{
				Results: successes,
			},
		}

		app.render(w, http.StatusOK, "hash.tmpl.html", data)

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return

	}

}

func hashFile(fh *multipart.FileHeader, results chan<- HashResult) {
	file, err := fh.Open()
	if err != nil {
		results <- HashResult{
			Filename: fh.Filename,
			Error:    "Error opening file.",
		}
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		results <- HashResult{
			Filename: fh.Filename,
			Error:    "Error reading file.",
		}
		return
	}
	hashedFile, err := hashutil.Hash(data)
	if err != nil {
		results <- HashResult{
			Filename: fh.Filename,
			Error:    fmt.Sprintf("Error: error hashing %s", fh.Filename),
		}
		return
	}
	results <- HashResult{
		Filename: fh.Filename,
		Hash:     hashedFile,
	}
}
