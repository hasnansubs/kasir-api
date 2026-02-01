package models

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	CategoryId int    `json:"categoryId"`
	Stock      int    `json:"stock"`
}

type GetProductResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Category string `json:"category"`
	Stock    int    `json:"stock"`
}
