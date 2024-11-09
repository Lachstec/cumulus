package db

import (
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
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

// migrate applies all migrations in the migrations directory.
func (m *Migrator) migrate() error {
	pg, err := postgres.WithInstance(m.DB.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	migrations, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		pg)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
