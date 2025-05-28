package db_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"workout-tracker/pkg/db"
)

func TestNew_DefaultEnv_Success(t *testing.T) {
	t.Setenv("DB_HOST", "centerbeam.proxy.rlwy.net")
	t.Setenv("DB_PORT", "35231")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "gesnsUZGmhNORJdxLZIBDXLBFSnGZnNJ")
	t.Setenv("DB_NAME", "railway")

	logger := zaptest.NewLogger(t).Sugar()
	database, err := db.New(logger)

	require.NoError(t, err)
	require.NotNil(t, database)
	require.NotNil(t, database.Pool)

	database.Pool.Close()
	require.NoError(t, err)
}

func TestNew_InvalidHost_ShouldFail(t *testing.T) {
	t.Setenv("DB_HOST", "invalid_host")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("DB_NAME", "workout_tracker")

	logger := zaptest.NewLogger(t).Sugar()
	database, err := db.New(logger)

	require.Error(t, err)
	require.Nil(t, database)
}
