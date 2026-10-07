import { expect, test } from "@playwright/test";
import { responseFor } from "./helpers";

test.describe("simple line page", () => {
  test("creates a line and reports success", async ({ page }) => {
    await page.goto("/simple");
    await expect(page.locator("#status-bar")).toHaveText("Idle");

    const upload = responseFor(page, "POST", /^\/upload\/lines/);
    await page.fill("#from-input", "Substation A");
    await page.fill("#to-input", "Substation B");
    await page.fill("#length-input", "10");
    await page.fill("#voltage-input", "400");
    await page.click("#submit-line");

    const response = await upload;
    expect(response.status()).toBe(200);

    // The page posts a single ndjson record, so assert on the decoded payload
    // instead of matching raw text.
    const payload = JSON.parse(response.request().postData() ?? "");
    expect(payload).toMatchObject({
      from: "Substation A",
      to: "Substation B",
      length: 10,
      voltage: 400,
    });

    await expect(page.locator("#status-bar")).toContainText("Successfully");
  });
});
