package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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
		case http.MethodPut:
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

			// didnt found
			http.Error(w, "product is not found", http.StatusNotFound)
		case http.MethodDelete:
			// find the product
			for i, product := range products {
				if id == product.ID {
					products = append(products[:i], products[i+1:]...)
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}

			// didnt found
			http.Error(w, "product is not found", http.StatusNotFound)
		}
	})

	// /api/products
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(products)
			return
		case http.MethodPost:
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
			return
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("failed to start the server")
	}

}
