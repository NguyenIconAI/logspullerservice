package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/NguyenIconAI/logspullerservice/constants"
	logHandler "github.com/NguyenIconAI/logspullerservice/pkg"
)

// handleReadLogFile handles the read log file endpoint
// @Summary Read log file
// @Description Reads the last N lines from a log file and returns them as a JSON array
// @Tags logs
// @Produce json
// @Param file query string true "Log file"
// @Param n query int true "Number of lines"
// @Param filter query string false "Filter"
// @Success 200 {array} string
// @Failure 400 {string} string "Invalid 'n' parameter"
// @Failure 400 {string} string "Missing 'file' parameter"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/log [get]
func (s *Server) handleReadLogFile(w http.ResponseWriter, r *http.Request) {
	numOfLinesStr := r.URL.Query().Get("n")
	numOfLines, err := strconv.Atoi(numOfLinesStr)
	if err != nil || numOfLines <= 0 || numOfLines > 1000 {
		http.Error(w, "Invalid 'n' parameter", http.StatusBadRequest)
		return
	}

	logFile := r.URL.Query().Get("file")
	if logFile == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		return
	}

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(filter)
		if !isAlphanumeric {
			http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
			return
		}
	}

	// Sanitize the file path to avoid security issues
	logFilePath, err := sanitizeFilePath(logFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read last N lines
	response, err := logHandler.ReadLastNLines(logFilePath, numOfLines, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(response) == 0 {
		http.Error(w, "No lines found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Sanitize the file path to avoid directory traversal attacks
func sanitizeFilePath(filePath string) (string, error) {
	cleanPath := cleanFilePath(filePath)

	// Ensure the cleaned path does not contain any directory traversal sequences
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("invalid file path")
	}

	fullPath := constants.LogFilesPath + "/" + cleanPath

	// Ensure the full path starts with the base directory
	if !strings.HasPrefix(fullPath, constants.LogFilesPath) {
		return "", fmt.Errorf("invalid file path")
	}

	return fullPath, nil
}

// Clean the file path by removing any redundant elements
func cleanFilePath(filePath string) string {
	// Remove leading and trailing slashes
	filePath = strings.Trim(filePath, "/")

	// Split the path into parts
	parts := strings.Split(filePath, "/")

	// Filter out any empty parts or current directory references
	var cleanParts []string
	for _, part := range parts {
		if part != "" && part != "." {
			cleanParts = append(cleanParts, part)
		}
	}

	// Join the cleaned parts back into a single path
	return strings.Join(cleanParts, "/")
}
