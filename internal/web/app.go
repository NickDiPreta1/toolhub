package web

import "log"

type Application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewApplication(infoLog, errorLog *log.Logger) *Application {
	return &Application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}
