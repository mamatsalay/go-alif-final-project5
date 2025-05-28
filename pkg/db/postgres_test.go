package db

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestMockPool_Ping_Close(t *testing.T) {
	mp := new(MockPool)
	ctx := t.Context()
	mp.On("Ping", ctx).Return(errors.New("ping error"))
	assert.EqualError(t, mp.Ping(ctx), "ping error")

	mp.On("Close").Return()
	mp.Close()
	mp.AssertExpectations(t)
}

func TestMockRow_Methods(t *testing.T) {
	r := new(MockRow)
	// FieldDescriptions
	desc := []pgconn.FieldDescription{{Name: "c"}}
	r.On("FieldDescriptions").Return(desc)
	assert.True(t, reflect.DeepEqual(r.FieldDescriptions(), desc))

	// Close
	r.On("Close").Return()
	r.Close()

	// CommandTag
	tag := pgconn.NewCommandTag("x")
	r.On("CommandTag").Return(tag)
	assert.Equal(t, tag, r.CommandTag())

	// Conn
	conn := &pgx.Conn{}
	r.On("Conn").Return(conn)
	assert.Same(t, conn, r.Conn())

	// Err
	r.On("Err").Return(errors.New("errr"))
	assert.EqualError(t, r.Err(), "errr")

	// RawValues
	vals := [][]byte{{1}, {2}}
	r.On("RawValues").Return(vals)
	assert.True(t, reflect.DeepEqual(r.RawValues(), vals))

	// Values
	vals2 := []interface{}{1, "x"}
	r.On("Values").Return(vals2, nil)
	v2, err := r.Values()
	assert.Equal(t, vals2, v2)
	assert.NoError(t, err)

	// Next
	r.On("Next").Return(true)
	assert.True(t, r.Next())

	// Scan
	errScan := errors.New("scan error")
	r.On("Scan", "d1", "d2").Return(nil)
	assert.NoError(t, r.Scan("d1", "d2"))

	r.On("Scan", "d").Return(errScan)
	err = r.Scan("d")
	assert.Equal(t, errScan, err)

	r.AssertExpectations(t)
}

func TestMockTx_Methods(t *testing.T) {
	tx := new(MockTx)
	ctx := t.Context()

	// Conn
	conn := &pgx.Conn{}
	tx.On("Conn").Return(conn)
	assert.Same(t, conn, tx.Conn())

	// Exec
	tag := pgconn.NewCommandTag("e")
	tx.On("Exec", ctx, "sql", 1).Return(tag, nil)
	tOut, err := tx.Exec(ctx, "sql", 1)
	assert.Equal(t, tag, tOut)
	assert.NoError(t, err)

	// Begin
	var dummyTx pgx.Tx = new(MockTx)
	tx.On("Begin", ctx).Return(dummyTx, nil)
	b, err := tx.Begin(ctx)
	assert.Same(t, dummyTx, b)
	assert.NoError(t, err)

	// Rollback
	errRb := errors.New("rb")
	tx.On("Rollback", ctx).Return(errRb)
	err = tx.Rollback(ctx)
	assert.Equal(t, errRb, err)

	// Commit
	errCm := errors.New("cm")
	tx.On("Commit", ctx).Return(errCm)
	err = tx.Commit(ctx)
	assert.Equal(t, errCm, err)

	// Query
	var dummyRows pgx.Rows = new(MockRow)
	dummyRows.(*MockRow).On("Close").Return(nil)
	tx.On("Query", ctx, "qs", 2).Return(dummyRows, nil)
	rows, err := tx.Query(ctx, "qs", 2)
	assert.Same(t, dummyRows, rows)
	assert.NoError(t, err)
	defer rows.Close()

	// QueryRow
	var dummyRow pgx.Row = new(MockRow)
	tx.On("QueryRow", ctx, "qr", 3).Return(dummyRow)
	r := tx.QueryRow(ctx, "qr", 3)
	assert.Same(t, dummyRow, r)

	// Prepare
	desc := &pgconn.StatementDescription{}
	tx.On("Prepare", ctx, "n", "s").Return(desc, nil)
	pd, err := tx.Prepare(ctx, "n", "s")
	assert.Same(t, desc, pd)
	assert.NoError(t, err)

	// SendBatch
	br := new(MockBatchResults)
	tx.On("SendBatch", ctx, mock.Anything).Return(br)
	res := tx.SendBatch(ctx, nil)
	assert.Same(t, br, res)

	tx.AssertExpectations(t)
}

func TestMockBatchResults_Methods(t *testing.T) {
	br := new(MockBatchResults)

	// Close
	br.On("Close").Return(nil)
	err := br.Close()
	assert.NoError(t, err)

	// Exec
	tag := pgconn.NewCommandTag("bt")
	br.On("Exec").Return(tag, nil)
	tm, err := br.Exec()
	assert.Equal(t, tag, tm)
	assert.NoError(t, err)

	// Query
	var dummyRows pgx.Rows = new(MockRow)
	dummyRows.(*MockRow).On("Close").Return(nil)
	br.On("Query").Return(dummyRows, errors.New("qerr"))
	r, err := br.Query()
	assert.Same(t, dummyRows, r)
	assert.EqualError(t, err, "qerr")
	defer r.Close()

	// QueryRow
	var dummyRow pgx.Row = new(MockRow)
	br.On("QueryRow").Return(dummyRow)
	rw := br.QueryRow()
	assert.Same(t, dummyRow, rw)

	br.AssertExpectations(t)
}
