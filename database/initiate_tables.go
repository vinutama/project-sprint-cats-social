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
		`
		CREATE TABLE IF NOT EXISTS cats (
			id VARCHAR(100) NOT NULL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			race VARCHAR(100) NOT NULL,
			sex VARCHAR(10) NOT NULL,
			age_in_month INT NOT NULL,
			user_id VARCHAR(100) NOT NULL,
			description VARCHAR(255) NOT NULL,
			image_urls TEXT NOT NULL,
			has_matched BOOL NOT NULL DEFAULT FALSE,
			is_deleted BOOL NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS matches (
			id VARCHAR(100) NOT NULL PRIMARY KEY,
			message VARCHAR(255) NOT NULL,
			status VARCHAR(100) NOT NULL DEFAULT 'requested',
			cat_issuer_id VARCHAR(100) NOT NULL,
			cat_receiver_id VARCHAR(100) NOT NULL,
			FOREIGN KEY (cat_issuer_id) REFERENCES cats(id) ON DELETE NO ACTION,
			FOREIGN KEY (cat_receiver_id) REFERENCES cats(id) ON DELETE NO ACTION,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
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
