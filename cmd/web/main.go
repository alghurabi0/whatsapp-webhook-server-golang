package main

import (
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}
	srv := http.Server{
		Handler:  app.routes(),
		ErrorLog: errorLog,
		Addr:     ":4002",
	}
	infoLog.Printf("app running on 4002\n")
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
