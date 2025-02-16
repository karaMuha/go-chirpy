package rest

import (
	"fmt"
	"net/http"

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

func (h *RestHandler) HandleResetViewCount(w http.ResponseWriter, r *http.Request) {
	h.appState.ResetFileServerHitsCound()
}
