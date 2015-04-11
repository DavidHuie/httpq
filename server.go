package httpq

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	httpq *Httpq
}

func NewServer(httpq *Httpq) *Server {
	return &Server{httpq}
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
}

func (s *Server) Size(w http.ResponseWriter, r *http.Request) {
	size, err := s.httpq.Size()
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sizeMap := map[string]uint64{"size": size}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(sizeMap); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
