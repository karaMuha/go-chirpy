package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/karaMuha/go-chirpy/internal/auth"
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
		appState:     appState,
		service:      service,
		userService:  userService,
		chirpService: chirpService,
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *RestHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := CreateUserDto{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, respErr := h.userService.CreateUser(r.Context(), data.Email, hashedPassword)
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
	Body string `json:"body"`
}

func (h *RestHandler) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	token, err := auth.GetBearerToken(headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, h.appState.Secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	data := CreateChirpsDto{}
	err = decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chrip, respErr := h.chirpService.CreateChrip(r.Context(), data.Body, userID.String())
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

func (h *RestHandler) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	sorting := r.URL.Query().Get("sort")
	if sorting == "" {
		sorting = "ASC"
	}
	authorID := r.URL.Query().Get("author_id")

	chirps, respErr := h.chirpService.GetAll(r.Context(), authorID, sorting)
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}

	respJson, err := json.Marshal(chirps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respJson)
}

func (h *RestHandler) HandleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirp, respErr := h.chirpService.GetByID(r.Context(), chirpID)
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}

	respJson, err := json.Marshal(chirp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respJson)
}

type LoginDto struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresIn int    `json:"expires_in_seconds,omitempty"`
}

func (h *RestHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := LoginDto{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expirationDuration := data.ExpiresIn
	if expirationDuration <= 0 || expirationDuration > 3600 {
		expirationDuration = 3600
	}
	user, respErr := h.userService.Login(r.Context(), data.Email, data.Password, expirationDuration)
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
	w.WriteHeader(200)
	w.Write(respJson)
}

func (h *RestHandler) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	token, err := auth.GetBearerToken(headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, h.appState.Secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data CreateUserDto
	err = decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser, respErr := h.userService.UpdateAccount(r.Context(), userID.String(), data.Email, data.Password)
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}

	respJson, err := json.Marshal(updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respJson)
}

func (h *RestHandler) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	token, err := auth.GetBearerToken(headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, h.appState.Secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chirpID := r.PathValue("chirpID")
	respErr := h.chirpService.Delete(r.Context(), userID.String(), chirpID)
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}

	w.WriteHeader(204)
}

type WebhookEvent struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (h *RestHandler) HandleUpgradeToRed(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if key != h.appState.PolkaKey {
		http.Error(w, "Key does not match", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var event WebhookEvent
	err = decoder.Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if event.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	respErr := h.userService.UpgradeToRed(r.Context(), event.Data.UserID)
	if respErr != nil {
		http.Error(w, respErr.Error, respErr.StatusCode)
		return
	}

	w.WriteHeader(204)
}
