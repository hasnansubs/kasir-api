package main

import (
	"context"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"

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

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	fmt.Println("starting kasir-api server....")
	fmt.Printf("server running on %v", config.PORT)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to kasir-api server!"))
	})

	// /api/products/{id}
	// /api/products
	http.HandleFunc("/api/products/", productHandler.HandleProductsDetail)
	http.HandleFunc("/api/products", productHandler.HandleProducts)

	// /api/categories/{id}
	// /api/categories
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoriesDetail)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)

	err = http.ListenAndServe(config.PORT, nil)
	if err != nil {
		fmt.Print("failed to start the server")
	}

}
