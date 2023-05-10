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
