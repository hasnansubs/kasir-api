package repositories

import "github.com/jackc/pgx/v5"

type ProductRepository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{db: db}
}
