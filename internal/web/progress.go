package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) progressDemo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := &templateData{}
		app.render(w, http.StatusOK, "progress.tmpl.html", data)
		return

	case http.MethodPost:
		var chunks, timeout int
		chunkString := r.FormValue("chunks")
		timeoutString := r.FormValue("timeout")

		chunks, err := strconv.Atoi(chunkString)
		if err != nil {
			chunks = 20
		}
		timeout, err = strconv.Atoi(timeoutString)
		if err != nil {
			timeout = 30
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		flusher, ok := w.(http.Flusher)
		if !ok {
			app.serverError(w, errors.New("Streaming not supported"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(timeout))
		defer cancel()

		fmt.Fprintf(w, "Starting operation\n")
		flusher.Flush()

		for i := 1; i <= chunks; i++ {
			err := progressHelper(ctx, i)
			if err != nil {
				fmt.Fprintf(w, "Cancelled at chunk %d of %d: %v\n", i, chunks, err)
				flusher.Flush()
				return
			}
			fmt.Fprintf(w, "Processed chunk %d of %d\n", i, chunks)
			flusher.Flush()
		}

		fmt.Fprintf(w, "Done writing %d chunks\n", chunks)
		flusher.Flush()

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func progressHelper(ctx context.Context, chunkNumber int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		time.Sleep(500 * time.Millisecond)
		return nil
	}
}
