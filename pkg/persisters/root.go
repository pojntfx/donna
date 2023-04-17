package persisters

//go:generate sqlc -f ../../sqlc.yaml generate

import (
	"context"
	"database/sql"

	"github.com/pojntfx/donna/internal/migrations"
	"github.com/pojntfx/donna/internal/models"
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

func (p *Persister) GetJournalEntries(ctx context.Context) ([]models.JournalEntry, error) {
	return p.queries.GetJournalEntries(ctx)
}

func (p *Persister) CreateJournalEntries(ctx context.Context, title, body string) error {
	return p.queries.CreateJournalEntries(ctx, models.CreateJournalEntriesParams{
		Title: title,
		Body:  body,
	})
}
