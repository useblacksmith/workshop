# The CI Speedrun — migration guide

You're going to take this repo's CI from **~9 minutes to ~1 minute**, in your
browser, in about 20 minutes of hands-on work — first by migrating the
hardware yourself, then by letting an agent clean up the workflow config.
No local setup needed.

Everything happens in a **throwaway GitHub org** that you create now and can
delete afterwards ([cleanup](#appendix-a--cleanup)). Nothing touches your
company's GitHub org, and no credit card is needed anywhere.

---

## Step 0 — Create your lab org (~2 min)

Blacksmith installs on GitHub **organizations**, not personal accounts, so
give yourself a disposable one:

1. Sign in at [github.com](https://github.com) (create an account if needed).
2. Click your avatar (top right) → **Settings** → **Organizations** → **New organization**
   — or go straight to <https://github.com/organizations/plan>.
3. Pick the **Free** plan.
4. Name it something like `<your-handle>-ci-lab`.
5. Skip inviting members — just you.

## Step 1 — Get your copy of this repo (~1 min)

1. On this repo's page, click **Use this template → Create a new repository**.
2. **Owner:** your new lab org. **Name:** `confab`. **Visibility:** Public
   (public repos get unlimited free GitHub Actions minutes).
3. Create the repository.

## Step 2 — Run the baseline (~1 min to start, ~9 min to finish)

1. In *your* repo: **Actions** tab → **CI** → **Run workflow** → run on `main`.
2. Watch five jobs fan out onto GitHub-hosted runners.
3. Note the long pole: the `docker` job, building an arm64 image under QEMU
   emulation.

☕ This takes ~9 minutes. That's the point. Back to the presentation.

## Step 3 — Sign up for Blacksmith (~3 min)

1. Go to [app.blacksmith.sh](https://app.blacksmith.sh) → **Sign in with GitHub**.
2. Install the Blacksmith GitHub App when prompted:
   - **Select your lab org** (not your company org — the app only sees what
     you scope it to).
3. That's it. There are no agents to deploy and no config files — any job
   whose `runs-on` label starts with `blacksmith-` now runs on Blacksmith.

You're on the free tier: 3,000 minutes/month, no credit card.

## Step 4 — Edit 1: swap the runner labels

Open `.github/workflows/ci.yml` in your repo (press `.` or use the ✏️ pencil
icon) and change every `runs-on` line:

| Job | Before | After |
|---|---|---|
| `web`, `integration`, `e2e`, `docker` | `runs-on: ubuntu-latest` | `runs-on: blacksmith-2vcpu-ubuntu-2404` |
| `rust` | `runs-on: ubuntu-latest` | `runs-on: blacksmith-4vcpu-ubuntu-2404` |

That's the whole migration for three of the five jobs. Bonus you didn't have
to configure: every `actions/cache` and `setup-node` cache in this file is
now served from Blacksmith's colocated cache — same code, ~4x faster
transfers.

Don't commit yet — two more edits.

> Lost? `git checkout step-1-runners` has this step done.

## Step 5 — Edit 2: put the cargo build on a sticky disk

Notice the `rust` job has **no caching at all** — it recompiles the entire
dependency tree every run. (Check your own repos before you judge.) The
expensive part of a Rust build is the `target/` directory, and the best home
for it is a **sticky disk** — a persistent NVMe volume that mounts into the
runner in seconds.

**Add** this step to the `rust` job, right before "Build and test":

```yaml
      - name: Mount sticky disk for build artifacts
        uses: useblacksmith/stickydisk@v1
        with:
          key: ${{ github.repository }}-cargo-target
          path: ./stats/target
```

> Lost? `git checkout step-2-stickydisk`.

## Step 6 — Edit 3: kill QEMU — build arm64 on real arm64

Replace the entire `docker` job with this: each platform builds on **native
hardware** via a matrix, and Blacksmith's Docker builder keeps your layer
cache on NVMe between runs.

```yaml
  docker:
    name: docker (${{ matrix.platform }})
    strategy:
      matrix:
        include:
          - platform: amd64
            runner: blacksmith-2vcpu-ubuntu-2404
          - platform: arm64
            runner: blacksmith-2vcpu-ubuntu-2404-arm
    runs-on: ${{ matrix.runner }}
    steps:
      - uses: actions/checkout@v7
      - name: Start timer
        run: echo "JOB_T0=$(date +%s)" >> "$GITHUB_ENV"

      - name: Set up Docker builder
        uses: useblacksmith/setup-docker-builder@v2
        with:
          cache-key: confab-api-${{ matrix.platform }}

      - name: Build image (native ${{ matrix.platform }})
        uses: useblacksmith/build-push-action@v2
        with:
          context: .
          file: api/Dockerfile
          platforms: linux/${{ matrix.platform }}
          push: false
          tags: confab-api:ci-${{ matrix.platform }}

      - name: Report duration
        if: always()
        run: echo "⏱ **docker/${{ matrix.platform }}** finished in **$(( $(date +%s) - JOB_T0 ))s**" >> "$GITHUB_STEP_SUMMARY"
```

Note what's *gone*: `setup-qemu-action`, `setup-buildx-action`, and any
`cache-from`/`cache-to` you'd normally maintain. Your config got shorter.

> Lost? `git checkout step-3-docker` — the finished state.

## Step 7 — Commit → cold run

Commit the edits to `main`. The push triggers your first Blacksmith run:
**everything is faster, but every cache is empty.** This is your worst-case
run (~3 min). While it builds, find your runs appearing in the
[Blacksmith dashboard](https://app.blacksmith.sh).

## Step 8 — Run again → warm run

**Actions** → **CI** → **Run workflow** once more. Now the Docker layer
cache and the sticky disk are primed: **~2 minutes wall clock.**

Fast — but look closer at the logs. `pnpm install` still downloads every
package. Playwright still installs three browsers (the tests use one). The
cargo registry still re-fetches. The hardware is fixed; the *workflow config*
is still the one your team wrote in a hurry two years ago. That's the next
step.

## Step 9 — Let the agent finish the job

You've been fixing this workflow by hand. Now watch the other half of the
story: comment on any PR or issue in your repo —

```
@codesmith this workflow wastes time on every run — find the config
problems and open a PR fixing them.
```

Codesmith reads your run history and step timings, then opens a PR doing
what a careful engineer would: dependency caches with lockfile keys,
Chromium-only browser install, cancel-superseded-runs concurrency. Review
the diff, merge, **Run workflow** one last time: **~1 minute.**

> No agent budget, or want to see the answer? `git checkout step-4-optimized`
> is the same PR, pre-baked.

Four runs on your screen: **~9:00 → ~3:00 → ~2:00 → ~1:00.** The first jump
was hardware. The last one was an agent doing CI hygiene — the same kind of
agent that's about to multiply your CI load is also the thing that keeps it
tuned.

Post your best wall-clock time to the leaderboard (QR on screen). Fastest
speedrun wins.

---

## Part 2 — Do it for real (homework)

You already have a Blacksmith account and 3,000 free minutes/month. Pick ONE
real workflow in your own org this week:

1. Install the Blacksmith app on the org (or ask your admin), scoped to one
   repo to start.
2. Swap `runs-on` labels — always the explicit form, e.g.
   `blacksmith-2vcpu-ubuntu-2404` (or use the migration wizard in the
   dashboard, which opens the PR for you).
3. If you build Docker images: `useblacksmith/setup-docker-builder@v2`
   (+ `cache-key`) and `useblacksmith/build-push-action@v2`, and delete your
   `cache-from`/`cache-to` lines.
4. If you have a big compiled artifact dir (cargo target, Bazel cache, Gradle
   cache): mount a `useblacksmith/stickydisk@v1`.
5. Keep `actions/cache`/`setup-node`/`setup-go` exactly as they are — they're
   automatically colocated.

Docs: <https://docs.blacksmith.sh> · Questions: come by the booth.

## Appendix A — Cleanup

Want to leave no trace? Takes one minute:

1. **Uninstall Blacksmith:** your lab org → **Settings** → **GitHub Apps** →
   Blacksmith → **Uninstall**.
2. **Delete the org:** your lab org → **Settings** → scroll down → **Delete
   this organization** (this deletes the repo too).
3. Optionally delete your Blacksmith account from the dashboard settings.

(We'd rather you kept it — you have 3,000 free minutes a month — but the
exit is always this easy.)
