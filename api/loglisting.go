package api

import (
	"encoding/json"
	"net/http"

	logHandler "github.com/NguyenIconAI/logspullerservice/pkg"
)

type GetLogFilesResponse struct {
	LogFiles []string `json:"files"`
}

// handleGetLogFiles returns a list of log files in a directory.
// GET /logs
// Response: 200 OK
//
//	{
//	  "files": [
//	    "/var/log/syslog",
//	    "/var/log/messages",
//	    "/var/log/nginx/access.log",
//	    "/var/log/nginx/error.log"
//	  ]
//	}
func (s *Server) handleGetLogFiles(w http.ResponseWriter, r *http.Request) {
	logFiles, err := logHandler.ListLogFiles("/var/log")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := GetLogFilesResponse{
		LogFiles: logFiles,
	}
	json.NewEncoder(w).Encode(response)
}
