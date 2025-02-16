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
	secret := os.Getenv("SECRET")
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	polkaKey := os.Getenv("POLKA_KEY")

	appState := state.NewAppState(platform)
	appState.Secret = secret
	appState.PolkaKey = polkaKey

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error while validating database connection: %v", err)
	}

	chirpRepo := repositories.NewChirpsRepository(db)
	userRepo := repositories.NewUsersRepository(db)

	userService := service.NewUsersService(userRepo, appState)
	chripsService := service.NewChripsService(chirpRepo)
	service := service.NewService()

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
	apiHandler.HandleFunc("GET /chirps", handler.HandleGetAllChirps)
	apiHandler.HandleFunc("GET /chirps/{chirpID}", handler.HandleGetChirpByID)
	apiHandler.HandleFunc("POST /login", handler.HandleLogin)
	apiHandler.HandleFunc("PUT /users", handler.HandleUpdateAccount)
	apiHandler.HandleFunc("DELETE /chirps/{chirpID}", handler.HandleDeleteChirp)
	apiHandler.HandleFunc("POST /polka/webhooks", handler.HandleUpgradeToRed)
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	adminHandler := http.NewServeMux()
	adminHandler.HandleFunc("GET /metrics", handler.HandleViewMetrics)
	adminHandler.HandleFunc("POST /reset", handler.HandleReset)
	mux.Handle("/admin/", http.StripPrefix("/admin", adminHandler))
}
