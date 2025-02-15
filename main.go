package main

import (
	"log"
	"net/http"

	"github.com/karaMuha/go-chirpy/rest"
	"github.com/karaMuha/go-chirpy/state"
)

func main() {
	appState := state.NewAppState()
	restHandler := rest.NewRestHandler(appState)
	mux := http.NewServeMux()
	setupEndpoints(mux, restHandler, appState)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func setupEndpoints(mux *http.ServeMux, handler rest.RestHandler, appState *state.AppState) {
	pathToStatic := http.Dir("./static")
	fileServer := http.FileServer(pathToStatic)
	fileHandlerWithMiddleware := appState.IncMetrics(fileServer)
	mux.Handle("GET /app", http.StripPrefix("/app", fileHandlerWithMiddleware))

	mux.HandleFunc("GET /healthz", handler.HandleHealthCheck)
	mux.HandleFunc("GET /metrics", handler.HandleViewMetrics)
	mux.HandleFunc("POST /reset", handler.HandleResetViewCount)
}
