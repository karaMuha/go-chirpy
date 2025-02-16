package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/karaMuha/go-chirpy/models"
	"github.com/karaMuha/go-chirpy/service"
	"github.com/karaMuha/go-chirpy/state"
)

type RestHandler struct {
	appState     *state.AppState
	service      service.Service
	userService  service.UsersService
	chirpService service.ChirpsService
}

func NewRestHandler(
	appState *state.AppState,
	service service.Service,
	userService service.UsersService,
	chirpService service.ChirpsService,
) RestHandler {
	return RestHandler{
		appState:    appState,
		service:     service,
		userService: userService,
	}
}

func (h *RestHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type: ", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
	w.WriteHeader(200)
}

func (h *RestHandler) HandleViewMetrics(w http.ResponseWriter, r *http.Request) {
	viewCount := h.appState.FileServerHitsCount()
	viewCountHtml := fmt.Sprintf(`
		<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
		</html>
	`, viewCount)

	w.Write([]byte(viewCountHtml))
	w.Header().Add("Content-Type:", "text/html")
	w.WriteHeader(200)
}

func (h *RestHandler) HandleReset(w http.ResponseWriter, r *http.Request) {
	h.appState.ResetFileServerHitsCound()

	if h.appState.Platform != "dev" {
		http.Error(w, "Not allowed outside of dev env", http.StatusForbidden)
		return
	}

	respErr := h.userService.ResetUsers(r.Context())
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}
}

func (h *RestHandler) HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	decorder := json.NewDecoder(r.Body)
	chirp := models.Chirp{}
	err := decorder.Decode(&chirp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, respErr := h.service.ValidateChirp(chirp)
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}

	respJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respJson)
}

type CreateUserDto struct {
	Email string `json:"email"`
}

func (h *RestHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := CreateUserDto{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, respErr := h.userService.CreateUser(r.Context(), data.Email)
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}

	respJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(respJson)
}

type CreateChirpsDto struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

func (h *RestHandler) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := CreateChirpsDto{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chrip, respErr := h.chirpService.CreateChrip(r.Context(), data.Body, data.UserID)
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}

	respJson, err := json.Marshal(chrip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(respJson)
}
