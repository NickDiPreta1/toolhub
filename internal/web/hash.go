package web

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/NickDiPreta1/toolhub/internal/tools/hashutil"
)

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
