import { expect, test } from "@playwright/test";

const TALKS = [
  { id: "keynote", title: "The Feedback Loop Is the Product", speaker: "Ada Okafor", room: "Main Stage", slot: "09:00", votes: 4 },
  { id: "arm64", title: "Escaping QEMU: Native Multi-Arch Builds", speaker: "Sofia Lindqvist", room: "Track 2", slot: "13:00", votes: 7 },
];

test.beforeEach(async ({ page }) => {
  await page.route("**/api/talks", (route) =>
    route.fulfill({ json: TALKS }),
  );
  await page.route("**/api/talks/*/vote", (route) =>
    route.fulfill({ json: { id: "keynote", votes: 5 } }),
  );
});

test("shows the schedule", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/Blacksmith-Demo/);
  await expect(page.getByText("The Feedback Loop Is the Product")).toBeVisible();
  await expect(page.getByText("Sofia Lindqvist · Track 2")).toBeVisible();
  await expect(page.getByTestId("total-votes")).toHaveText("11");
});

test("voting bumps the count", async ({ page }) => {
  await page.goto("/");
  const voteButton = page.getByRole("button", {
    name: "Vote for The Feedback Loop Is the Product",
  });
  await expect(voteButton).toHaveText("▲ 4");
  await voteButton.click();
  await expect(voteButton).toHaveText("▲ 5");
  await expect(page.getByTestId("total-votes")).toHaveText("12");
});
