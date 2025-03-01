package tests

import (
	"learning/cmd/server"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV"
	value := "some_value"
	err := os.Setenv(key, value)
	if err != nil {
		return
	}

	result, err := server.GetEnv(key)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	err = os.Unsetenv(key)
	if err != nil {
		return
	}

	result, err = server.GetEnv(key)

	if result != "" {
		t.Errorf("Expected an empty string, but got %s", result)
	}
}
