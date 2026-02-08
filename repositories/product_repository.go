package repositories

import (
	"context"
	"errors"
	"kasir-api/models"

	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAllProducts(nameFilter string) (products []models.GetProductResponse, err error) {
	args := []any{}
	query := "SELECT p.id, p.name, p.price, p.stock, c.name FROM products p JOIN categories c on p.category_id = c.id"
	if nameFilter != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+nameFilter+"%")

	}
	rows, err := repo.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.GetProductResponse
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repo *ProductRepository) AddProduct(newProductRequest models.Product) (id int, err error) {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES (@name, @price, @stock, @categoryId) RETURNING id"
	args := pgx.NamedArgs{
		"name":       newProductRequest.Name,
		"price":      newProductRequest.Price,
		"stock":      newProductRequest.Stock,
		"categoryId": newProductRequest.CategoryId,
	}

	row := repo.db.QueryRow(context.Background(), query, args)
	err = row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (repo *ProductRepository) GetProductById(id int) (product models.GetProductResponse, err error) {
	query := "SELECT p.id, p.name, p.price, p.stock, c.name FROM products p JOIN categories c on p.category_id = c.id WHERE p.id=$1"
	row := repo.db.QueryRow(context.Background(), query, id)
	err = row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category)
	if err != nil {
		if err == pgx.ErrNoRows {
			return product, errors.New("NOT_FOUND")
		}
		return product, err
	}

	return product, nil
}

func (repo *ProductRepository) EditProduct(newProduct models.Product) (product models.Product, err error) {
	query := "UPDATE products SET name=@name, price=@price, stock=@stock, category_id=@categoryId WHERE id=@id RETURNING id, name, price, stock"
	args := pgx.NamedArgs{
		"name":  newProduct.Name,
		"price": newProduct.Price,
		"stock": newProduct.Stock,
		"id":    newProduct.ID,
	}

	row := repo.db.QueryRow(context.Background(), query, args)

	err = row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		if err == pgx.ErrNoRows {
			return product, errors.New("NOT_FOUND")
		}
		return product, err
	}

	return product, nil
}

func (repo *ProductRepository) DeleteProduct(id int) (err error) {
	query := "DELETE from products WHERE id=$1"

	commandTag, err := repo.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("NOT_FOUND")
	}

	return nil
}
