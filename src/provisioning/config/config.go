package config

import (
	"encoding/json"
	"os"
)

// Config is the main configuration type for the application.
// All settings or credentials the backend needs is to be supplied here.
type Config struct {
	// Db configuration to access the primary database
	Db DbConfig
}

// LoadConfig loads the config from the file at given path.
// File is expected to be in JSON format. When an error occurs,
// this function panics as config is mandatory for the service to work.
func LoadConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		panic("Failed to open config file. aborting.")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("Unexpected error when closing file")
		}
	}(file)

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		panic("Invalid config file format")
	}

	return &config
}
