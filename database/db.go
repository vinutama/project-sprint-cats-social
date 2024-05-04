package database

import (
	cfg "cats-social/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbName            = cfg.EnvConfigs.DbName
	dbHost            = cfg.EnvConfigs.DbHost
	dbPass            = cfg.EnvConfigs.DbPassword
	dbUser            = cfg.EnvConfigs.DbUser
	dbPort            = cfg.EnvConfigs.DbPort
	dbParams            = cfg.EnvConfigs.DbParams
	dbTimeout         = 30 * time.Second
	dbMaxConnLifeTime = 60 * time.Minute
	dbMaxConnIdleTime = 5 * time.Minute
	dbMaxConn         = int32(3000)
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
