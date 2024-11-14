package config

import "fmt"

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

// ConnectionURI prints a PostgreSQL connection string from a DbConfig.
func (db *DbConfig) ConnectionURI() string {
	s := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", db.User, db.Password, db.Host, db.Port)
	return s
}
