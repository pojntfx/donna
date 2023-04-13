package psql

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pojntfx/networkmate/internal/migrations"
	migrate "github.com/rubenv/sql-migrate"
)

//go:generate sqlboiler psql -o ../../internal/models -c ../../sqlboiler.yaml
//go:generate go-bindata -pkg migrations -o ../../internal/migrations/migrations.go ../migrations

type RootPersister struct {
	db *sql.DB
}

func NewRootPersister() *RootPersister {
	return &RootPersister{}
}

func (p *RootPersister) Open(dbURL string) error {
	// Connect to the DB
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	// Configure the db
	db.SetMaxOpenConns(1) // Prevent "database locked" errors

	// Run migrations
	if _, err := migrate.Exec(db, "postgres", migrate.AssetMigrationSource{
		Asset:    migrations.Asset,
		AssetDir: migrations.AssetDir,
		Dir:      "../migrations",
	}, migrate.Up); err != nil {
		return err
	}

	p.db = db

	return nil
}
