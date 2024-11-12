package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_HOST", "sample.com")
	os.Setenv("DB_PORT", "1337")
	os.Setenv("DB_USER", "sample_user")
	os.Setenv("DB_PASS", "sample_pass")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
	}()

	cfg := LoadConfig()

	assert.Equal(t, "sample.com", cfg.Db.Host)
	assert.Equal(t, "1337", cfg.Db.Port)
	assert.Equal(t, "sample_user", cfg.Db.User)
	assert.Equal(t, "sample_pass", cfg.Db.Password)
}

func TestFallbackValues(t *testing.T) {
	cfg := LoadConfig()

	assert.Equal(t, "localhost", cfg.Db.Host)
	assert.Equal(t, "5432", cfg.Db.Port)
	assert.Equal(t, "postgres", cfg.Db.User)
	assert.Equal(t, "postgres", cfg.Db.Password)
}
