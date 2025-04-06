package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"errors"
	"shop-api/internal/models"
	"shop-api/internal/service"

	"github.com/go-chi/chi/v5"
)

var ErrProductNotFound = errors.New("product not found")

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetProducts godoc
// @Summary Получить все продукты
// @Description Возвращает список всех продуктов
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {string} string
// @Router /products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	// Добавляем заголовки для отслеживания кэша
	if h.service.IsFromCache() {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products", http.StatusInternalServerError)
		return
	}
}

// CreateProduct godoc
// @Summary Создать новый продукт
// @Description Создает новый продукт в магазине
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequest true "Данные продукта"
// @Success 201 {object} map[string]int
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.service.CreateProduct(r.Context(), &req)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

// GetProduct godoc
// @Summary Получить продукт по ID
// @Description Возвращает информацию о продукте по его ID
// @Tags products
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} models.Product
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		if err == ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary Обновить продукт
// @Description Обновляет информацию о продукте
// @Tags products
// @Accept json
// @Param id path int true "ID продукта"
// @Param product body models.UpdateProductRequest true "Данные для обновления"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProduct(r.Context(), id, &req); err != nil {
		if err == ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteProduct godoc
// @Summary Удалить продукт
// @Description Удаляет продукт по его ID
// @Tags products
// @Param id path int true "ID продукта"
// @Success 204 "No Content"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		if err == ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
