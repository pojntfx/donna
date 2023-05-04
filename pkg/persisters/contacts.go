package persisters

import (
	"context"

	"github.com/pojntfx/donna/pkg/models"
)

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

func (p *Persister) GetContact(ctx context.Context, id int32, namespace string) (models.Contact, error) {
	return p.queries.GetContact(ctx, models.GetContactParams{
		ID:        id,
		Namespace: namespace,
	})
}

func (p *Persister) DeleteContact(ctx context.Context, id int32, namespace string) error {
	return p.queries.DeleteContact(ctx, models.DeleteContactParams{
		ID:        id,
		Namespace: namespace,
	})
}
