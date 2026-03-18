package repository

import (
	"context"
	"fmt"

	"catalog-product/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (id, name, price, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, product.ID, product.Name, product.Price, product.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save product to DB: %w", err)
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product from DB: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("product with ID %s not found", id)
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	query := `
		SELECT id, name, price, created_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product list: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows from DB: %w", err)
	}

	return products, nil
}
