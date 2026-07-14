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

	if err := ensureDatabaseSchema(ctx, pool); err != nil {
		t.Fatalf("failed to ensure database schema: %v", err)
	}

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

func TestLoadProjects(t *testing.T) {
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

	if err := ensureDatabaseSchema(ctx, pool); err != nil {
		t.Fatalf("failed to ensure database schema: %v", err)
	}

	if _, err := pool.Exec(ctx, `
INSERT INTO "Project" (name, type)
SELECT 'Example Project', 'internal'
WHERE NOT EXISTS (SELECT 1 FROM "Project")
`); err != nil {
		t.Fatalf("failed to seed project row: %v", err)
	}

	projects, err := loadProjects(ctx, pool)
	if err != nil {
		t.Fatalf("loadProjects returned error: %v", err)
	}
	if len(projects) == 0 {
		t.Fatal("expected at least one project from the database")
	}
	if projects[0].Name == "" || projects[0].Type == "" {
		t.Fatalf("expected the first project to have name and type, got %#v", projects[0])
	}
}
