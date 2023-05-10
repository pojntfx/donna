package persisters

import (
	"context"

	"github.com/pojntfx/donna/pkg/models"
)

// TODO: Use transaction and check if contact belongs to namespace before trying to continue
func (p *Persister) CreateDebt(ctx context.Context, amount float64, currency string, contactID int32, namespace string) (int32, error) {
	return p.queries.CreateDebt(ctx, models.CreateDebtParams{
		Amount:    amount,
		Currency:  currency,
		ContactID: contactID,
	})
}
