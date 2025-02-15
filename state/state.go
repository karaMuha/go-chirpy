package state

import (
	"log"
	"net/http"
	"sync/atomic"
)

type AppState struct {
	fileserverHits atomic.Int32
}

func NewAppState() *AppState {
	return &AppState{}
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
