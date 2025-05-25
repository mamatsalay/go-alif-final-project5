package db

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

const timeOut = time.Second * 5

type DB struct {
	Pool *pgxpool.Pool
}

func New(logger *zap.SugaredLogger) (*DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
		port = "5432"
		user = "postgres"
		password = "postgres"
		dbname = "workout_tracker"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s",
		user, password, net.JoinHostPort(host, port), dbname)
	logger.Infof("Connecting to PostgreSQL...",
		"host", host,
		"port", port,
		"user", user,
		"dbname", dbname)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Errorw("error parsing database config", "error", err)
		return nil, fmt.Errorf("error parsing database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Errorw("error connecting to database", "error", err)
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Errorw("error pinging database", "error", err)
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	logger.Info("PostgreSQL connection established")
	return &DB{Pool: pool}, nil
}
