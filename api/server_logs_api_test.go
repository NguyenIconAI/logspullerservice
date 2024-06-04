//go:build integration

package api

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// Response structure
type LogFilesResponse struct {
	Files []string `json:"files"`
}

// Integration test to check if server returns a list of log files
func IntegTest_GetLogFiles(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/v1/logs")
	if err != nil {
		t.Errorf("error: %v", err)
		t.FailNow()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("error: %v", err)
		t.FailNow()
	}

	var actualResponse LogFilesResponse
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
		t.FailNow()
	}

	expectedFiles := []string{
		"/var/log/syslog",
		"/var/log/auth.log",
		"/var/log/kern.log",
	}

	for _, expectedFile := range expectedFiles {
		if !contains(actualResponse.Files, expectedFile) {
			t.Errorf("expected file %s not found in response", expectedFile)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
