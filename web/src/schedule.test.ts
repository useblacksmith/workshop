import { describe, expect, it } from "vitest";
import { FALLBACK_TALKS, sortBySlot, totalVotes } from "./schedule";

describe("sortBySlot", () => {
  it("orders talks by time slot", () => {
    const sorted = sortBySlot(FALLBACK_TALKS);
    expect(sorted[0].id).toBe("keynote");
    expect(sorted[sorted.length - 1].id).toBe("postgres");
  });

  it("does not mutate the input", () => {
    const input = [...FALLBACK_TALKS].reverse();
    const copy = [...input];
    sortBySlot(input);
    expect(input).toEqual(copy);
  });
});

describe("totalVotes", () => {
  it("sums votes across talks", () => {
    const talks = FALLBACK_TALKS.map((t, i) => ({ ...t, votes: i }));
    expect(totalVotes(talks)).toBe(15);
  });

  it("is zero for an empty schedule", () => {
    expect(totalVotes([])).toBe(0);
  });
});
