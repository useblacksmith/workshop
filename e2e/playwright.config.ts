import { defineConfig, devices } from "@playwright/test";

// Serves the built web app with `vite preview`; API calls are mocked
// per-test with page.route(), so no backend is needed.
export default defineConfig({
  testDir: "./tests",
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? "github" : "list",
  use: {
    baseURL: "http://localhost:4173",
    trace: "on-first-retry",
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  webServer: {
    command: "pnpm --dir ../web preview --port 4173 --strictPort",
    url: "http://localhost:4173",
    reuseExistingServer: !process.env.CI,
  },
});
