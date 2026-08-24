export interface Talk {
  id: string;
  title: string;
  speaker: string;
  room: string;
  slot: string;
  votes: number;
}

// Local fallback so the UI renders without the API (e.g. `vite preview`).
export const FALLBACK_TALKS: Talk[] = [
  { id: "keynote", title: "The Feedback Loop Is the Product", speaker: "Ada Okafor", room: "Main Stage", slot: "09:00", votes: 0 },
  { id: "ci-agents", title: "Your CI Was Sized for Humans", speaker: "Jonas Weber", room: "Track 2", slot: "10:30", votes: 0 },
  { id: "monorepo", title: "Monorepos Without Tears", speaker: "Priya Nair", room: "Track 1", slot: "11:15", votes: 0 },
  { id: "arm64", title: "Escaping QEMU: Native Multi-Arch Builds", speaker: "Sofia Lindqvist", room: "Track 2", slot: "13:00", votes: 0 },
  { id: "rust-perf", title: "Rust Compile Times: A Support Group", speaker: "Marco Bianchi", room: "Track 3", slot: "14:30", votes: 0 },
  { id: "postgres", title: "Postgres Is All You Need", speaker: "Elif Demir", room: "Track 1", slot: "16:00", votes: 0 },
];

export async function fetchTalks(): Promise<Talk[]> {
  try {
    const res = await fetch("/api/talks");
    if (!res.ok) throw new Error(`status ${res.status}`);
    return (await res.json()) as Talk[];
  } catch {
    return FALLBACK_TALKS;
  }
}

export function sortBySlot(talks: Talk[]): Talk[] {
  return [...talks].sort(
    (a, b) => a.slot.localeCompare(b.slot) || a.title.localeCompare(b.title),
  );
}

export function totalVotes(talks: Talk[]): number {
  return talks.reduce((sum, t) => sum + t.votes, 0);
}
