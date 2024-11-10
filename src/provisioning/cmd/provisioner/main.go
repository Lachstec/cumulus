package main

import (
	"fmt"
	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.LoadConfig("config.json")
	s, err := sqlx.Open("pgx", cfg.Db.ConnectionURI())
	if err != nil {
		panic(err)
	}
	mig := db.NewMigrator(s)

	err = mig.Migrate()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database schema has been created")
}
