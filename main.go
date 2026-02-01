package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// CONFIG
type Config struct {
	PORT        string `mapstructure:"APP_PORT"`
	DB_USER     string `mapstructure:"DB_USERNAME"`
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_NAME     string `mapstructure:"DB_NAME"`
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

func getCategoryById(id int, w http.ResponseWriter, _ *http.Request) {
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
	viper.AutomaticEnv()
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		PORT:        viper.GetString("APP_PORT"),
		DB_USER:     viper.GetString("DB_USERNAME"),
		DB_PASSWORD: viper.GetString("DB_PASSWORD"),
		DB_HOST:     viper.GetString("DB_HOST"),
		DB_PORT:     viper.GetString("DB_PORT"),
		DB_NAME:     viper.GetString("DB_NAME"),
	}

	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	db, err := database.InitDB(connStr)
	if err != nil {
		log.Fatal("failed to initialize database:", err.Error())
	}
	defer db.Close(context.Background())

	productRepo := repositories.NewRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	fmt.Println("starting kasir-api server....")
	fmt.Printf("server running on %v", config.PORT)

	// /api/products/{id}
	http.HandleFunc("/api/products/", productHandler.HandleProductsDetail)
	// /api/products
	http.HandleFunc("/api/products", productHandler.HandleProducts)

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

	err = http.ListenAndServe(config.PORT, nil)
	if err != nil {
		fmt.Print("failed to start the server")
	}

}
