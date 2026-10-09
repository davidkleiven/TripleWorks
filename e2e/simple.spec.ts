import { expect, test } from "@playwright/test";
import { responseFor } from "./helpers";

test.describe("simple line page", () => {
  test("creates a line and reports success", async ({ page, request }) => {
    await page.goto("/simple");
    await expect(page.locator("#status-bar")).toHaveText("Idle");

    const upload = responseFor(page, "POST", /\/upload\/lines.*/);
    await page.selectOption("#model-selection", "1");
    await page.fill("#from-input", "Substation A");
    await page.fill("#to-input", "Substation B");
    await page.fill("#length-input", "10");
    await page.fill("#voltage-input", "400");
    await page.click("#submit-line");

    const response = await upload;
    expect(response.status()).toBe(200);

    // The chosen model travels as a query parameter, next to commit.
    expect(new URL(response.url()).searchParams.get("modelId")).toBe("1");

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

    // A 200 only proves the upload was accepted; the line is only persisted
    // when the request carries commit=true, so look it up afterwards.
    const lines = await request.get("/entities?kind=ACLineSegment");
    expect(lines.status()).toBe(200);
    // The line name is derived from the substation pair and the voltage.
    expect(await lines.text()).toContain("Substation A-Substation B (400 kV)");
  });

  test("creates substations", async ({ page, request }) => {
    await page.goto("/simple");
    await expect(page.locator("#status-bar")).toHaveText("Idle");

    const upload = responseFor(page, "POST", /\/upload\/substations.*/);
    await page.selectOption("#model-selection", "1");
    await page.fill("#name-input", "New substation");
    await page.fill("#region-input", "NO2");
    await page.fill("#long-input", "56.23");
    await page.fill("#lat-input", "6.27");
    await page.click("#submit-substation");

    const response = await upload;
    expect(response.status()).toBe(200);
    const lines = await request.get("/entities?kind=Substation");
    expect(lines.status()).toBe(200);
    expect(await lines.text()).toContain("New substation");
  });
});

test.describe("simple page file upload", () => {
  test("creates many lines from a jsonl file", async ({ page, request }) => {
    await page.goto("/simple");
    await expect(page.locator("#status-bar")).toHaveText("Idle");

    const upload = responseFor(page, "POST", /\/upload\/lines.*/);
    await page.selectOption("#model-selection", "1");

    // Names must not collide with the single-item tests: the server silently
    // drops records whose mrids already exist in the model.
    const body = [
      { from: "Bulk Sub A", to: "Bulk Sub B", length: 10, voltage: 400 },
      { from: "Bulk Sub C", to: "Bulk Sub D", length: 20, voltage: 220 },
    ]
      .map((record) => JSON.stringify(record))
      .join("\n");

    await page.setInputFiles("#line-file", {
      name: "lines.jsonl",
      mimeType: "application/x-ndjson",
      buffer: Buffer.from(body),
    });

    const response = await upload;
    expect(response.status()).toBe(200);
    expect(new URL(response.url()).searchParams.get("modelId")).toBe("1");

    // The whole file travels as one ndjson request.
    expect(response.request().postData()).toBe(body + "\n");

    await expect(page.locator("#status-bar")).toContainText(
      "Successfully created 2 items",
    );

    const lines = await request.get("/entities?kind=ACLineSegment");
    expect(lines.status()).toBe(200);
    expect(await lines.text()).toContain("Bulk Sub A-Bulk Sub B (400 kV)");
    expect(await lines.text()).toContain("Bulk Sub C-Bulk Sub D (220 kV)");
  });

  test("creates many substations from a jsonl file", async ({
    page,
    request,
  }) => {
    await page.goto("/simple");
    await expect(page.locator("#status-bar")).toHaveText("Idle");

    const upload = responseFor(page, "POST", /\/upload\/substations.*/);
    await page.selectOption("#model-selection", "1");

    const body = [
      { name: "Bulk Sub E", region: "NO2", x: 10.5, y: 60.1 },
      { name: "Bulk Sub F", region: "NO2", x: 11.5, y: 61.1 },
    ]
      .map((record) => JSON.stringify(record))
      .join("\n");

    await page.setInputFiles("#substation-file", {
      name: "substations.jsonl",
      mimeType: "application/x-ndjson",
      buffer: Buffer.from(body),
    });

    const response = await upload;
    expect(response.status()).toBe(200);
    await expect(page.locator("#status-bar")).toContainText(
      "Successfully created 2 items",
    );

    const substations = await request.get("/entities?kind=Substation");
    expect(substations.status()).toBe(200);
    expect(await substations.text()).toContain("Bulk Sub E");
    expect(await substations.text()).toContain("Bulk Sub F");
  });

  test("reports the server error when a record is malformed", async ({
    page,
  }) => {
    await page.goto("/simple");
    await expect(page.locator("#status-bar")).toHaveText("Idle");

    const upload = responseFor(page, "POST", /\/upload\/substations.*/);
    await page.selectOption("#model-selection", "1");

    await page.setInputFiles("#substation-file", {
      name: "broken.jsonl",
      mimeType: "application/x-ndjson",
      buffer: Buffer.from(
        '{"name":"Good Sub","region":"NO2","x":1,"y":2}\nnot json\n',
      ),
    });

    const response = await upload;
    expect(response.status()).toBe(400);
    await expect(page.locator("#status-bar")).toContainText(
      "Could not unmarshal line",
    );
  });
});
