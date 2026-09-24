# The CI Speedrun: workshop guide

Three steps, all in your browser, no local setup: migrate CI to
[Blacksmith](https://blacksmith.sh), put the expensive parts on persistent
disks, then let an agent finish the caching stack.

You'll do this either on **your own repo** (best: you leave with real CI
migrated) or on **our demo repo** (guaranteed to work for everyone).

---

## Step 0: Setup (~10 min, do this first)

1. Get on the wifi, sign in at [github.com](https://github.com).
2. **Get an org you own.** The Blacksmith app installs on organizations, and
   installing it needs **owner** rights (or an owner nearby who'll approve
   your install request fast); personal side-project orgs count. No org you
   own? Create a throwaway: github.com → avatar → Settings → Organizations →
   **New organization** → Free plan. Name it
   `<github username>-blacksmith-workshop`. Skip inviting members.
   ([Appendix A](#appendix-a-cleanup) deletes all of this in one minute.)
3. **Choose your repo:**

**Option A: your own repo** (pick this if you can):
- A repo with real CI that you're **allowed to experiment on**, in the org
  you own. Nothing has to merge today; everything lands as PRs you control.

**Option B: our demo repo** (pick this otherwise):
- On this repo: **Use this template → Create a new repository.**
  Owner: your new org. Visibility: **Public** (free unlimited Actions
  minutes). It's a small real app (Go API + Postgres, Rust service,
  TypeScript frontend, Playwright) with deliberately typical CI.

4. **Sign in to Blacksmith** at
   [app.blacksmith.sh/?ref=wearedevelopers-2026](https://app.blacksmith.sh/?ref=wearedevelopers-2026)
   with GitHub. Use this exact link (or the QR on screen); the `ref` unlocks
   workshop access with no card required. **Already signed in to Blacksmith
   from before?** Use this
   [force sign-in link](https://dashboardbackend.blacksmith.sh/login/github?redirect=https%3A%2F%2Fapp.blacksmith.sh%2F%3Fref%3Dwearedevelopers-2026)
   instead: it re-runs sign-in and attaches the workshop access even while
   you're logged in. Installing the app is Step 1.

**Both options, run your baseline now:** if the repo just landed in your org
(templated or brought over), creating it already triggered the first run: see
the Actions tab and click the newest run (tab empty? your workflow has no
`push` trigger; use the next path). On a repo that already lived there:
**Actions → CI → Run workflow**, or push an empty commit / open a PR. When it
finishes, write down the **per-job durations** (not wall clock): your
before-numbers.

## Step 1: Migrate the runners

1. Install the Blacksmith GitHub App from [app.blacksmith.sh](https://app.blacksmith.sh),
   **scoped to just your chosen repo**; the app only sees what you select.
2. Open the **migration wizard** in the Blacksmith dashboard (the popup you
   saw after installing, or the top of the left panel) and walk its five
   steps:
   1. **Repository:** pick your repo.
   2. **Runners:** pick a size per job; the **4vcpu default is a good pick**
      for today.
   3. **Docker builds:** make sure Docker build caching is **off**.
   4. **Git checkout:** make sure checkout caching is **off**. (Both are
      Step 3: the agent adds them, so you can read the diff.)
   5. **Pull request:** hit **Generate pull request**.
3. The PR run is your Step 1 number. **Everything today stays on this PR**;
   Steps 2 and 3 stack onto it, and nothing needs to merge.

Bonus, zero config: every existing `actions/cache` / `setup-node` cache is
now served **colocated** with the runner, ~4x faster transfers.

> Fell behind? [Appendix B](#appendix-b-the-finished-workflow) has the finished workflow; paste it in on a new branch and open a PR.

## Step 2: Sticky disks

A **sticky disk** is a persistent NVMe volume that mounts into your runner in
seconds, exactly as the last run left it. Ask the agent to add them: comment
on your open PR:

```
@codesmith mount sticky disks for the expensive paths in this workflow.
```

Disks only for now; the rest of the caching stack is Step 3. Here is what it
will add for the demo repo's `rust` job (one step per disk, a key and a path):

```yaml
      - name: Mount sticky disk
        uses: useblacksmith/stickydisk@v1
        with:
          key: ${{ github.repository }}-build-cache
          path: ./<expensive-directory>
```

**Recipe sheet: what to mount**

| Ecosystem | Path to persist |
|---|---|
| Rust | `target/` |
| Turborepo / Nx | `.turbo/` / `.nx/cache` |
| Gradle | `~/.gradle/caches` |
| Go | `~/.cache/go-build` |
| Cypress / Playwright browsers | `~/.cache/Cypress` / `~/.cache/ms-playwright` |
| Docker layers | don't mount; swap to `useblacksmith/setup-docker-builder@v2` (with a `cache-key`) + `useblacksmith/build-push-action@v2`; the layer cache persists automatically |

The disk pays off on the **second** run: the first one only fills it. When
the run goes green, hit **Re-run all jobs** (top right of the run page) and
write down the warm time: that's your Step 2 number.

> Fell behind? [Appendix B](#appendix-b-the-finished-workflow) has the finished workflow.

## Step 3: Let Codesmith configure the rest

Comment on the **same PR** as Step 2, so everything stacks into one
reviewable PR:

```
@codesmith find any remaining CI optimizations in this workflow: swap
checkout to useblacksmith/checkout, enable Docker layer caching with the
Blacksmith build actions, and add any caches I am missing.
```

Name what you want: the agent does exactly what you ask and nothing more. It
pushes commits to the PR; that run fills the new caches, so let it finish,
then **Re-run all jobs** once more and write down the warm time: your final
number. The PR can stay open; merge it later (or never) on your own terms.

Want to check the agent's work, or do it by hand instead? Across Steps 2 and
3 its commits add, all included in [Appendix B](#appendix-b-the-finished-workflow):

1. A `concurrency` block so superseded runs cancel themselves.
2. `cache: pnpm` on both `actions/setup-node` steps (lockfile-keyed).
3. An `actions/cache` step for the cargo registry (`~/.cargo/registry` + `~/.cargo/git`).
4. Playwright: the browser directory goes on a **sticky disk**, so warm runs
   skip the download entirely.
5. Checkout swapped to `useblacksmith/checkout@v1`: a drop-in for
   `actions/checkout` that clones from a git mirror cached next to the runner,
   which big repos feel the most.

Rule of thumb in the finished workflow: **disks for big state** (cargo
`target/`, Docker layers, Playwright browsers), **Actions cache for small
state** (pnpm store, cargo registry).

## Compare your results

Compare **per-job durations** between your first GitHub
run and your final run.

---

## Keep going

- After a week of real runs, ask Codesmith to `/rightsize` your longest
  workflow: it reads per-step CPU and memory headroom and recommends a runner
  size per job. (The report needs run history, which is why it's homework and
  not a workshop step.)
- Migrated the demo repo today? Do your real repo this week; you already
  have the account, the app, and 3,000 free minutes/month, and Step 1 took
  you ten minutes.
- Migrated a real repo today? Expand the app's scope to the next repo.
- Docs: <https://docs.blacksmith.sh> · Questions: find us at the booth.

## Appendix A: Cleanup

Used a throwaway org and want to leave no trace?

1. **Uninstall Blacksmith:** org → Settings → GitHub Apps → Blacksmith → Uninstall.
2. **Delete the org:** org → Settings → Delete this organization (removes the repo too).
3. Optionally delete your Blacksmith account from the dashboard settings.

(We'd rather you kept the free minutes, but the exit is always this easy.)

## Appendix B: the finished workflow

The complete `ci.yml` after Steps 1 to 3. Paste it over yours at any point to
catch up (use a new branch and open a PR so CI runs on it):

```yaml
# Blacksmith-Demo CI: the finished state after Steps 1 to 3. Three sticky
# disks (cargo target, Docker layers, Playwright browsers) plus the Actions
# cache for small state.
# Paste this whole file over .github/workflows/ci.yml (on a new branch,
# then open a PR) to catch up at any point.
name: CI

on:
  push:
    branches: [main]
  pull_request:
  workflow_dispatch:

concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true

jobs:
  web:
    name: web (pnpm + Vite + vitest)
    runs-on: blacksmith-2vcpu-ubuntu-2404
    steps:
      - uses: useblacksmith/checkout@v1
      - name: Start timer
        run: echo "JOB_T0=$(date +%s)" >> "$GITHUB_ENV"

      - uses: pnpm/action-setup@v4
      - uses: actions/setup-node@v7
        with:
          node-version: 22
          cache: pnpm

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Unit tests
        run: pnpm --filter @blacksmith-demo/web test

      - name: Build
        run: pnpm --filter @blacksmith-demo/web build

      - name: Report duration
        if: always()
        run: echo "⏱ **web** finished in **$(( $(date +%s) - JOB_T0 ))s**" >> "$GITHUB_STEP_SUMMARY"

  rust:
    name: rust (cargo build + test)
    runs-on: blacksmith-4vcpu-ubuntu-2404
    steps:
      - uses: useblacksmith/checkout@v1
      - name: Start timer
        run: echo "JOB_T0=$(date +%s)" >> "$GITHUB_ENV"

      - name: Build context
        run: echo "BLACKSMITH-DEMO-MARKER building stats service from a cold start"

      - name: Cache cargo registry
        uses: actions/cache@v4
        with:
          path: |
            ~/.cargo/registry
            ~/.cargo/git
          key: ${{ runner.os }}-cargo-${{ hashFiles('stats/Cargo.lock') }}

      - name: Mount sticky disk for build artifacts
        uses: useblacksmith/stickydisk@v1
        with:
          key: ${{ github.repository }}-cargo-target
          path: ./stats/target

      - name: Build and test
        working-directory: stats
        run: cargo test --locked

      - name: Report duration
        if: always()
        run: echo "⏱ **rust** finished in **$(( $(date +%s) - JOB_T0 ))s**" >> "$GITHUB_STEP_SUMMARY"

  integration:
    name: integration (Go + Postgres)
    runs-on: blacksmith-2vcpu-ubuntu-2404
    services:
      postgres:
        image: postgres:17-alpine
        env:
          POSTGRES_PASSWORD: demo
          POSTGRES_DB: demo
        ports:
          - 5432:5432
        options: >-
          --health-cmd "pg_isready -U postgres"
          --health-interval 5s
          --health-timeout 5s
          --health-retries 10
    steps:
      - uses: useblacksmith/checkout@v1
      - name: Start timer
        run: echo "JOB_T0=$(date +%s)" >> "$GITHUB_ENV"

      - uses: actions/setup-go@v6
        with:
          go-version-file: api/go.mod
          cache-dependency-path: api/go.sum

      - name: Tests (unit + Postgres integration)
        working-directory: api
        env:
          DATABASE_URL: postgres://postgres:demo@localhost:5432/demo
        run: go test -v ./...

      - name: Report duration
        if: always()
        run: echo "⏱ **integration** finished in **$(( $(date +%s) - JOB_T0 ))s**" >> "$GITHUB_STEP_SUMMARY"

  e2e:
    name: e2e (Playwright)
    runs-on: blacksmith-2vcpu-ubuntu-2404
    timeout-minutes: 8
    steps:
      - uses: useblacksmith/checkout@v1
      - name: Start timer
        run: echo "JOB_T0=$(date +%s)" >> "$GITHUB_ENV"

      - uses: pnpm/action-setup@v4
      - uses: actions/setup-node@v7
        with:
          node-version: 22
          cache: pnpm

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Mount sticky disk for Playwright browsers
        uses: useblacksmith/stickydisk@v1
        with:
          key: ${{ github.repository }}-playwright
          path: /home/runner/.cache/ms-playwright

      - name: Install Chromium
        run: pnpm --filter @blacksmith-demo/e2e exec playwright install --with-deps chromium

      - name: Build web app
        run: pnpm --filter @blacksmith-demo/web build

      - name: Run e2e tests
        run: pnpm --filter @blacksmith-demo/e2e test

      - name: Report duration
        if: always()
        run: echo "⏱ **e2e** finished in **$(( $(date +%s) - JOB_T0 ))s**" >> "$GITHUB_STEP_SUMMARY"

  docker:
    name: docker (container image)
    runs-on: blacksmith-2vcpu-ubuntu-2404
    steps:
      - uses: useblacksmith/checkout@v1
      - name: Start timer
        run: echo "JOB_T0=$(date +%s)" >> "$GITHUB_ENV"

      # Blacksmith's builder keeps the Docker layer cache on NVMe between runs.
      - name: Set up Docker builder
        uses: useblacksmith/setup-docker-builder@v2
        with:
          cache-key: blacksmith-demo-api

      - name: Build image
        uses: useblacksmith/build-push-action@v2
        with:
          context: .
          file: api/Dockerfile
          platforms: linux/amd64
          push: false
          tags: blacksmith-demo-api:ci

      - name: Report duration
        if: always()
        run: echo "⏱ **docker** finished in **$(( $(date +%s) - JOB_T0 ))s**" >> "$GITHUB_STEP_SUMMARY"
```
