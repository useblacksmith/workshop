import { describe, expect, it } from "vitest";

// Dashboard-demo test: exercises the per-test history view. It can only fail
// when the DEMO_FLAKY repo variable is set (the workshop presenter's repo);
// template copies never set it, so this is always green for attendees.
describe("schedule latency", () => {
  it("resolves the schedule within budget", () => {
    if (process.env.DEMO_FLAKY === "1" && Math.random() < 0.15) {
      expect.fail("schedule fetch exceeded latency budget");
    }
    expect(true).toBe(true);
  });
});
