package main

import (
	"log"
	"os"
)

//ILogger ...
type ILogger interface {
	Log(v ...interface{})
	Logln(v ...interface{})
	Logf(format string, v ...interface{})
}

//APILogger ...
type APILogger struct {
	log *log.Logger
}

//Log ...
func (t *APILogger) Log(v ...interface{}) {
	if t.log != nil {
		t.log.Print(v...)
	}
}

//Logln ...
func (t *APILogger) Logln(v ...interface{}) {
	if t.log != nil {
		t.log.Println(v...)
	}
}

//Logf ..
func (t *APILogger) Logf(format string, v ...interface{}) {
	if t.log != nil {
		t.log.Printf(format, v...)
	}
}

//NewAPILogger ...
func NewAPILogger(prefix string) *APILogger {
	return &APILogger{
		log: log.New(os.Stdout, prefix, log.Ldate|log.Ltime),
	}
}
