package config

// LoggingConfig controls hpw the application Logging behaves
type LoggingConfig struct {
	// Environment controls where logs get sent to. When it is "dev", logs get printed on stdout. On "prod",
	// logs get sent to an external service.
	Environment string
}
