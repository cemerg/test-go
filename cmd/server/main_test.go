package main

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestLoadUsers(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = getDatabaseDSN()
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	users, err := loadUsers(ctx, pool)
	if err != nil {
		t.Fatalf("loadUsers returned error: %v", err)
	}
	if len(users) == 0 {
		t.Fatal("expected at least one user from the database")
	}
	if users[0].Name == "" {
		t.Fatalf("expected the first user to have a name, got %#v", users[0])
	}
}
