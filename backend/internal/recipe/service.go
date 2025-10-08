package recipe

import (
	"context"
	"errors"
)

var ErrRepositoryUnavailable = errors.New("recipe repository not configured")

// Service exposes higher-level recipe operations shared across transports.
type Service interface {
	List(ctx context.Context, limit, offset int) ([]*Recipe, error)
	Get(ctx context.Context, id string) (*Recipe, error)
	Create(ctx context.Context, name string) (*Recipe, error)
	Rename(ctx context.Context, id, newName string) (*Recipe, error)
	SoftDelete(ctx context.Context, id string) error
	Ping(ctx context.Context) error
}

type service struct {
	repo Repository
}

// NewService constructs a recipe Service from the given Repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, limit, offset int) ([]*Recipe, error) {
	if err := s.ensureRepo(); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *service) Get(ctx context.Context, id string) (*Recipe, error) {
	if err := s.ensureRepo(); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *service) Create(ctx context.Context, name string) (*Recipe, error) {
	if err := s.ensureRepo(); err != nil {
		return nil, err
	}

	id, err := s.repo.Create(ctx, name)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) Rename(ctx context.Context, id, newName string) (*Recipe, error) {
	if err := s.ensureRepo(); err != nil {
		return nil, err
	}

	if err := s.repo.Rename(ctx, id, newName); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) SoftDelete(ctx context.Context, id string) error {
	if err := s.ensureRepo(); err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, id)
}

func (s *service) Ping(ctx context.Context) error {
	if err := s.ensureRepo(); err != nil {
		return err
	}
	return s.repo.Ping(ctx)
}

func (s *service) ensureRepo() error {
	if s.repo == nil {
		return ErrRepositoryUnavailable
	}
	return nil
}
