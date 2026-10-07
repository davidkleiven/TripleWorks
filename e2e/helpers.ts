import { Page } from "@playwright/test";

/**
 * Waits for a response to `method` + `path`, where `path` is either an exact
 * pathname or a regular expression tested against the pathname. Register the
 * waiter before triggering the action it should observe.
 */
export function responseFor(
  page: Page,
  method: string,
  path: string | RegExp,
): Promise<import("@playwright/test").Response> {
  return page.waitForResponse((response) => {
    const pathname = new URL(response.url()).pathname;
    const methodMatches = response.request().method() === method;
    const pathMatches =
      path instanceof RegExp ? path.test(pathname) : pathname === path;
    return methodMatches && pathMatches;
  });
}

/** Counts non-overlapping occurrences of `needle` in `haystack`. */
export function countOccurrences(haystack: string, needle: string): number {
  return haystack.split(needle).length - 1;
}

/**
 * Waits until htmx has processed `selector`.
 *
 * htmx attaches its listeners to swapped-in content asynchronously after the
 * swap, so an element can be present and visible before it reacts to events.
 * Typing into a searchable picker that early is silently dropped, so wait for
 * the internal data marker htmx records on every element it processes.
 */
export async function waitForHtmxProcessed(
  page: Page,
  selector: string,
): Promise<void> {
  await page.waitForFunction((sel) => {
    const element = document.querySelector(sel) as unknown as Record<
      string,
      unknown
    > | null;
    return !!element?.["htmx-internal-data"];
  }, selector);
}
