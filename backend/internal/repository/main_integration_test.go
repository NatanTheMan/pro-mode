//go:build integration

package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migratePost "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	ctr, err := postgres.Run(
		ctx, "postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(err)
	}

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	// Run your REAL migrations here, not a copy of the schema,
	// so the unique index on LOWER(name) is what gets tested.
	if err := runMigrations(dsn); err != nil {
		panic(err)
	}

	code := m.Run()

	_ = ctr.Terminate(ctx)
	os.Exit(code)
}

func runMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratePost.WithInstance(
		db,
		&migratePost.Config{},
	)
	if err != nil {
		return err
	}

	migrationsPath, err := filepath.Abs("../../migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+filepath.ToSlash(migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil &&
		!errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func resetDB(t *testing.T) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `TRUNCATE categories RESTART IDENTITY CASCADE`); err != nil {
		t.Fatal(err)
	}
}
