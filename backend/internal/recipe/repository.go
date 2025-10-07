package recipe

import "context"

type Repository interface {
	Create(ctx context.Context, name string) (string, error)
	GetByID(ctx context.Context, id string) (*Recipe, error)
	GetByName(ctx context.Context, name string) (*Recipe, error)
	List(ctx context.Context, limit, offset int) ([]*Recipe, error)
	Rename(ctx context.Context, id string, newName string) error
	SoftDelete(ctx context.Context, id string) error
	Ping(ctx context.Context) error
}
