package persisters

import (
	"context"

	"github.com/pojntfx/senbara/senbara-forms/pkg/models"
)

type AllUserData struct {
	JournalEntries []models.JournalEntry
	Contacts       []models.Contact
	Debts          []models.GetDebtsForNamespaceRow
	Activities     []models.GetActivitiesForNamespaceRow
}

func (p *Persister) GetAllUserData(ctx context.Context, namespace string) (AllUserData, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return AllUserData{}, err
	}
	defer tx.Rollback()

	allUserData := AllUserData{}

	qtx := p.queries.WithTx(tx)

	allUserData.JournalEntries, err = qtx.GetJournalEntries(ctx, namespace)
	if err != nil {
		return AllUserData{}, err
	}

	allUserData.Contacts, err = qtx.GetContacts(ctx, namespace)
	if err != nil {
		return AllUserData{}, err
	}

	allUserData.Debts, err = qtx.GetDebtsForNamespace(ctx, namespace)
	if err != nil {
		return AllUserData{}, err
	}

	allUserData.Activities, err = qtx.GetActivitiesForNamespace(ctx, namespace)
	if err != nil {
		return AllUserData{}, err
	}

	if err := tx.Commit(); err != nil {
		return AllUserData{}, err
	}

	return allUserData, nil
}
