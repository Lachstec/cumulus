package config

// DbConfig contains all information to connect to the application database.
type DbConfig struct {
	// Host where the database can be reached
	Host string
	// Port on which the database listens
	Port string
	// User to use for authentication
	User string
	// Password to use for authentication
	Password string
}
