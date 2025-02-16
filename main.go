package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/karaMuha/go-chirpy/rest"
	"github.com/karaMuha/go-chirpy/service"
	"github.com/karaMuha/go-chirpy/sql/repositories"
	"github.com/karaMuha/go-chirpy/state"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}

	platform := os.Getenv("PLATFORM")

	userRepo := repositories.NewUsersRepository(db)
	chirpRepo := repositories.NewChirpsRepository(db)

	userService := service.NewUsersService(userRepo)
	chripsService := service.NewChripsService(chirpRepo)
	service := service.NewService()

	appState := state.NewAppState(platform)

	restHandler := rest.NewRestHandler(appState, service, userService, chripsService)
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
	apiHandler.HandleFunc("POST /validate_chirp", handler.HandleValidateChirp)
	apiHandler.HandleFunc("POST /users", handler.HandleCreateUser)
	apiHandler.HandleFunc("POST /chirps", handler.HandleCreateChirp)
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	adminHandler := http.NewServeMux()
	adminHandler.HandleFunc("GET /metrics", handler.HandleViewMetrics)
	adminHandler.HandleFunc("POST /reset", handler.HandleReset)
	mux.Handle("/admin/", http.StripPrefix("/admin", adminHandler))
}
