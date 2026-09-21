# Blacksmith-Demo: the CI Speedrun workshop

A small conference schedule & voting app with deliberately typical CI. It's
the fallback repo for the hands-on Blacksmith workshop: if you can't migrate
your own repo, you migrate this one: same four steps, guaranteed to work.

**👉 Doing the workshop? Everything is in [MIGRATION.md](MIGRATION.md).**

## What's inside

| Path | What | CI job |
|---|---|---|
| `api/` | Go API (Gin + pgx) serving the schedule and votes | `integration` (tested against a real Postgres), `docker` |
| `stats/` | Rust service (Axum) computing vote analytics | `rust` |
| `web/` | React + TypeScript frontend (pnpm, Vite, vitest) | `web` |
| `e2e/` | Playwright browser tests | `e2e` |
| `api/Dockerfile` | Container image build | `docker` |

One workflow, five parallel jobs: a miniature of a real team's CI, on
GitHub-hosted runners with stock actions. It's unoptimized in exactly the
ways real workflows are (no dependency caches, every browser installed,
cargo compiling from scratch); that's not sabotage, it's realism, and
fixing it is the workshop.

## The four steps

1. Swap runner labels to Blacksmith
2. Cargo build on a sticky disk + persistent Docker layer cache
3. The PR `@codesmith` opens: caches, Chromium-only, concurrency
4. Right-sizing happens in the Blacksmith dashboard

Fell behind? Appendix C of [MIGRATION.md](MIGRATION.md) has the finished
workflow to paste in.

## Run it locally

```bash
cd api && go run .                       # API (in-memory store)
pnpm install && pnpm --filter @blacksmith-demo/web dev   # web, proxies /api
cd stats && cargo run                    # stats service

# tests
cd api && go test ./...
cd stats && cargo test
pnpm test          # web unit tests
pnpm e2e           # Playwright (first: pnpm --filter @blacksmith-demo/e2e exec playwright install chromium)
```
