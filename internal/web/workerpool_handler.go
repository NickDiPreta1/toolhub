package web

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/NickDiPreta1/toolhub/internal/tools/encodingutil"
	"github.com/NickDiPreta1/toolhub/internal/tools/hashutil"
	"github.com/NickDiPreta1/toolhub/internal/workerpool"
)

type WorkerPoolResult struct {
	JobID    int
	Filename string
	Content  string
	// Duration time.Duration TODO
	Error string
}

type WorkerPoolData struct {
	Results     []WorkerPoolResult
	WorkerCount int
	Error       string
}

func (app *Application) workerPool(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := &templateData{
			ToolData: &WorkerPoolData{},
		}
		app.render(w, http.StatusOK, "workerpool.tmpl.html", data)
		return

	case http.MethodPost:
		const maxUploadSize = 10 * 1024 * 1024
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			data := &templateData{
				ToolData: &WorkerPoolData{
					Error: "File too large or invalid upload.",
				},
			}
			app.render(w, http.StatusBadRequest, "workerpool.tmpl.html", data)
			return
		}

		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			data := &templateData{
				ToolData: &WorkerPoolData{
					Error: "Error: please upload at least a few files.",
				},
			}
			app.render(w, http.StatusBadRequest, "workerpool.tmpl.html", data)
			return
		}

		workersStr := r.FormValue("workers")
		workerCount, err := strconv.Atoi(workersStr)
		if err != nil {
			workerCount = 3
		}

		functionChoice := r.FormValue("function")
		if functionChoice == "" {
			functionChoice = "hash"
		}

		var processFunc func([]byte) ([]byte, error)
		switch functionChoice {
		case "hash":
			processFunc = func(b []byte) ([]byte, error) {
				result, err := hashutil.Hash(b)
				return []byte(result), err
			}
		case "uppercase":
			processFunc = func(b []byte) ([]byte, error) {
				upper := bytes.ToUpper(b)
				return upper, nil
			}
		case "base64encode":
			processFunc = func(b []byte) ([]byte, error) {
				encoded := encodingutil.Encode(string(b))
				return []byte(encoded), nil
			}
		case "base64decode":
			processFunc = func(b []byte) ([]byte, error) {
				decoded, err := encodingutil.Decode(string(b))
				if err != nil {
					return nil, err
				}
				return []byte(decoded), nil
			}
		default:
			processFunc = func(b []byte) ([]byte, error) {
				result, err := hashutil.Hash(b)
				return []byte(result), err
			}
		}

		ctx := r.Context()
		pool := workerpool.NewPool(workerCount, len(files))
		resultsChan := pool.Start(ctx)
		jobFilenames := make(map[int]string)

		var results []WorkerPoolResult
		done := make(chan struct{})

		go func() {
			for result := range resultsChan {
				filename := jobFilenames[result.JobID]
				res := WorkerPoolResult{
					JobID:    result.JobID,
					Filename: filename,
					Content:  string(result.Content),
				}
				results = append(results, res)
			}
			close(done)
		}()

		for i, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				app.errorLog.Printf("Failed to open %s: %v", fileHeader.Filename, err)
				continue
			}
			content, err := io.ReadAll(file)
			file.Close()
			jobFilenames[i] = fileHeader.Filename
			pool.Submit(workerpool.Job{
				ID:      i,
				Content: content,
				Func:    processFunc,
			})
		}

		pool.Shutdown()
		<-done

		data := &templateData{
			ToolData: &WorkerPoolData{
				Results:     results,
				WorkerCount: workerCount,
			},
		}

		app.render(w, http.StatusOK, "workerpool.tmpl.html", data)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
