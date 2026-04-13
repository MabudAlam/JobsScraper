package handler

import (
	"context"
	"crypto/subtle"
	"log"
	"net/http"
	"os"

	scraperpkg "jobscraper/internal/scraper"

	"github.com/gin-gonic/gin"
)

func SyncAll(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	expectedPassword := os.Getenv("SYNC_PASSWORD")
	if expectedPassword == "" {
		log.Printf("SYNC_PASSWORD environment variable is not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sync not configured"})
		return
	}

	if authHeader == "" || subtle.ConstantTimeCompare([]byte(authHeader), []byte("Bearer "+expectedPassword)) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	scrapers := scraperpkg.GetAllScrapers()
	pool := scraperpkg.NewPool(scrapers, 4, 15)

	log.Println("Starting parallel sync")
	results, _ := pool.Run(context.Background())

	totalCount := 0
	for _, r := range results {
		totalCount += r.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "synced successfully",
		"count":   totalCount,
		"results": results,
	})
}

func SyncWithPassword(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	expectedPassword := os.Getenv("SYNC_PASSWORD")
	if expectedPassword == "" {
		log.Printf("SYNC_PASSWORD environment variable is not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sync not configured"})
		return
	}

	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(expectedPassword)) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	scrapers := scraperpkg.GetAllScrapers()
	pool := scraperpkg.NewPool(scrapers, 4, 15)

	log.Println("Starting parallel sync")
	results, _ := pool.Run(context.Background())

	totalCount := 0
	for _, r := range results {
		totalCount += r.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "synced successfully",
		"count":   totalCount,
		"results": results,
	})
}
