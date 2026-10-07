import path from "node:path";
import { expect, test } from "@playwright/test";
import { responseFor } from "./helpers";

const JSON_PATCH_FIXTURE = path.resolve(__dirname, "fixtures/json_patch.json");

test.describe("can apply json patch", () => {
  test("can apply json patch", async ({ page }) => {
    await page.goto("/patch-form");
    await page.setInputFiles('input[type="file"]', JSON_PATCH_FIXTURE);

    const patched = responseFor(page, "PATCH", "/resource");
    await page.click('button[type="submit"]');
    expect((await patched).status()).toBe(200);
  });

  test("triggers connect-dangling lines on click", async ({ page }) => {
    await page.goto("/patch-form");
    await expect(page.locator("#model-selection")).toBeVisible();

    const connected = responseFor(page, "POST", "/connect-dangling");
    await page.click("#connect-dangling-lines-btn");

    const response = await connected;
    expect(response.request().postData()).toContain("modelId=");
    expect(response.status()).toBe(200);
    expect(response.headers()["content-type"]).toContain("text/html");
    expect(await response.text()).toContain("Inserted");
  });
});
