package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// PRODUCT
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Indomie Goreng", Price: 3500, Stock: 10},
	{ID: 2, Name: "Vit 1000mL", Price: 3000, Stock: 40},
}

func getProducts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	newId := len(products) + 1
	newProduct.ID = newId
	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   newProduct,
	})
}

func editProduct(id int, w http.ResponseWriter, r *http.Request) {
	// find the product
	for i, product := range products {
		if id == product.ID {
			var updatedProduct Product
			err := json.NewDecoder(r.Body).Decode(&updatedProduct)
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
			}

			updatedProduct.ID = id
			products[i] = updatedProduct

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data":   products[i],
			})
			return
		}
	}

	http.Error(w, "product is not found", http.StatusNotFound)
}

func deleteProduct(id int, w http.ResponseWriter, _ *http.Request) {
	// find the product
	for i, product := range products {
		if id == product.ID {
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "product is not found", http.StatusNotFound)
}

func getProductById(id int, w http.ResponseWriter, r *http.Request) {
	for i, product := range products {
		if id == product.ID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data":   products[i],
			})
			return
		}
	}

	// didnt found
	http.Error(w, "product is not found", http.StatusNotFound)
}

// CATEGORY
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Electronic", Description: "Electronic stuff"},
	{ID: 2, Name: "Groceries", Description: "Grocery struff"},
}

func getCategories(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	newId := len(categories) + 1
	newCategory.ID = newId
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   newCategory,
	})
}

func editCategory(id int, w http.ResponseWriter, r *http.Request) {
	// find the product
	for i, category := range categories {
		if id == category.ID {
			var updatedcategory Category
			err := json.NewDecoder(r.Body).Decode(&updatedcategory)
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
			}

			updatedcategory.ID = id
			categories[i] = updatedcategory

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data":   categories[i],
			})
			return
		}
	}

	http.Error(w, "category is not found", http.StatusNotFound)
}

func deleteCategory(id int, w http.ResponseWriter, _ *http.Request) {
	// find the product
	for i, category := range categories {
		if id == category.ID {
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "category is not found", http.StatusNotFound)
}

func getCategoryById(id int, w http.ResponseWriter, r *http.Request) {
	for i, category := range categories {
		if id == category.ID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data":   categories[i],
			})
			return
		}
	}

	// didnt found
	http.Error(w, "category is not found", http.StatusNotFound)
}

func main() {
	fmt.Println("starting kasir-api server....")
	fmt.Println("server running on 8080")

	// /api/products/{id}
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
		}

		switch r.Method {
		case http.MethodGet:
			getProductById(id, w, r)
		case http.MethodPut:
			editProduct(id, w, r)
		case http.MethodDelete:
			deleteProduct(id, w, r)
		}
	})
	// /api/products
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProducts(w, r)
		case http.MethodPost:
			addProduct(w, r)
		}
	})

	// /api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
		}

		switch r.Method {
		case http.MethodGet:
			getCategoryById(id, w, r)
		case http.MethodPut:
			editCategory(id, w, r)
		case http.MethodDelete:
			deleteCategory(id, w, r)
		}
	})
	// /api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCategories(w, r)
		case http.MethodPost:
			addCategory(w, r)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("failed to start the server")
	}

}
