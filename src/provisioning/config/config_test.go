package config

import (
	"os"
	"testing"
)

func TestConfigLoader(t *testing.T) {
	t.Parallel()

	sampleConfig := `{
		"Db": {
			"Host": "localhost",
			"Port": "5432",
			"User": "postgres",
			"Password": "sampletext"
		}
	}`

	temp, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file for testing: %v", err)
	}
	defer os.Remove(temp.Name())

	if _, err := temp.Write([]byte(sampleConfig)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := temp.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	config := LoadConfig(temp.Name())

	if config.Db.Host != "localhost" {
		t.Errorf("Expected Db.Host to be 'localhost', got '%s'", config.Db.Host)
	}
	if config.Db.Port != "5432" {
		t.Errorf("Expected Db.Port to be '5432', got %s", config.Db.Port)
	}
	if config.Db.User != "postgres" {
		t.Errorf("Expected Db.User to be 'testuser', got '%s'", config.Db.User)
	}
	if config.Db.Password != "sampletext" {
		t.Errorf("Expected Db.Password to be 'sampletext', got '%s'", config.Db.Password)
	}
}

func TestConfigLoader_FileNotFound(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for missing config, but didn't get one")
		}
	}()
	LoadConfig("not_existing.json")
}

func TestConfigLoader_BadJson(t *testing.T) {
	t.Parallel()

	temp, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file for testing: %v", err)
	}
	defer os.Remove(temp.Name())

	if _, err := temp.Write([]byte(`{ invalid json }`)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := temp.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid config, but didn't get one")
		}
	}()
	LoadConfig(temp.Name())
}
