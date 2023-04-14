package persisters

//go:generate sqlc -f ../../sqlc.yaml generate

import (
	"context"
	"database/sql"

	"github.com/pojntfx/networkmate/internal/migrations"
	"github.com/pojntfx/networkmate/internal/models"
	"github.com/pressly/goose/v3"
)

type Persister struct {
	dbaddr  string
	queries *models.Queries
}

func NewPersister(dbaddr string) *Persister {
	return &Persister{
		dbaddr: dbaddr,
	}
}

func (p *Persister) Init() error {
	db, err := sql.Open("postgres", p.dbaddr)
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		return err
	}

	p.queries = models.New(db)

	return nil
}

func (p *Persister) GetContacts(ctx context.Context) ([]models.Contact, error) {
	return p.queries.GetContacts(ctx)
}
