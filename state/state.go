package state

import (
	"log"
	"net/http"
	"sync/atomic"
)

type AppState struct {
	fileserverHits atomic.Int32
	Platform       string
	Secret         string
	PolkaKey       string
}

func NewAppState(platform string) *AppState {
	return &AppState{
		Platform: platform,
	}
}

func (s *AppState) IncMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.fileserverHits.Add(1)
		log.Println(s.fileserverHits.Load())
		next.ServeHTTP(w, r)
	})
}

func (s *AppState) FileServerHitsCount() int32 {
	return s.fileserverHits.Load()
}

func (s *AppState) ResetFileServerHitsCound() {
	s.fileserverHits.Store(0)
}
