package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/NguyenIconAI/logspullerservice/utils"
)

// RemoteLogRequest represents the request body for reading a remote log file
// swagger:parameters RemoteLogRequest
type RemoteLogRequest struct {
	File   string     `json:"file" validate:"required"`
	N      int        `json:"n"`
	Filter string     `json:"filter"`
	Hosts  []HostInfo `json:"hosts" validate:"required"`
}

// HostInfo represents the host information
// swagger:parameters HostInfo
type HostInfo struct {
	HostName string `json:"host_name" validate:"required"`
	ApiKey   string `json:"api_key" validate:"required"`
}

// RemoteLogResponse represents the response body for reading a remote log file
// swagger:response RemoteLogResponse
type RemoteLogResponse struct {
	Status     string   `json:"status"`
	StatusCode int      `json:"status_code"`
	Logs       []string `json:"logs"`
}

// @Summary Read remote log file
// @Description Reads the last N lines from a log file on a remote host and returns them as a JSON array
// @Tags logs
// @Accept json
// @Produce json
// @Param request body RemoteLogRequest true "Request body"
// @Success 200 {object} map[string]RemoteLogResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /v1/remotelog [post]
func (s *Server) handleReadRemoteLogFile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request RemoteLogRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// we don't validate the request here because the hosts will validate them

	response := make(map[string]RemoteLogResponse)

	wg := sync.WaitGroup{}
	wg.Add(len(request.Hosts))
	// goroutine to get logs from each host concurrently
	for _, host := range request.Hosts {
		go func(host HostInfo) {
			defer wg.Done()
			response[host.HostName] = getLogFromHost(host, request.File, request.N, request.Filter)
		}(host)
	}
	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Simple rest client to get logs from a remote host
func getLogFromHost(host HostInfo, file string, n int, filter string) (response RemoteLogResponse) {
	defer func() {
		if r := recover(); r != nil {
			response = RemoteLogResponse{
				Status:     "panic occurred",
				StatusCode: http.StatusInternalServerError,
				Logs:       []string{fmt.Sprintf("Recovered from panic: %v", r)},
			}
		}
	}()

	bearer := "Bearer " + host.ApiKey
	url := fmt.Sprintf(host.HostName+"/v1/log?file=%s&n=%d&filter=%s", file, n, filter)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := utils.NewRetryableClient()
	client.Timeout = time.Second * 30
	resp, err := client.Do(req)
	if err != nil {
		return RemoteLogResponse{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RemoteLogResponse{
			Status:     "cannot read response body",
			StatusCode: http.StatusInternalServerError,
		}
	}

	var actualResponse []string
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		return RemoteLogResponse{
			Status:     "cannot unmarshal response",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return RemoteLogResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Logs:       actualResponse,
	}
}
