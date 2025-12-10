package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool
var queriesFS embed.FS

func ConnectPostgres() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalln("Error connecting to database", err)
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	DB = pool
	log.Println("Connected to PostgreSQL successfully!")

	if err := createTables(pool); err != nil {
		log.Fatalln("Error creating tables", err)
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return pool, nil
}

func createTables(pool *pgxpool.Pool) error {
	ctx := context.Background()
	sqlInitTables, _ := LoadQuery("init_tables.sql")

	_, err := pool.Exec(ctx, sqlInitTables)
	if err != nil {
		return err
	}

	log.Println("Tables created successfully!")
	return nil
}

func LoadQuery(name string) (string, error) {
	path := "./internal/db/queries/" + name

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading file", err)
		return "", fmt.Errorf("error reading file: %v", err)
	}

	log.Println("Query read successfully!")

	return string(data), nil
}
