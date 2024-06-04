package api

import (
	"encoding/json"
	"net/http"

	"github.com/NguyenIconAI/logspullerservice/middleware"
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
	http.Handle("GET /health", middleware.LogMiddleware(http.HandlerFunc(s.handleHealthCheck)))
	http.Handle("GET /v1/logs", middleware.LogMiddleware(middleware.AuthMiddleware(http.HandlerFunc(s.handleGetLogFiles))))
	http.Handle("GET /v1/log", middleware.LogMiddleware(middleware.AuthMiddleware(http.HandlerFunc(s.handleReadLogFile))))
	http.Handle("POST /v1/remotelog", middleware.LogMiddleware(middleware.AuthMiddleware(http.HandlerFunc(s.handleReadRemoteLogFile))))
	// TODO: Adding logging middleware
	return http.ListenAndServe(s.port, nil)
}

// handleHealthCheck handles the health check endpoint
// @Summary Health check endpoint
// @Description Returns the status of the server
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
