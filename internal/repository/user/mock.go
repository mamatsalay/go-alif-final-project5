package user

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	callArgs := make([]interface{}, 0, 2+len(args))
	callArgs = append(callArgs, ctx, sql)
	for _, a := range args {
		callArgs = append(callArgs, a)
	}
	called := m.Called(callArgs...)
	return called.Get(0).(pgx.Row)
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	callArgs := make([]interface{}, 0, 2+len(args))
	callArgs = append(callArgs, ctx, sql)
	for _, a := range args {
		callArgs = append(callArgs, a)
	}
	called := m.Called(callArgs...)
	return called.Get(0).(pgx.Rows), called.Error(1)
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	callArgs := make([]interface{}, 0, 2+len(args))
	callArgs = append(callArgs, ctx, sql)
	for _, a := range args {
		callArgs = append(callArgs, a)
	}
	called := m.Called(callArgs...)
	return called.Get(0).(pgconn.CommandTag), called.Error(1)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}
