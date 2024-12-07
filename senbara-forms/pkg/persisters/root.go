package persisters

//go:generate sqlc -f ../../sqlc.yaml generate

import (
	"database/sql"

	"github.com/pojntfx/senbara/senbara-forms/pkg/migrations"
	"github.com/pojntfx/senbara/senbara-forms/pkg/models"
	"github.com/pressly/goose/v3"
)

type Persister struct {
	pgaddr  string
	queries *models.Queries
	db      *sql.DB
}

func NewPersister(pgaddr string) *Persister {
	return &Persister{
		pgaddr: pgaddr,
	}
}

func (p *Persister) Init() error {
	var err error
	p.db, err = sql.Open("postgres", p.pgaddr)
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(p.db, "."); err != nil {
		return err
	}

	p.queries = models.New(p.db)

	return nil
}
