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
	fsHandler := http.FileServer(pathToStatic)
	fsHandlerWithMiddleware := appState.IncMetrics(fsHandler)
	mux.Handle("/app/", http.StripPrefix("/app", fsHandlerWithMiddleware))

	apiHandler := http.NewServeMux()
	apiHandler.HandleFunc("GET /healthz", handler.HandleHealthCheck)
	apiHandler.HandleFunc("GET /metrics", handler.HandleViewMetrics)
	apiHandler.HandleFunc("POST /reset", handler.HandleResetViewCount)

	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	adminHandler := http.NewServeMux()
	adminHandler.HandleFunc("GET /healthz", handler.HandleHealthCheck)
	adminHandler.HandleFunc("GET /metrics", handler.HandleViewMetrics)
	adminHandler.HandleFunc("POST /reset", handler.HandleResetViewCount)

	mux.Handle("/admin/", http.StripPrefix("/admin", adminHandler))
}
