package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB holds a test Postgres container and connection pool.
type TestDB struct {
	Pool      *pgxpool.Pool
	Container testcontainers.Container
}

// TestRedis holds a test Redis container and client.
type TestRedis struct {
	Client    *redis.Client
	Container testcontainers.Container
}

// SetupPostgres spins up an ephemeral Postgres container, runs migrations, and returns a pool.
func SetupPostgres(t *testing.T) *TestDB {
	t.Helper()
	ctx := context.Background()

	// Read migration file
	migrationPath, err := findMigrationFile()
	if err != nil {
		t.Fatalf("failed to find migration file: %v", err)
	}

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("volo_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithInitScripts(migrationPath),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
		container.Terminate(ctx)
	})

	return &TestDB{Pool: pool, Container: container}
}

// SetupRedis spins up an ephemeral Redis container and returns a client.
func SetupRedis(t *testing.T) *TestRedis {
	t.Helper()
	ctx := context.Background()

	container, err := tcredis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(15*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("failed to get redis endpoint: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	t.Cleanup(func() {
		client.Close()
		container.Terminate(ctx)
	})

	return &TestRedis{Client: client, Container: container}
}

// findMigrationFile locates the migration SQL file relative to the test.
func findMigrationFile() (string, error) {
	// Try common relative paths from test directories
	paths := []string{
		"../../migrations/001_initial.sql",
		"../../../migrations/001_initial.sql",
		"migrations/001_initial.sql",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	// Try absolute from project root
	cwd, _ := os.Getwd()
	return "", fmt.Errorf("migration file not found (cwd: %s)", cwd)
}
