// +build integration

package api

import (
	"io"
	"net/http"
	"testing"
)


// This is a wrapper function that calls the actual integration test function.
func TestIntegTests(t *testing.T) {
	t.Run("HealthCheck", IntegTest_HealthCheck)
}

// Integration test to check if server is up and running
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

	bodyString := string(body)
	if bodyString != "{\"status\":\"OK\"}\n" {
		t.Errorf("unexpected response body: %v", bodyString)
		t.FailNow()
	}
}
