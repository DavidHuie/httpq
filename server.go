package httpq

import (
	"log"
	"net/http"
)

type Server struct {
	httpq *Httpq
}

func (s *Server) Push(w http.ResponseWriter, r *http.Request) {
	if err := s.httpq.PushRequest(r); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) Pop(w http.ResponseWriter, r *http.Request) {
	requestBytes, err := s.httpq.PopRequestBytes()
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(requestBytes); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
