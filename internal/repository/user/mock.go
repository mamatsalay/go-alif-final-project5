package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

const (
	lengthOffset = 2
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	callArgs := make([]interface{}, 0, lengthOffset+len(args))
	callArgs = append(callArgs, ctx, sql)
	callArgs = append(callArgs, args...)
	called := m.Called(callArgs...)

	row, ok := called.Get(0).(pgx.Row)
	if !ok {
		return nil
	}
	return row
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	callArgs := make([]interface{}, 0, lengthOffset+len(args))
	callArgs = append(callArgs, ctx, sql)
	callArgs = append(callArgs, args...)
	called := m.Called(callArgs...)

	rows, ok := called.Get(0).(pgx.Rows)
	if !ok {
		return nil, errors.New("mock: invalid pgx.Rows return")
	}
	return rows, fmt.Errorf("error invalid pgx.Rows return: %w", called.Error(1))
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	callArgs := make([]interface{}, 0, lengthOffset+len(args))
	callArgs = append(callArgs, ctx, sql)
	callArgs = append(callArgs, args...)
	called := m.Called(callArgs...)

	tag, ok := called.Get(0).(pgconn.CommandTag)
	if !ok {
		return pgconn.CommandTag{}, errors.New("mock: invalid CommandTag return")
	}
	return tag, fmt.Errorf("error invalid CommandTag return: %w", called.Error(1))
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return fmt.Errorf("error scanning row: %w", args.Error(0))
}
