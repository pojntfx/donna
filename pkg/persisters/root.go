package persisters

//go:generate sqlc -f ../../sqlc.yaml generate

import (
	"context"
	"database/sql"

	"github.com/pojntfx/donna/pkg/migrations"
	"github.com/pojntfx/donna/pkg/models"
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

func (p *Persister) GetJournalEntries(ctx context.Context, namespace string) ([]models.JournalEntry, error) {
	return p.queries.GetJournalEntries(ctx, namespace)
}

func (p *Persister) CreateJournalEntry(ctx context.Context, title, body string, rating int32, namespace string) (int32, error) {
	return p.queries.CreateJournalEntry(ctx, models.CreateJournalEntryParams{
		Title:     title,
		Body:      body,
		Rating:    rating,
		Namespace: namespace,
	})
}

func (p *Persister) DeleteJournalEntry(ctx context.Context, id int32, namespace string) error {
	return p.queries.DeleteJournalEntry(ctx, models.DeleteJournalEntryParams{
		ID:        id,
		Namespace: namespace,
	})
}

func (p *Persister) GetJournalEntry(ctx context.Context, id int32, namespace string) (models.JournalEntry, error) {
	return p.queries.GetJournalEntry(ctx, models.GetJournalEntryParams{
		ID:        id,
		Namespace: namespace,
	})
}

func (p *Persister) UpdateJournalEntry(ctx context.Context, id int32, title, body string, rating int32, namespace string) error {
	return p.queries.UpdateJournalEntry(ctx, models.UpdateJournalEntryParams{
		ID:        id,
		Namespace: namespace,
		Title:     title,
		Body:      body,
		Rating:    rating,
	})
}

func (p *Persister) GetContacts(ctx context.Context, namespace string) ([]models.Contact, error) {
	return p.queries.GetContacts(ctx, namespace)
}

func (p *Persister) CreateContact(
	ctx context.Context,
	firstName string,
	lastName string,
	nickname string,
	email string,
	pronouns string,
	namespace string,
) (int32, error) {
	return p.queries.CreateContact(ctx, models.CreateContactParams{
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickname,
		Email:     email,
		Pronouns:  pronouns,
		Namespace: namespace,
	})
}
