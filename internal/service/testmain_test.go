package service_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/adamaso/wallet-service/internal/infrastructure/eventstore"
	"github.com/adamaso/wallet-service/internal/infrastructure/projection"
	"github.com/adamaso/wallet-service/internal/service"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("wallet_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		panic("start postgres container: " + err.Error())
	}
	defer container.Terminate(ctx) //nolint:errcheck

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("get connection string: " + err.Error())
	}

	if err := runMigrations(connStr); err != nil {
		panic("run migrations: " + err.Error())
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		panic("create pool: " + err.Error())
	}
	defer testPool.Close()

	os.Exit(m.Run())
}

// newTestService returns a Service wired to the test database with a clean slate.
func newTestService(t *testing.T) *service.Service {
	t.Helper()
	truncate(t)
	repo := eventstore.NewWalletRepository(testPool, projection.NewWalletProjector())
	store := projection.NewWalletViewRepository(testPool)
	return service.NewService(repo, store)
}

func truncate(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE events, wallet_views")
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}
}

func runMigrations(connStr string) error {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	goose.SetDialect("postgres") //nolint:errcheck
	return goose.Up(db, "../../migrations")
}
