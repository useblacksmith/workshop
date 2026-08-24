package main

import (
	"context"
	"os"
	"testing"
)

// Integration tests against a real Postgres. They run when DATABASE_URL is
// set (in CI: a service container) and skip otherwise, so `go test ./...`
// still passes locally without a database.
func pgTestStore(t *testing.T) *pgStore {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping Postgres integration test")
	}
	s, err := newPGStore(context.Background(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() {
		s.pool.Exec(context.Background(), `DROP TABLE IF EXISTS talks`)
		s.pool.Close()
	})
	return s
}

func TestPGTalksSeeded(t *testing.T) {
	s := pgTestStore(t)

	talks, err := s.Talks(context.Background())
	if err != nil {
		t.Fatalf("talks: %v", err)
	}
	if len(talks) != 6 {
		t.Fatalf("got %d talks, want 6", len(talks))
	}
}

func TestPGVotePersists(t *testing.T) {
	s := pgTestStore(t)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		votes, err := s.Vote(ctx, "rust-perf")
		if err != nil {
			t.Fatalf("vote %d: %v", i, err)
		}
		if votes != i {
			t.Fatalf("got %d votes after %d votes cast", votes, i)
		}
	}

	if _, err := s.Vote(ctx, "nope"); err != errNotFound {
		t.Fatalf("got %v, want errNotFound", err)
	}
}
