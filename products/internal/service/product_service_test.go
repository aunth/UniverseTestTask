package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"catalog-product/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of ProductRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, product *model.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Product), args.Error(1)
}

// MockBroker is a mock of MessageBroker
type MockBroker struct {
	mock.Mock
}

func (m *MockBroker) PublishProductCreated(ctx context.Context, product *model.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockBroker) PublishProductDeleted(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductService_CreateProduct(t *testing.T) {
	repo := new(MockRepository)
	broker := new(MockBroker)
	svc := NewProductService(repo, broker)

	ctx := context.Background()
	name := "Test Product"
	price := 99.99

	repo.On("Create", ctx, mock.AnythingOfType("*model.Product")).Return(nil)
	broker.On("PublishProductCreated", ctx, mock.AnythingOfType("*model.Product")).Return(nil)

	product, err := svc.CreateProduct(ctx, name, price)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, name, product.Name)
	assert.Equal(t, price, product.Price)
	repo.AssertExpectations(t)
}

func TestProductService_DeleteProduct(t *testing.T) {
	repo := new(MockRepository)
	broker := new(MockBroker)
	svc := NewProductService(repo, broker)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)
	broker.On("PublishProductDeleted", ctx, id).Return(nil)

	err := svc.DeleteProduct(ctx, id)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProductService_ListProducts(t *testing.T) {
	repo := new(MockRepository)
	svc := NewProductService(repo, nil)

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedProducts := []*model.Product{
		{ID: uuid.New(), Name: "P1", Price: 10, CreatedAt: time.Now()},
		{ID: uuid.New(), Name: "P2", Price: 20, CreatedAt: time.Now()},
	}

	repo.On("List", ctx, limit, offset).Return(expectedProducts, nil)

	products, err := svc.ListProducts(ctx, limit, offset)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "P1", products[0].Name)
	repo.AssertExpectations(t)
}

func TestProductService_CreateProduct_Error(t *testing.T) {
	repo := new(MockRepository)
	svc := NewProductService(repo, nil)

	ctx := context.Background()
	
	repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

	product, err := svc.CreateProduct(ctx, "Fail", 10)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "db error")
}
