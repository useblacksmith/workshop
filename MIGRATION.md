# The CI Speedrun — workshop guide

Four steps, all in your browser, no local setup: migrate CI to
[Blacksmith](https://blacksmith.sh), put the expensive parts on persistent
disks, then let an agent finish the tuning and right-size your runners.

You'll do this either on **your own repo** (best — you leave with real CI
migrated) or on **our demo repo** (guaranteed to work for everyone).

---

## Step 0 — Setup (~10 min, do this first)

1. Get on the wifi, sign in at [github.com](https://github.com).
2. Sign in at [app.blacksmith.sh](https://app.blacksmith.sh) with GitHub.
   Don't install anything yet.
3. **Choose your repo:**

**Option A — your own repo** (pick this if you can):
- You need a GitHub **organization you can install apps on** (org admin, or
  an admin who'll approve fast) — personal side-project orgs count.
- And a repo in it with real CI that you're **allowed to experiment on**.

**Option B — our demo repo** (pick this otherwise):
1. Create a throwaway org: github.com → avatar → Settings → Organizations →
   **New organization** → Free plan. Name it `<github username>-blacksmith-workshop`. Skip
   inviting members. (Blacksmith installs on orgs, not personal accounts.
   [Appendix A](#appendix-a--cleanup) deletes all of this in one minute.)
2. On this repo: **Use this template → Create a new repository.**
   Owner: your new org. Visibility: **Public** (free unlimited Actions
   minutes). It's a small real app — Go API + Postgres, Rust service,
   TypeScript frontend, Playwright — with deliberately typical CI.

**Both options:** run the workflow once now on GitHub's runners
(**Actions → CI → Run workflow**) and note the **per-job durations** — job
times, not the run's wall clock, are your before-numbers for the scoreboard.

## Step 1 — Migrate the runners

1. Install the Blacksmith GitHub App from [app.blacksmith.sh](https://app.blacksmith.sh),
   **scoped to just your chosen repo** — the app only sees what you select.
2. Swap the runner labels:
   - **Your repo:** use the **migration wizard** in the Blacksmith dashboard —
     it opens the PR for you. Or hand-edit: every `runs-on: ubuntu-latest`
     becomes `runs-on: blacksmith-2vcpu-ubuntu-2404` (always the explicit
     label; per-job sizing like `4vcpu` is a feature, not a typo).
   - **Demo repo:** hand-edit `.github/workflows/ci.yml` in the GitHub web
     editor (press `.`) — five `runs-on` lines; give `rust` the
     `blacksmith-4vcpu-ubuntu-2404` label.
3. Merge/commit and run the workflow again.

Free bonus you didn't configure: every existing `actions/cache` /
`setup-node` cache is now served from a cache **colocated** with the runner —
same code, ~4x faster transfers.

> Demo repo, fell behind? `step-1-runners` has this done.

## Step 2 — Sticky disks

A **sticky disk** is a persistent NVMe volume that mounts into your runner in
seconds, with everything exactly as the last run left it. Give one to your
most expensive directory:

```yaml
      - name: Mount sticky disk
        uses: useblacksmith/stickydisk@v1
        with:
          key: ${{ github.repository }}-build-cache
          path: ./<expensive-directory>
```

**Recipe sheet — what to mount:**

| Ecosystem | Path to persist |
|---|---|
| Rust | `target/` |
| Turborepo / Nx | `.turbo/` / `.nx/cache` |
| Gradle | `~/.gradle/caches` |
| Go | `~/.cache/go-build` |
| Cypress / Playwright browsers | `~/.cache/Cypress` / `~/.cache/ms-playwright` |
| Docker layers | don't mount — swap to `useblacksmith/setup-docker-builder@v2` (with a `cache-key`) + `useblacksmith/build-push-action@v2`; the layer cache persists automatically |

**Demo repo:** mount `./stats/target` in the `rust` job, and make the Docker
builder swap shown above in the `docker` job.

Two expectations to set: the disk pays off on the **second** run, and on orgs
with **sticky-disk branch protection** enabled, disks commit only from the
default branch — PR runs read but don't warm, so merge before you measure.

> Demo repo checkpoint: `step-2-stickydisk`.

## Step 3 — Let Codesmith configure the rest

First, workshop credits: scan the QR on screen and submit your **GitHub
username** (make sure the org you installed today is the one selected in
your Blacksmith dashboard). Credits land on your org within a couple of
minutes.

Then comment on any PR or issue in your repo:

```
@codesmith find the expensive, uncached parts of this workflow and
open a PR configuring sticky disks and caching for them.
```

Codesmith reads your run history — step timings, cache misses, oversized
installs — and opens a PR. Review the diff, compare it with what you mounted
by hand in Step 2, and merge.

> Demo repo: `step-3-agent-optimized` mirrors the agent's PR if you'd rather
> not spend credits.

## Step 4 — Right-size your longest workflow

You now have at least two Blacksmith runs of history — enough for
right-sizing. In the Blacksmith dashboard, run **rightsize** against your
longest-running workflow: it analyzes per-step CPU and memory headroom and
recommends a runner size per job (bigger where you're compute-bound, smaller
where you're paying for idle cores). Review the recommendation card, apply
the rows you agree with — it edits the YAML and opens the PR — merge, and
run one last time.

## The scoreboard

Compare **per-job durations** (never wall clock) between your first GitHub
run and your final run. Post your biggest percentage speedup to the
leaderboard (QR on screen). Biggest speedup wins.

---

## Keep going

- Migrated the demo repo today? Do your real repo this week — you already
  have the account, the app, and 3,000 free minutes/month, and Step 1 took
  you ten minutes.
- Migrated a real repo today? Expand the app's scope to the next repo.
- Docs: <https://docs.blacksmith.sh> · Questions: find us at the booth.

## Appendix A — Cleanup

Used a throwaway org and want to leave no trace?

1. **Uninstall Blacksmith:** org → Settings → GitHub Apps → Blacksmith → Uninstall.
2. **Delete the org:** org → Settings → Delete this organization (removes the repo too).
3. Optionally delete your Blacksmith account from the dashboard settings.

(We'd rather you kept the free minutes — but the exit is always this easy.)

## Appendix B — Multi-arch images without QEMU

If you build `arm64` images on x86 runners today, you're paying the QEMU
emulation tax (~10x slower than native). Blacksmith has native arm64
runners — build each platform on its own hardware:

```yaml
  docker:
    strategy:
      matrix:
        include:
          - { platform: amd64, runner: blacksmith-2vcpu-ubuntu-2404 }
          - { platform: arm64, runner: blacksmith-2vcpu-ubuntu-2404-arm }
    runs-on: ${{ matrix.runner }}
    steps:
      - uses: actions/checkout@v7
      - uses: useblacksmith/setup-docker-builder@v2
        with:
          cache-key: my-image-${{ matrix.platform }}
      - uses: useblacksmith/build-push-action@v2
        with:
          platforms: linux/${{ matrix.platform }}
          ...
```
