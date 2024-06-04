//go:build integration

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/NguyenIconAI/logspullerservice/constants"
)

// Response structure
type LogFilesResponse struct {
	Files []string `json:"files"`
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Integration test to check if server returns a list of log files
func IntegTest_GetLogFiles(t *testing.T) {
	bearer := "Bearer " + os.Getenv(constants.ApiKeyEnvVar)
	url := "http://localhost:3000/v1/logs"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
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

func IntegTest_GetLogContent(t *testing.T) {
	bearer := "Bearer " + os.Getenv(constants.ApiKeyEnvVar)
	expectedLength := 10
	url := fmt.Sprintf("http://localhost:3000/v1/log?file=auth.log&n=%d&filter=session", expectedLength)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err :err= client.Do(req)
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

	var actualResponse []string
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
		t.FailNow()
	}

	if len(actualResponse) != expectedLength {
		t.Errorf("expected %d lines, got %d", expectedLength, len(actualResponse))
	}
}
