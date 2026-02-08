package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProducts(w, r)
	case http.MethodPost:
		h.AddProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleProductsDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	switch r.Method {
	case http.MethodGet:
		h.GetProduct(id, w, r)
	case http.MethodPut:
		h.EditProduct(id, w, r)
	case http.MethodDelete:
		h.DeleteProduct(id, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	products, err := h.service.GetProductsService(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	id, err := h.service.AddProduct(newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newProduct.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   newProduct,
	})
}

func (h *ProductHandler) GetProduct(id int, w http.ResponseWriter, r *http.Request) {
	product, err := h.service.GetProductById(id)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			http.Error(w, "product is not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   product,
	})
}

func (h *ProductHandler) EditProduct(id int, w http.ResponseWriter, r *http.Request) {
	var updatedProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	updatedProduct.ID = id
	product, err := h.service.EditProduct(updatedProduct)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			http.Error(w, " no product to update", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   product,
	})
}

func (h *ProductHandler) DeleteProduct(id int, w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteProduct(id)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			http.Error(w, " no product to delete", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
