import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { App } from "./App";

describe("App", () => {
  beforeEach(() => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () =>
          Promise.resolve([
            { id: "keynote", title: "The Feedback Loop Is the Product", speaker: "Ada Okafor", room: "Main Stage", slot: "09:00", votes: 2 },
            { id: "arm64", title: "Escaping QEMU: Native Multi-Arch Builds", speaker: "Sofia Lindqvist", room: "Track 2", slot: "13:00", votes: 3 },
          ]),
      }),
    );
  });

  it("renders the schedule and vote totals", async () => {
    render(<App />);
    expect(await screen.findByText("The Feedback Loop Is the Product")).toBeInTheDocument();
    expect(screen.getByText("Escaping QEMU: Native Multi-Arch Builds")).toBeInTheDocument();
    expect(screen.getByTestId("total-votes")).toHaveTextContent("5");
  });
});
