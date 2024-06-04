package api

import (
	"encoding/json"
	"net/http"

	logHandler "github.com/NguyenIconAI/logspullerservice/pkg"
)

type GetLogFilesResponse struct {
	LogFiles []string `json:"files"`
}

// handleGetLogFiles handles the get log files endpoint
// @Summary Get log files
// @Description Returns a list of log files in a directory
// @Tags logs
// @Produce json
// @Success 200 {object} GetLogFilesResponse
// @Failure 500 {string} string "Internal server error"
// @Router /v1/logs [get]
func (s *Server) handleGetLogFiles(w http.ResponseWriter, r *http.Request) {
	logFiles, err := logHandler.ListLogFiles("/var/log")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(logFiles) == 0 {
		http.Error(w, "No lines found", http.StatusNotFound)
		return
	}

	response := GetLogFilesResponse{
		LogFiles: logFiles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
