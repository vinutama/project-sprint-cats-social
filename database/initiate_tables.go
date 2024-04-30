package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitiateTables(dbPool *pgxpool.Pool) error {
	// Define table creation queries
	queries := []string{
		`
        CREATE TABLE IF NOT EXISTS users (
            id VARCHAR(100) PRIMARY KEY NOT NULL,
            name TEXT NOT NULL,
			email VARCHAR(50) UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
        `,
		// Add more table creation queries here if needed
	}

	// Execute table creation queries
	for _, query := range queries {
		_, err := dbPool.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	return nil
}
