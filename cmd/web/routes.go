package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /webhook", app.verifyHook)
	mux.HandleFunc("POST /webhook", app.processPayload)

	mux.HandleFunc("GET /chat", app.chat)
	mux.HandleFunc("POST /message", app.sendMessage)
	return mux
}
