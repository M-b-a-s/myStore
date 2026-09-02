package products

import (
	"context"
	repo "github/M-b-a-s/myStore/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	GetProductByID(ctx context.Context, id int64) (repo.Product, error)
	CreateProduct(ctx context.Context, input repo.CreateProductParams) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	// Do something with the retrieved products, e.g., format them or perform additional processing
	_ = products
	return products, nil
}

func (s *svc) GetProductByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *svc) CreateProduct(ctx context.Context, input repo.CreateProductParams) (repo.Product, error) {
	return s.repo.CreateProduct(ctx, input)
}
