//go:build integration

package api

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

// Integration test to check if server health endpoint returns the correct status
func IntegTest_HealthCheck(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/health")
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

	var actualResponse HealthCheckResponse
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
		t.FailNow()
	}

	expectedResponse := HealthCheckResponse{
		Status: "OK",
	}

	if actualResponse != expectedResponse {
		t.Errorf("unexpected response body: got %v, want %v", actualResponse, expectedResponse)
	}
}
