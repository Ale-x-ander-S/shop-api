package repository

import (
	"context"
	"shop-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.CreateProductRequest) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx,
		`INSERT INTO products (name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		product.Name, product.Description, product.Price, product.Stock, time.Now(), time.Now(),
	).Scan(&id)
	return id, err
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	product := &models.Product{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, description, price, stock, created_at, updated_at
		FROM products WHERE id = $1`,
		id,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return product, err
}

func (r *ProductRepository) Update(ctx context.Context, id int64, product *models.UpdateProductRequest) error {
	_, err := r.db.Exec(ctx,
		`UPDATE products 
		SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			price = COALESCE($3, price),
			stock = COALESCE($4, stock),
			updated_at = $5
		WHERE id = $6`,
		product.Name, product.Description, product.Price, product.Stock, time.Now(), id,
	)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	return err
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, description, price, stock, created_at, updated_at
		FROM products`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
