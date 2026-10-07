import { defineConfig, devices } from "@playwright/test";

// The suite shares a single SQLite database (tripleworks-e2e.db) that is
// deleted and reseeded by `make e2e`, and several specs mutate that shared
// state. Keep the workers at 1 so the specs run in isolation from each other.
export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false,
  workers: 1,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI
    ? [["github"], ["html", { open: "never" }]]
    : [["list"], ["html", { open: "never" }]],
  use: {
    baseURL: "http://localhost:36000",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
  webServer: {
    command: "make e2e",
    url: "http://localhost:36000",
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
});
