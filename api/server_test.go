//go:build integration

package api

import (
	"testing"
)

// This is a wrapper function that calls the actual integration test function.
func TestIntegTests(t *testing.T) {
	t.Run("HealthCheck", IntegTest_HealthCheck)
	t.Run("GetLogFiles", IntegTest_GetLogFiles)
}
