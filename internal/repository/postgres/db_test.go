package postgres_test

import (
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/app"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	// err := godotenv.Load(".env")
	err := config.MustLoad("../../../.env")
	require.NoError(t, err)

	dns := config.Cfg.DatabaseTest.DSN()
	dir := config.Cfg.DatabaseTest.MigrationDir

	db, err := sqlx.Open("pgx", dns)
	require.NoError(t, err)

	err = app.MigrateUp(db, dir)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = app.MigrateDown(db, dir)
		require.NoError(t, err)
		require.NoError(t, db.Close())
	})

	return db
}
