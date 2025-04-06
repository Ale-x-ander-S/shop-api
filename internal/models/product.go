package models

import (
	"time"
)

type Product struct {
	ID          int64     `json:"id" redis:"id"`
	Name        string    `json:"name" redis:"name"`
	Description string    `json:"description" redis:"description"`
	Price       float64   `json:"price" redis:"price"`
	Stock       int       `json:"stock" redis:"stock"`
	Category    string    `json:"category" redis:"category"`
	ImageURL    string    `json:"image_url" redis:"image_url"`
	CreatedAt   time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" redis:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	Category    string  `json:"category" binding:"required"`
	ImageURL    string  `json:"image_url"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Stock       int     `json:"stock,omitempty" binding:"min=0"`
	Category    string  `json:"category,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
}
