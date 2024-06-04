package api

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	port string
}

// Create new server
func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

// Start a server instance
func (s *Server) Start() error {
	http.HandleFunc("/health", s.handleHealthCheck)
	http.HandleFunc("/logs", s.handleGetLogFiles)
	// TODO: Adding authentication and logging middleware
	return http.ListenAndServe(s.port, nil)
}

// Return OK
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(status)
}
