package repository

import (
	"context"
	"errors"
	"shop-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProductNotFound = errors.New("product not found")

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

// PostgresProductRepository реализует интерфейс ProductRepository
type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Create(ctx context.Context, product *models.Product) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO products (name, description, price, stock, category) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id`,
		product.Name, product.Description, product.Price, product.Stock, product.Category).
		Scan(&product.ID)
	return err
}

func (r *PostgresProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	err := r.db.QueryRow(ctx,
		`SELECT id, name, description, price, stock, category, created_at, updated_at 
		 FROM products 
		 WHERE id = $1`,
		id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return &product, nil
}

func (r *PostgresProductRepository) Update(ctx context.Context, product *models.Product) error {
	result, err := r.db.Exec(ctx,
		`UPDATE products 
		 SET name = $1, description = $2, price = $3, stock = $4, category = $5, updated_at = NOW()
		 WHERE id = $6`,
		product.Name, product.Description, product.Price, product.Stock, product.Category, product.ID)
	if err != nil {
		return err
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *PostgresProductRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *PostgresProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
	// Искусственная задержка для демонстрации
	time.Sleep(2 * time.Second)

	rows, err := r.db.Query(ctx,
		`SELECT id, name, description, price, stock, category, created_at, updated_at 
		 FROM products 
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
