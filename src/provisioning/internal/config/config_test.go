package config

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_HOST", "sample.com")
	os.Setenv("DB_PORT", "1337")
	os.Setenv("DB_USER", "sample_user")
	os.Setenv("DB_PASS", "sample_pass")
	os.Setenv("AUTH0_URL", "https://auth0.com/test")
	os.Setenv("AUTH0_AUDIENCE", "https://mc-hosting.zip")
	os.Setenv("AUTH0_SECRET", "sample_secret")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("AUTH0_URL")
		os.Unsetenv("AUTH0_AUDIENCE")
		os.Unsetenv("AUTH0_SECRET")
	}()

	cfg := LoadConfig()

	assert.Equal(t, "sample.com", cfg.Db.Host)
	assert.Equal(t, "1337", cfg.Db.Port)
	assert.Equal(t, "sample_user", cfg.Db.User)
	assert.Equal(t, "sample_pass", cfg.Db.Password)
	assert.Equal(t, url.URL{Scheme: "https", Host: "auth0.com", Path: "/test"}, cfg.Auth0.AuthURL)
	assert.Equal(t, "https://mc-hosting.zip", cfg.Auth0.Audience)
	assert.Equal(t, "sample_secret", cfg.Auth0.Secret)
}

func TestFallbackValues(t *testing.T) {
	cfg := LoadConfig()

	assert.Equal(t, "localhost", cfg.Db.Host)
	assert.Equal(t, "5432", cfg.Db.Port)
	assert.Equal(t, "postgres", cfg.Db.User)
	assert.Equal(t, "postgres", cfg.Db.Password)
	assert.Equal(t, url.URL{Scheme: "http", Host: "localhost"}, cfg.Auth0.AuthURL)
	assert.Equal(t, "http://localhost", cfg.Auth0.Audience)
	assert.Equal(t, "secret", cfg.Auth0.Secret)
}
