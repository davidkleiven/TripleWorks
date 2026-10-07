import { expect, Page, test } from "@playwright/test";
import { countOccurrences, responseFor, waitForHtmxProcessed } from "./helpers";

const UNCONNECTED_LINE_MRID = "ce8e57c7-8f6c-42c3-8b8e-e06aa39f0da3";

async function openSubstationConnector(page: Page) {
  await page.goto("/");
  await page.click("#list-all-btn");
  await page
    .locator(`button[hx-get*="${UNCONNECTED_LINE_MRID}"]`)
    .filter({ hasText: "Con sub." })
    .click();
  await expect(page.locator("h2")).toHaveText("Substation connector");
  await waitForHtmxProcessed(page, "#from-substation-search-input");
  await waitForHtmxProcessed(page, "#to-substation-search-input");
}

/**
 * Fills both searchable pickers and picks the requested substations. Passing
 * the new names replaces the previous selection.
 */
async function pickSubstations(
  page: Page,
  from: string,
  to: string,
): Promise<void> {
  await page.locator("#from-substation-search-input").fill(from);
  await page.locator("#to-substation-search-input").fill(to);

  const fromResults = page.locator("#from-substation-results");
  const toResults = page.locator("#to-substation-results");
  await expect(fromResults.locator("span").first()).toBeVisible();
  await expect(toResults.locator("span").first()).toBeVisible();

  await fromResults.locator("span", { hasText: from }).first().click();
  await toResults.locator("span", { hasText: to }).first().click();

  await expect(page.locator("#from-substation-display")).toContainText(from);
  await expect(page.locator("#to-substation-display")).toContainText(to);
}

// The two tests below share one server and one database, and the move test
// asserts on the result of connecting the line. Declaration order matters, so
// keep these in the same file and run with a single worker.
test.describe("substation connector works", () => {
  test("connect sends correct body", async ({ page }) => {
    await openSubstationConnector(page);
    await pickSubstations(page, "Substation A", "Substation B");

    await expect(page.locator("#from-substation-display")).not.toContainText(
      "No selection",
    );
    await expect(page.locator("#to-substation-display")).not.toContainText(
      "No selection",
    );

    const connectRequest = responseFor(page, "POST", /^\/connect\//);
    await page.click("#connect-substations-btn");

    const response = await connectRequest;
    const body = response.request().postData() ?? "";
    expect(body).toContain("modelId=");
    expect(countOccurrences(body, "substation-mrid=")).toBe(2);
    expect(response.status()).toBe(200);

    await expect(page.locator("#status-message")).toContainText(
      "Successfully committed",
    );
  });

  test("move sends correct body", async ({ page }) => {
    await openSubstationConnector(page);
    await pickSubstations(page, "Substation A", "Substation B");

    // Moving to the substations the line already is connected to does nothing
    const noopMove = responseFor(page, "POST", /^\/move\//);
    await page.click("#move-substations-btn");

    const noopResponse = await noopMove;
    const noopBody = noopResponse.request().postData() ?? "";
    expect(countOccurrences(noopBody, "substation-mrid=")).toBe(2);
    expect(noopResponse.status()).toBe(200);
    await expect(page.locator("#status-message")).toContainText(
      "Nothing to move",
    );

    // Move both ends of the line to Substation C and Substation D
    await pickSubstations(page, "Substation C", "Substation D");

    const moveRequest = responseFor(page, "POST", /^\/move\//);
    await page.click("#move-substations-btn");

    expect((await moveRequest).status()).toBe(200);
    await expect(page.locator("#status-message")).toContainText(
      "Successfully moved",
    );
  });
});
