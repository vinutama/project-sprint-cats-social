package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

var (
	dbName            = viper.GetString("DB_NAME")
	dbHost            = viper.GetString("DB_HOST")
	dbPass            = viper.GetString("DB_PASSWORD")
	dbUser            = viper.GetString("DB_USERNAME")
	dbPort            = viper.GetString("DB_PORT")
	dbParams          = viper.GetString("DB_PARAMS")
	dbTimeout         = 30 * time.Second
	dbMaxConnLifeTime = 2 * time.Minute
	dbMaxConnIdleTime = 5 * time.Second
	dbMaxConn         = int32(100)
	dbMinConn         = int32(0)
)

func GetConnPool() *pgxpool.Pool {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dbUser, dbPass, dbHost, dbPort, dbName, dbParams)
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Error when parsing config DB URL: ", err)
	}

	cfg.MaxConnLifetime = dbMaxConnLifeTime
	cfg.MaxConnIdleTime = dbMaxConnIdleTime
	cfg.MaxConns = dbMaxConn
	cfg.MinConns = dbMinConn
	cfg.ConnConfig.ConnectTimeout = dbTimeout

	dbPool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal("Error when creating Database Pool Context: ", err)
	}

	return dbPool
}
