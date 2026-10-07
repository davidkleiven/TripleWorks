import { expect, test } from "@playwright/test";
import { responseFor } from "./helpers";

const CONNECTED_LINE_MRID = "4e832836-ef53-458e-9711-903982551fcf";

test.describe("test map", () => {
  test("can add and remove production", async ({ page }) => {
    await page.goto("/map");

    // initMap creates the markers and exposes them on window, which is the only
    // way to open a popup without clicking the map canvas. The markers are
    // added asynchronously by the inline map data, so wait for them first.
    await page.waitForFunction(
      () => (window as any).substationMarkers?.length > 0,
    );
    await page.evaluate(() => (window as any).substationMarkers[0].openPopup());

    // htmx processes the popup contents on popupopen, so the button only
    // becomes active once the popup has been rendered.
    const produce = page.getByRole("button", { name: "Produce" });
    await expect(produce).toBeVisible();

    const activeForm = page.locator("#active-production-form");
    const numberInputs = activeForm.locator('input[type="number"]');
    const hiddenInputs = activeForm.locator('input[type="hidden"]');

    // The production endpoint answers with HX-Trigger-After-Swap, which makes
    // the page post the current injections to /flow.
    const flowRequest = responseFor(page, "POST", "/flow");
    await produce.click();

    // Only counts active
    await expect(numberInputs).toHaveCount(1);
    await expect(hiddenInputs).toHaveCount(1);

    const flowResponse = await flowRequest;
    expect(flowResponse.status()).toBe(200);
    const body = await flowResponse.json();
    expect(body.flow).toHaveProperty(CONNECTED_LINE_MRID);

    // Trigger flow recalculation
    const recalc = responseFor(page, "POST", "/flow");
    await page.click("#recalc-flow-btn");
    expect((await recalc).status()).toBe(200);

    // Check that that a flow label exists
    await expect(page.locator(".flow-value").first()).toBeVisible();

    await activeForm.locator('button[id^="delete-"]').click();
    await expect(page.locator(".flow-value")).toHaveCount(0);
    await expect(numberInputs).toHaveCount(0);
    await expect(hiddenInputs).toHaveCount(0);
  });
});
