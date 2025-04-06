package repository

import (
	"context"
	"errors"
	"shop-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProductNotFound = errors.New("product not found")

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	var product models.Product
	err := r.db.QueryRow(ctx,
		`INSERT INTO products (name, description, price, stock) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, name, description, price, stock, created_at, updated_at`,
		req.Name, req.Description, req.Price, req.Stock,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	var product models.Product
	err := r.db.QueryRow(ctx,
		`SELECT id, name, description, price, stock, created_at, updated_at 
		 FROM products 
		 WHERE id = $1`,
		id,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int64, req *models.UpdateProductRequest) error {
	result, err := r.db.Exec(ctx,
		`UPDATE products 
		 SET name = $1, description = $2, price = $3, stock = $4, updated_at = NOW() 
		 WHERE id = $5`,
		req.Name, req.Description, req.Price, req.Stock, id,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	return err
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, description, price, stock, created_at, updated_at 
		 FROM products 
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
