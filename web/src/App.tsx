import { useEffect, useState } from "react";
import { fetchTalks, sortBySlot, totalVotes, type Talk } from "./schedule";

export function App() {
  const [talks, setTalks] = useState<Talk[]>([]);

  useEffect(() => {
    fetchTalks().then((t) => setTalks(sortBySlot(t)));
  }, []);

  async function vote(id: string) {
    setTalks((prev) =>
      prev.map((t) => (t.id === id ? { ...t, votes: t.votes + 1 } : t)),
    );
    try {
      await fetch(`/api/talks/${id}/vote`, { method: "POST" });
    } catch {
      // Offline / preview mode: keep the optimistic update.
    }
  }

  return (
    <main style={{ fontFamily: "system-ui, sans-serif", maxWidth: 640, margin: "2rem auto", padding: "0 1rem" }}>
      <h1>Blacksmith-Demo</h1>
      <p>
        Vote for the talks you want to see. Total votes:{" "}
        <strong data-testid="total-votes">{totalVotes(talks)}</strong>
      </p>
      <ul style={{ listStyle: "none", padding: 0 }}>
        {talks.map((talk) => (
          <li
            key={talk.id}
            style={{ display: "flex", alignItems: "center", gap: "1rem", padding: "0.75rem 0", borderBottom: "1px solid #ddd" }}
          >
            <span style={{ minWidth: "3.5rem", color: "#666" }}>{talk.slot}</span>
            <span style={{ flex: 1 }}>
              <strong>{talk.title}</strong>
              <br />
              <small>
                {talk.speaker} · {talk.room}
              </small>
            </span>
            <button onClick={() => vote(talk.id)} aria-label={`Vote for ${talk.title}`}>
              ▲ {talk.votes}
            </button>
          </li>
        ))}
      </ul>
    </main>
  );
}
