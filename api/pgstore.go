package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgStore struct {
	pool *pgxpool.Pool
}

func newPGStore(ctx context.Context, databaseURL string) (*pgStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	s := &pgStore{pool: pool}
	if err := s.migrate(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *pgStore) migrate(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS talks (
			id      TEXT PRIMARY KEY,
			title   TEXT NOT NULL,
			speaker TEXT NOT NULL,
			room    TEXT NOT NULL,
			slot    TEXT NOT NULL,
			votes   INT  NOT NULL DEFAULT 0
		)`)
	if err != nil {
		return err
	}
	for _, t := range seedTalks() {
		_, err := s.pool.Exec(ctx,
			`INSERT INTO talks (id, title, speaker, room, slot) VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (id) DO NOTHING`,
			t.ID, t.Title, t.Speaker, t.Room, t.Slot)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *pgStore) Talks(ctx context.Context) ([]Talk, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, title, speaker, room, slot, votes FROM talks ORDER BY slot`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var talks []Talk
	for rows.Next() {
		var t Talk
		if err := rows.Scan(&t.ID, &t.Title, &t.Speaker, &t.Room, &t.Slot, &t.Votes); err != nil {
			return nil, err
		}
		talks = append(talks, t)
	}
	return talks, rows.Err()
}

func (s *pgStore) Vote(ctx context.Context, id string) (int, error) {
	var votes int
	err := s.pool.QueryRow(ctx,
		`UPDATE talks SET votes = votes + 1 WHERE id = $1 RETURNING votes`, id).Scan(&votes)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errNotFound
	}
	if err != nil {
		return 0, err
	}
	return votes, nil
}
