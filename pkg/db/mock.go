package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return fmt.Errorf("%w", args.Error(0))
}

func (m *MockPool) Close() {
	_ = m.Called()
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	called := m.Called(ctx, sql, args)
	row, _ := called.Get(0).(pgx.Row)
	return row
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	called := m.Called(ctx, sql, args)
	rows, _ := called.Get(0).(pgx.Rows)
	return rows, fmt.Errorf("%w", called.Error(1))
}

func (m *MockPool) Begin(ctx context.Context) (pgx.Tx, error) {
	called := m.Called(ctx)
	tx, _ := called.Get(0).(pgx.Tx)
	return tx, fmt.Errorf("%w", called.Error(1))
}

func (m *MockPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	called := m.Called(ctx, sql, arguments)
	tag, _ := called.Get(0).(pgconn.CommandTag)
	return tag, fmt.Errorf("%w", called.Error(1))
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	fd, _ := args.Get(0).([]pgconn.FieldDescription)
	return fd
}

func (m *MockRow) Close() {
	_ = m.Called()
}

func (m *MockRow) CommandTag() pgconn.CommandTag {
	args := m.Called()
	tag, _ := args.Get(0).(pgconn.CommandTag)
	return tag
}

func (m *MockRow) Conn() *pgx.Conn {
	args := m.Called()
	conn, _ := args.Get(0).(*pgx.Conn)
	return conn
}

func (m *MockRow) Err() error {
	args := m.Called()
	return fmt.Errorf("%w", args.Error(0))
}

func (m *MockRow) RawValues() [][]byte {
	args := m.Called()
	vals, _ := args.Get(0).([][]byte)
	return vals
}

func (m *MockRow) Values() ([]interface{}, error) {
	args := m.Called()
	vals, _ := args.Get(0).([]interface{})
	return vals, fmt.Errorf("%w", args.Error(1))
}

func (m *MockRow) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return fmt.Errorf("%w", args.Error(0))
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Conn() *pgx.Conn {
	args := m.Called()
	conn, _ := args.Get(0).(*pgx.Conn)
	return conn
}

func (m *MockTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	called := m.Called(ctx, sql, args)
	tag, _ := called.Get(0).(pgconn.CommandTag)
	return tag, fmt.Errorf("%w", called.Error(1))
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	tx, _ := args.Get(0).(pgx.Tx)
	return tx, fmt.Errorf("%w", args.Error(1))
}

func (m *MockTx) Rollback(ctx context.Context) error {
	return fmt.Errorf("%w", m.Called(ctx).Error(0))
}

func (m *MockTx) Commit(ctx context.Context) error {
	return fmt.Errorf("%w", m.Called(ctx).Error(0))
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	called := m.Called(ctx, sql, args)
	rows, _ := called.Get(0).(pgx.Rows)
	return rows, fmt.Errorf("%w", called.Error(1))
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	called := m.Called(ctx, sql, args)
	row, _ := called.Get(0).(pgx.Row)
	return row
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier,
	columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columnNames, rowSrc)
	n, _ := args.Get(0).(int64)
	return n, fmt.Errorf("%w", args.Error(1))
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	args := m.Called()
	lo, _ := args.Get(0).(pgx.LargeObjects)
	return lo
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	sd, _ := args.Get(0).(*pgconn.StatementDescription)
	return sd, fmt.Errorf("%w", args.Error(1))
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	results, _ := args.Get(0).(pgx.BatchResults)
	return results
}

type MockBatchResults struct {
	mock.Mock
}

func (m *MockBatchResults) Close() error {
	args := m.Called()
	return fmt.Errorf("%w", args.Error(0))
}

func (m *MockBatchResults) Exec() (pgconn.CommandTag, error) {
	args := m.Called()
	tag, _ := args.Get(0).(pgconn.CommandTag)
	return tag, fmt.Errorf("%w", args.Error(1))
}

func (m *MockBatchResults) Query() (pgx.Rows, error) {
	args := m.Called()
	rows, _ := args.Get(0).(pgx.Rows)
	return rows, fmt.Errorf("%w", args.Error(1))
}

func (m *MockBatchResults) QueryRow() pgx.Row {
	args := m.Called()
	row, _ := args.Get(0).(pgx.Row)
	return row
}
