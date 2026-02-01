package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func InitDB(connectionString string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
