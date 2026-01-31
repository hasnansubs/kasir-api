package main

import (
	"fmt"
	"net/http"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

func main() {
	fmt.Println("starting kasir-api server....")
	fmt.Println("server running on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("failed to start the server")
	}

}
