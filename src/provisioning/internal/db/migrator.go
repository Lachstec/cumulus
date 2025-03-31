package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib" // Importing pgx for the driver
	"github.com/jmoiron/sqlx"
)

// Migrator ensures the expected database schema is present by
// applying all sql scripts located in migrations. Migrations need to follow
// the naming scheme {version}_{name}.{up/down}.sql. The up variant should
// create all resources needed by the migration and the down variant should
// undo every change made by the migration.
type Migrator struct {
	*sqlx.DB
}

// NewMigrator creates a new Migrator that uses the connection
// specified by db.
func NewMigrator(db *sqlx.DB) *Migrator {
	return &Migrator{db}
}

// Migrate applies all migrations in the directory dir to the database.
func (m *Migrator) Migrate(dir string) error {
	// Using pgx with the "postgres" driver from golang-migrate
	pg, err := postgres.WithInstance(m.DB.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("file://%s", dir)

	// Ensure that we're using the correct driver
	migrations, err := migrate.NewWithDatabaseInstance(
		path,       // Path to migration files
		"postgres", // Use "postgres" as the driver name
		pg)
	if err != nil {
		return err
	}

	// Apply migrations
	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
