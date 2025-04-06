package handlers

import (
	"encoding/json"
	"net/http"
	"shop-api/internal/models"
	"shop-api/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	service  *service.ProductService
	validate *validator.Validate
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service:  service,
		validate: validator.New(),
	}
}

// @Summary Создать новый продукт
// @Description Создает новый продукт в магазине
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequest true "Данные продукта"
// @Success 201 {object} map[string]int64
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateProduct(r.Context(), &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// @Summary Обновить продукт
// @Description Обновляет информацию о продукте
// @Tags products
// @Accept json
// @Param id path int true "ID продукта"
// @Param product body models.UpdateProductRequest true "Данные для обновления"
// @Success 204
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProduct(r.Context(), id, &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Удалить продукт
// @Description Удаляет продукт по его ID
// @Tags products
// @Param id path int true "ID продукта"
// @Success 204
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Получить все продукты
// @Description Возвращает список всех продуктов
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {string} string
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
