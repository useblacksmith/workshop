package main

import (
	"context"
	"errors"
	"sync"
)

var errNotFound = errors.New("not found")

// Talk is one schedule entry.
type Talk struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Speaker string `json:"speaker"`
	Room    string `json:"room"`
	Slot    string `json:"slot"`
	Votes   int    `json:"votes"`
}

// Store is the persistence interface; backed by memory or Postgres.
type Store interface {
	Talks(ctx context.Context) ([]Talk, error)
	Vote(ctx context.Context, id string) (int, error)
}

// newStore returns a Postgres-backed store when databaseURL is set,
// otherwise an in-memory store seeded with the demo schedule.
func newStore(ctx context.Context, databaseURL string) (Store, error) {
	if databaseURL != "" {
		return newPGStore(ctx, databaseURL)
	}
	return newMemStore(), nil
}

func seedTalks() []Talk {
	return []Talk{
		{ID: "keynote", Title: "The Feedback Loop Is the Product", Speaker: "Ada Okafor", Room: "Main Stage", Slot: "09:00"},
		{ID: "ci-agents", Title: "Your CI Was Sized for Humans", Speaker: "Jonas Weber", Room: "Track 2", Slot: "10:30"},
		{ID: "monorepo", Title: "Monorepos Without Tears", Speaker: "Priya Nair", Room: "Track 1", Slot: "11:15"},
		{ID: "arm64", Title: "Escaping QEMU: Native Multi-Arch Builds", Speaker: "Sofia Lindqvist", Room: "Track 2", Slot: "13:00"},
		{ID: "rust-perf", Title: "Rust Compile Times: A Support Group", Speaker: "Marco Bianchi", Room: "Track 3", Slot: "14:30"},
		{ID: "postgres", Title: "Postgres Is All You Need", Speaker: "Elif Demir", Room: "Track 1", Slot: "16:00"},
	}
}

type memStore struct {
	mu    sync.Mutex
	talks []Talk
}

func newMemStore() *memStore {
	return &memStore{talks: seedTalks()}
}

func (s *memStore) Talks(ctx context.Context) ([]Talk, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Talk, len(s.talks))
	copy(out, s.talks)
	return out, nil
}

func (s *memStore) Vote(ctx context.Context, id string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.talks {
		if s.talks[i].ID == id {
			s.talks[i].Votes++
			return s.talks[i].Votes, nil
		}
	}
	return 0, errNotFound
}
