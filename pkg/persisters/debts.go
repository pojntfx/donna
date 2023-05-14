package persisters

import (
	"context"

	"github.com/pojntfx/donna/pkg/models"
)

func (p *Persister) CreateDebt(
	ctx context.Context,

	amount float64,
	currency string,

	contactID int32,
	namespace string,
) (int32, error) {
	return p.queries.CreateDebt(ctx, models.CreateDebtParams{
		ID:        contactID,
		Namespace: namespace,
		Amount:    amount,
		Currency:  currency,
	})
}

func (p *Persister) GetDebts(
	ctx context.Context,

	contactID int32,
	namespace string,
) ([]models.GetDebtsRow, error) {
	return p.queries.GetDebts(ctx, models.GetDebtsParams{
		ID:        contactID,
		Namespace: namespace,
	})
}
