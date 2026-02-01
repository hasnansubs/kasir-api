package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetCategories(w, r)
	case http.MethodPost:
		h.AddCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) HandleCategoriesDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	switch r.Method {
	case http.MethodGet:
		h.GetCategory(id, w, r)
	case http.MethodPut:
		h.EditCategory(id, w, r)
	case http.MethodDelete:
		h.DeleteCategory(id, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	Categories, err := h.service.GetCategoriesService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Categories)
}

func (h *CategoryHandler) AddCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	id, err := h.service.AddCategory(newCategory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newCategory.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   newCategory,
	})
}

func (h *CategoryHandler) GetCategory(id int, w http.ResponseWriter, r *http.Request) {
	category, err := h.service.GetCategoryById(id)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			http.Error(w, "Category is not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   category,
	})
}

func (h *CategoryHandler) EditCategory(id int, w http.ResponseWriter, r *http.Request) {
	var updatedCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	updatedCategory.ID = id
	category, err := h.service.EditCategory(updatedCategory)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			http.Error(w, " no category to update", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   category,
	})
}

func (h *CategoryHandler) DeleteCategory(id int, w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteCategory(id)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			http.Error(w, " no category to delete", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
