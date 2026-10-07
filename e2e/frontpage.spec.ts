import { expect, test } from "@playwright/test";
import { responseFor } from "./helpers";

test.describe("Load front page", () => {
  test("Has title TripleWorks", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("h1")).toHaveText("TripleWorks");
  });

  test("can create new item", async ({ page }) => {
    await page.goto("/");

    // The type select is populated by htmx on load, so the option may not be
    // in the DOM yet on first paint.
    await page.selectOption("#type-select", "ReportingGroup");
    await page.click("#new-btn");
    await page.fill("#commit-input", "create an empty reporting group");

    const commit = responseFor(page, "POST", "/commit");
    await page.click("#commit-btn");
    expect((await commit).status()).toBe(200);

    await expect(page.locator("#status-message")).toContainText(
      "Successfully updated",
    );
  });

  test("can download xiidm", async ({ request }) => {
    const response = await request.get("/xiidm");
    expect(response.status()).toBe(200);
    expect((await response.body()).length).toBeGreaterThan(10);
  });

  test("can retrieve cross region ptdfs", async ({ request }) => {
    const response = await request.get("/cross-region-ptdf");
    expect(response.status()).toBe(200);
  });
});
