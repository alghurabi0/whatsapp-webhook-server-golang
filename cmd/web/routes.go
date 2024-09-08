package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /webhook", app.verifyHook)
	mux.HandleFunc("POST /webhook", app.processPayload)
	mux.HandleFunc("POST /send", app.sendMessage)
	return mux
}
