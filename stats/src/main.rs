//! Blacksmith-Demo stats: pulls the schedule from the API and serves vote analytics.

use axum::{extract::State, routing::get, Json, Router};
use chrono::Utc;
use serde::{Deserialize, Serialize};
use std::sync::Arc;

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Talk {
    pub id: String,
    pub title: String,
    pub speaker: String,
    pub room: String,
    pub slot: String,
    pub votes: i64,
}

#[derive(Debug, Serialize, PartialEq)]
pub struct Stats {
    pub total_votes: i64,
    pub talk_count: usize,
    pub mean_votes: f64,
    pub leader: Option<String>,
    pub votes_by_room: Vec<RoomVotes>,
}

#[derive(Debug, Serialize, PartialEq)]
pub struct RoomVotes {
    pub room: String,
    pub votes: i64,
}

/// Aggregate vote statistics across the schedule.
pub fn compute_stats(talks: &[Talk]) -> Stats {
    let total_votes: i64 = talks.iter().map(|t| t.votes).sum();
    let talk_count = talks.len();
    let mean_votes = if talk_count == 0 {
        0.0
    } else {
        total_votes as f64 / talk_count as f64
    };
    let leader = talks
        .iter()
        .max_by_key(|t| t.votes)
        .filter(|t| t.votes > 0)
        .map(|t| t.id.clone());

    let mut by_room: Vec<RoomVotes> = Vec::new();
    for t in talks {
        match by_room.iter_mut().find(|r| r.room == t.room) {
            Some(r) => r.votes += t.votes,
            None => by_room.push(RoomVotes {
                room: t.room.clone(),
                votes: t.votes,
            }),
        }
    }
    by_room.sort_by(|a, b| b.votes.cmp(&a.votes).then(a.room.cmp(&b.room)));

    Stats {
        total_votes,
        talk_count,
        mean_votes,
        leader,
        votes_by_room: by_room,
    }
}

struct AppState {
    api_url: String,
}

async fn stats_handler(State(state): State<Arc<AppState>>) -> Json<serde_json::Value> {
    let talks: Vec<Talk> = match reqwest::get(format!("{}/api/talks", state.api_url)).await {
        Ok(resp) => resp.json().await.unwrap_or_default(),
        Err(_) => Vec::new(),
    };
    let stats = compute_stats(&talks);
    Json(serde_json::json!({
        "generated_at": Utc::now().to_rfc3339(),
        "stats": stats,
    }))
}

async fn health_handler() -> Json<serde_json::Value> {
    Json(serde_json::json!({ "ok": true }))
}

#[tokio::main]
async fn main() {
    let api_url =
        std::env::var("DEMO_API_URL").unwrap_or_else(|_| "http://localhost:3000".into());
    let state = Arc::new(AppState { api_url });

    let app = Router::new()
        .route("/stats", get(stats_handler))
        .route("/health", get(health_handler))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3001").await.unwrap();
    println!("blacksmith-demo-stats listening on :3001");
    axum::serve(listener, app).await.unwrap();
}

#[cfg(test)]
mod tests {
    use super::*;

    fn talk(id: &str, room: &str, votes: i64) -> Talk {
        Talk {
            id: id.into(),
            title: format!("Talk {id}"),
            speaker: "Speaker".into(),
            room: room.into(),
            slot: "09:00".into(),
            votes,
        }
    }

    #[test]
    fn empty_schedule() {
        let stats = compute_stats(&[]);
        assert_eq!(stats.total_votes, 0);
        assert_eq!(stats.talk_count, 0);
        assert_eq!(stats.mean_votes, 0.0);
        assert_eq!(stats.leader, None);
        assert!(stats.votes_by_room.is_empty());
    }

    #[test]
    fn totals_and_leader() {
        let talks = vec![
            talk("a", "Main Stage", 10),
            talk("b", "Track 1", 4),
            talk("c", "Track 1", 6),
        ];
        let stats = compute_stats(&talks);
        assert_eq!(stats.total_votes, 20);
        assert_eq!(stats.talk_count, 3);
        assert!((stats.mean_votes - 20.0 / 3.0).abs() < 1e-9);
        assert_eq!(stats.leader, Some("a".into()));
    }

    #[test]
    fn no_leader_without_votes() {
        let talks = vec![talk("a", "Main Stage", 0), talk("b", "Track 1", 0)];
        assert_eq!(compute_stats(&talks).leader, None);
    }

    #[test]
    fn votes_grouped_by_room_sorted_desc() {
        let talks = vec![
            talk("a", "Track 1", 4),
            talk("b", "Track 2", 9),
            talk("c", "Track 1", 3),
        ];
        let stats = compute_stats(&talks);
        assert_eq!(
            stats.votes_by_room,
            vec![
                RoomVotes { room: "Track 2".into(), votes: 9 },
                RoomVotes { room: "Track 1".into(), votes: 7 },
            ]
        );
    }
}
