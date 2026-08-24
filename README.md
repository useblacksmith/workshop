# Confab — the CI Speedrun workshop

A small conference schedule & live-voting app with a deliberately honest CI
pipeline. It exists so you can *feel* slow CI, then fix it yourself:
migrate this repo's workflow from GitHub-hosted runners to
[Blacksmith](https://blacksmith.sh) in three copy-paste edits, then let
Codesmith fix the config debt, and watch the same pipeline go from
**~9 minutes → ~3 minutes → ~2 minutes → ~1 minute**.

**👉 Doing the workshop? Everything you need is in [MIGRATION.md](MIGRATION.md).**

## What's inside

| Path | What | CI job |
|---|---|---|
| `api/` | Go API (Gin + pgx) serving the schedule and votes | `integration` (tested against a real Postgres), `docker` |
| `stats/` | Rust service (Axum) computing vote analytics | `rust` |
| `web/` | React + TypeScript frontend (pnpm, Vite, vitest) | `web` |
| `e2e/` | Playwright browser tests | `e2e` |
| `api/Dockerfile` | Multi-platform (amd64 + arm64) image build | `docker` |

One workflow, five parallel jobs — a miniature version of a real team's CI.
On GitHub-hosted runners the arm64 Docker leg runs under QEMU emulation,
which is why the `docker` job is the long pole. That's not sabotage; it's
how most teams build multi-arch images today.

## Run it locally

```bash
# API (in-memory store)
cd api && go run .

# Web (proxies /api to :3000)
pnpm install && pnpm --filter @confab/web dev

# Stats service
cd stats && cargo run

# Tests
cd api && go test ./...
cd stats && cargo test
pnpm test          # web unit tests
pnpm e2e           # Playwright (needs: pnpm --filter @confab/e2e exec playwright install chromium)
```

## The migration steps

Each step has a matching branch if you fall behind:

1. [`step-1-runners`](../../tree/step-1-runners) — swap runner labels
2. [`step-2-stickydisk`](../../tree/step-2-stickydisk) — put the cargo build on a sticky disk
3. [`step-3-docker`](../../tree/step-3-docker) — native arm64 builds, no more QEMU
4. [`step-4-optimized`](../../tree/step-4-optimized) — what `@codesmith` fixes: the caches,
   browser bloat, and concurrency this workflow "forgot" (like most real ones do)

See [MIGRATION.md](MIGRATION.md) for the full guide, including creating a
throwaway GitHub org so nothing touches your company's setup — and how to
clean everything up afterwards.
