package exercise

import (
	"context"

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
	return called.Get(0).(pgx.Row)
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	callArgs := make([]interface{}, 0, lengthOffset+len(args))
	callArgs = append(callArgs, ctx, sql)
	callArgs = append(callArgs, args...)

	called := m.Called(callArgs...)
	return called.Get(0).(pgx.Rows), called.Error(1)
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	callArgs := make([]interface{}, 0, lengthOffset+len(args))
	callArgs = append(callArgs, ctx, sql)
	callArgs = append(callArgs, args...)

	called := m.Called(callArgs...)
	return called.Get(0).(pgconn.CommandTag), called.Error(1)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRow) Close() {
	m.Called()
}

func (m *MockRow) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

func (m *MockRow) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

func (m *MockRow) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRow) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}

func (m *MockRow) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockRow) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}
