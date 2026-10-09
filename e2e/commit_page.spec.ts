import { test, expect } from "@playwright/test";

test("test", async ({ page }) => {
  await page.goto("http://localhost:36000/history");
  await expect(
    page.getByRole("cell", { name: "2", exact: true }),
  ).toBeVisible();
  await page
    .getByRole("row", { name: "2 Empty commit for deletion" })
    .getByRole("button")
    .click();
  await expect(
    page.getByRole("cell", { name: "2", exact: true }),
  ).not.toBeAttached();
});
