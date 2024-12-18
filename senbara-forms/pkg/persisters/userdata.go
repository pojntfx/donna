package persisters

import (
	"context"

	"github.com/pojntfx/senbara/senbara-forms/pkg/models"
)

func (p *Persister) GetUserData(
	ctx context.Context,

	namespace string,

	onJournalEntry func(journalEntry models.ExportedJournalEntry) error,
	onContact func(contact models.ExportedContact) error,
	onDebt func(debt models.ExportedDebt) error,
	onActivity func(activity models.ExportedActivity) error,
) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := p.queries.WithTx(tx)

	journalEntries, err := qtx.GetJournalEntriesExportForNamespace(ctx, namespace)
	if err != nil {
		return err
	}

	for _, journalEntry := range journalEntries {
		if err := onJournalEntry(models.ExportedJournalEntry{
			ID:        journalEntry.ID,
			Title:     journalEntry.Title,
			Date:      journalEntry.Date,
			Body:      journalEntry.Body,
			Rating:    journalEntry.Rating,
			Namespace: journalEntry.Namespace,
		}); err != nil {
			return err
		}
	}

	contacts, err := qtx.GetContactsExportForNamespace(ctx, namespace)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		if err := onContact(models.ExportedContact{
			ID:        contact.ID,
			FirstName: contact.FirstName,
			LastName:  contact.LastName,
			Nickname:  contact.Nickname,
			Email:     contact.Email,
			Pronouns:  contact.Pronouns,
			Namespace: contact.Namespace,
			Birthday:  contact.Birthday,
			Address:   contact.Address,
			Notes:     contact.Notes,
		}); err != nil {
			return err
		}
	}

	debts, err := qtx.GetDebtsExportForNamespace(ctx, namespace)
	if err != nil {
		return err
	}

	for _, debt := range debts {
		if err := onDebt(models.ExportedDebt{
			ID:          debt.ID,
			Amount:      debt.Amount,
			Currency:    debt.Currency,
			Description: debt.Description,
			ContactID:   debt.ContactID,
		}); err != nil {
			return err
		}
	}

	activities, err := qtx.GetActivitiesExportForNamespace(ctx, namespace)
	if err != nil {
		return err
	}

	for _, activity := range activities {
		if err := onActivity(models.ExportedActivity{
			ID:          activity.ID,
			Name:        activity.Name,
			Date:        activity.Date,
			Description: activity.Description,
			ContactID:   activity.ContactID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (p *Persister) DeleteUserData(ctx context.Context, namespace string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := p.queries.WithTx(tx)

	if err := qtx.DeleteActivitiesForNamespace(ctx, namespace); err != nil {
		return err
	}

	if err := qtx.DeleteDebtsForNamespace(ctx, namespace); err != nil {
		return err
	}

	if err := qtx.DeleteContactsForNamespace(ctx, namespace); err != nil {
		return err
	}

	if err := qtx.DeleteJournalEntriesForNamespace(ctx, namespace); err != nil {
		return err
	}

	return tx.Commit()
}
