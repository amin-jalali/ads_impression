package tests

import (
	"learning/cmd/server"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test case when the environment variable is set
	key := "TEST_ENV"
	value := "some_value"
	os.Setenv(key, value) // Set the environment variable

	err, result := server.GetEnv(key)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if result != value {
		t.Errorf("Expected %s, but got %s", value, result)
	}

	// Test case when the environment variable is not set
	os.Unsetenv(key) // Unset the environment variable

	err, result = server.GetEnv(key)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
	if result != "" {
		t.Errorf("Expected an empty string, but got %s", result)
	}
}
