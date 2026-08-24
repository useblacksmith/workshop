// Blacksmith-Demo API: serves the conference schedule and collects talk votes.
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	store, err := newStore(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("store init: %v", err)
	}

	r := newRouter(store)
	addr := ":3000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("blacksmith-demo-api listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func newRouter(store Store) *gin.Engine {
	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.GET("/api/talks", func(c *gin.Context) {
		talks, err := store.Talks(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, talks)
	})

	r.POST("/api/talks/:id/vote", func(c *gin.Context) {
		count, err := store.Vote(c.Request.Context(), c.Param("id"))
		if err == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "no such talk"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "votes": count})
	})

	return r
}
