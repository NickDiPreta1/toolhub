// Package main is the entry point for the toolhub web application.
// It starts an HTTP server that provides various web-based utilities
// including text processing, encoding, and file conversion tools.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/NickDiPreta1/toolhub/internal/web"
)

func main() {
	addr := flag.String("addr", ":4000", "TCP address to serve on")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app, err := web.NewApplication(infoLog, errorLog)
	if err != nil {
		errorLog.Fatal(err)
	}

	srv := http.Server{
		Addr:     *addr,
		Handler:  app.Routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
