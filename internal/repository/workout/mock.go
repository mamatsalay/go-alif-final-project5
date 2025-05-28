package workout

import (
	"context"
	"fmt"
	"log"

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
	callArgs := append([]interface{}{ctx, sql}, args...)
	called := m.Called(callArgs...)
	return called.Get(0).(pgx.Row)
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	callArgs := append([]interface{}{ctx, sql}, args...)
	called := m.Called(callArgs...)
	return called.Get(0).(pgx.Rows), called.Error(1)
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	callArgs := append([]interface{}{ctx, sql}, args...)
	called := m.Called(callArgs...)
	return called.Get(0).(pgconn.CommandTag), called.Error(1)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	fields, ok := args.Get(0).([]pgconn.FieldDescription)
	if !ok {
		return nil
	}
	return fields
}

func (m *MockRow) Close() {
	m.Called()
}

func (m *MockRow) CommandTag() pgconn.CommandTag {
	args := m.Called()
	values, ok := args.Get(0).(pgconn.CommandTag)
	if !ok {
		log.Fatal("invalid type for pgconn.CommandTag")
		return values
	}
	return values
}

func (m *MockRow) Conn() *pgx.Conn {
	args := m.Called()
	conn, ok := args.Get(0).(*pgx.Conn)
	if !ok {
		return nil
	}
	return conn
}

func (m *MockRow) Err() error {
	args := m.Called()
	return fmt.Errorf("mock error: %w", args.Error(0))
}

func (m *MockRow) RawValues() [][]byte {
	args := m.Called()
	values, ok := args.Get(0).([][]byte)
	if !ok {
		return nil
	}
	return values
}

func (m *MockRow) Values() ([]interface{}, error) {
	args := m.Called()

	raw := args.Get(0)
	values, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected []interface{} but got %T", raw)
	}

	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock error: %w", err)
	}

	return values, nil
}

func (m *MockRow) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("error scanning row: %w", err)
	}
	return nil
}
