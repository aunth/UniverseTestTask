package service

import (
	"context"
	"fmt"
	"time"

	"catalog-product/internal/metrics"
	"catalog-product/internal/model"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*model.Product, error)
}

type MessageBroker interface {
	PublishProductCreated(ctx context.Context, product *model.Product) error
	PublishProductDeleted(ctx context.Context, id uuid.UUID) error
}

type ProductService struct {
	repo   ProductRepository
	broker MessageBroker
}

func NewProductService(repo ProductRepository, broker MessageBroker) *ProductService {
	return &ProductService{
		repo:   repo,
		broker: broker,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string, price float64) (*model.Product, error) {
	product := &model.Product{
		ID:        uuid.New(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now().UTC(),
	}

	err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("service failed to create product: %w", err)
	}

	if s.broker != nil {
		err = s.broker.PublishProductCreated(ctx, product)
		if err != nil {
			fmt.Printf("warning: failed to publish creation event: %v\n", err)
		}
	}

	metrics.ProductsCreated.Inc()

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service failed to delete product: %w", err)
	}

	if s.broker != nil {
		err = s.broker.PublishProductDeleted(ctx, id)
		if err != nil {
			fmt.Printf("warning: failed to publish deletion event: %v\n", err)
		}
	}

	metrics.ProductsDeleted.Inc()

	return nil
}

func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	products, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service failed to retrieve product list: %w", err)
	}

	return products, nil
}
