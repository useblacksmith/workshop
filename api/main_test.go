package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestTalksEndpoint(t *testing.T) {
	r := newRouter(newMemStore())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/talks", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", w.Code)
	}
	var talks []Talk
	if err := json.Unmarshal(w.Body.Bytes(), &talks); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(talks) != 6 {
		t.Fatalf("got %d talks, want 6", len(talks))
	}
}

func TestVoteEndpoint(t *testing.T) {
	r := newRouter(newMemStore())

	for i := 1; i <= 3; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/talks/keynote/vote", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("got status %d, want 200", w.Code)
		}
		var resp struct {
			Votes int `json:"votes"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if resp.Votes != i {
			t.Fatalf("got %d votes after %d votes cast", resp.Votes, i)
		}
	}
}

func TestVoteUnknownTalk(t *testing.T) {
	r := newRouter(newMemStore())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/talks/nope/vote", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want 404", w.Code)
	}
}
