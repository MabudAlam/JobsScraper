package main

import (
	"log"
	"net/http"
	"os"

	"ashbyimpl/scrapers/ashby/scheduler"
	"ashbyimpl/scrapers/ashby/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := setupRouter()

	if err := r.Run(":" + "8080"); err != nil {
		log.Panicf("error: %s", err)
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); exists == false {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	r.GET("/syncall", func(c *gin.Context) {
		password := c.Query("password")

		correctPassword := os.Getenv("SYNC_WITH_SQL_PASSWORD")

		if password != correctPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		scheduler.RunPipeline()
		c.JSON(http.StatusOK, gin.H{"status": "sync completed"})
	})

	r.GET("/getallJobsFromSQL", func(c *gin.Context) {
		_ = store.InitDB()
		jobs, err := store.GetAllActiveJobs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"jobs": jobs})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Job Scraper API"})
	})

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
