package persisters

import (
	"context"
	"time"

	"github.com/pojntfx/donna/pkg/models"
)

func (p *Persister) GetPendingTodos(ctx context.Context, namespace string) ([]models.Todo, error) {
	return p.queries.GetPendingTodos(ctx, namespace)
}

func (p *Persister) GetDoneTodos(ctx context.Context, namespace string) ([]models.Todo, error) {
	return p.queries.GetDoneTodos(ctx, namespace)
}

func (p *Persister) CreateTodo(ctx context.Context, name string, deadline time.Time, importance int32, namespace string) (int32, error) {
	return p.queries.CreateTodo(ctx, models.CreateTodoParams{
		Name:       name,
		Deadline:   deadline,
		Importance: importance,
		Namespace:  namespace,
	})
}

func (p *Persister) DeleteTodo(ctx context.Context, id int32, namespace string) error {
	return p.queries.DeleteTodo(ctx, models.DeleteTodoParams{
		ID:        id,
		Namespace: namespace,
	})
}

func (p *Persister) GetTodo(ctx context.Context, id int32, namespace string) (models.Todo, error) {
	return p.queries.GetTodo(ctx, models.GetTodoParams{
		ID:        id,
		Namespace: namespace,
	})
}

func (p *Persister) CloseTodo(ctx context.Context, id int32, namespace string) error {
	return p.queries.CloseTodo(ctx, models.CloseTodoParams{
		ID:        id,
		Namespace: namespace,
	})
}

func (p *Persister) UpdateTodo(ctx context.Context, id int32, name string, deadline time.Time, importance int32, namespace string) error {
	return p.queries.UpdateTodo(ctx, models.UpdateTodoParams{
		ID:         id,
		Name:       name,
		Deadline:   deadline,
		Importance: importance,
		Namespace:  namespace,
	})
}
