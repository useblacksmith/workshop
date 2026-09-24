# Blacksmith-Demo: the CI Speedrun workshop

A small conference schedule & voting app with deliberately typical CI. It's
the fallback repo for the hands-on Blacksmith workshop: if you can't migrate
your own repo, you migrate this one: same three steps, guaranteed to work.

**👉 Doing the workshop? Everything is in [MIGRATION.md](MIGRATION.md).**

## Copy-paste kit

Everything you'd otherwise type by hand during the workshop, in one place.
Copy from here instead of transcribing from the screen.

**Links**

| What | Link |
|---|---|
| Sign in to Blacksmith (use this exact link; the `ref` unlocks workshop access, no card) | <https://app.blacksmith.sh/?ref=wearedevelopers-2026> |
| Already had a Blacksmith account? Sign in with this link instead (works while logged in) | <https://dashboardbackend.blacksmith.sh/login/github?redirect=https%3A%2F%2Fapp.blacksmith.sh%2F%3Fref%3Dwearedevelopers-2026> |
| Create a new GitHub org (pick the Free plan) | <https://github.com/account/organizations/new> |
| This repo (hit "Use this template") | <https://github.com/useblacksmith/workshop> |
| Sticky disks guide | <https://docs.blacksmith.sh/blacksmith-caching/dependencies-sticky-disks> |
| Docs home | <https://docs.blacksmith.sh> |

**Names and labels**

```text
<your github username>-blacksmith-workshop     # throwaway org name
blacksmith-4vcpu-ubuntu-2404                   # default runner label
```

**Codesmith prompts** (comment on your open migration PR, or paste into a
Codesmith chat in the dashboard)

Step 2, sticky disks only:

```text
@codesmith mount sticky disks for the expensive paths in this workflow.
```

Step 3, the rest of the caching stack (same PR as Step 2):

```text
@codesmith find any remaining CI optimizations in this workflow: swap
checkout to useblacksmith/checkout, enable Docker layer caching with the
Blacksmith build actions, and add any caches I am missing.
```

**Baseline run** (if your workflow has no Run workflow button):

```bash
git commit --allow-empty -m "baseline" && git push
```

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
ways real workflows are (no dependency caches, browsers re-downloaded every run,
cargo compiling from scratch); that's not sabotage, it's realism, and
fixing it is the workshop.

## The three steps

1. Swap runner labels to Blacksmith
2. Cargo build on a sticky disk + persistent Docker layer cache
3. The PR `@codesmith` opens: the full caching stack (sticky disks, checkout, Docker layers)

Later, with a week of run history: ask Codesmith to `/rightsize` your longest workflow.

Fell behind? Appendix B of [MIGRATION.md](MIGRATION.md) has the finished
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
