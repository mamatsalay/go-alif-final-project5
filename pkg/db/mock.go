package db

import (
	"context"
)

type mockPool struct{}

func (m *mockPool) Ping(ctx context.Context) error {
	return nil
}
