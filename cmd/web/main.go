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

	app := web.NewApplication(infoLog, errorLog)

	srv := http.Server{
		Addr:     *addr,
		Handler:  app.Routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
