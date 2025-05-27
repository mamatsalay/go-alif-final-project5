package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var newPoolFunc = pgxpool.NewWithConfig

type pinger interface {
	Ping(context.Context) error
}

type mockPool struct{}

func (m *mockPool) Ping(ctx context.Context) error {
	return nil
}
