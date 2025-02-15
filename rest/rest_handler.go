package rest

import (
	"net/http"
	"strconv"

	"github.com/karaMuha/go-chirpy/state"
)

type RestHandler struct {
	appState *state.AppState
}

func NewRestHandler(appState *state.AppState) RestHandler {
	return RestHandler{
		appState: appState,
	}
}

func (h *RestHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type: ", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
	w.WriteHeader(200)
}

func (h *RestHandler) HandleViewMetrics(w http.ResponseWriter, r *http.Request) {
	viewCount := h.appState.FileServerHitsCount()
	result := "Hits: " + strconv.Itoa(int(viewCount))

	w.Write([]byte(result))
	w.WriteHeader(200)
}

func (h *RestHandler) HandleResetViewCount(w http.ResponseWriter, r *http.Request) {
	h.appState.ResetFileServerHitsCound()
}
