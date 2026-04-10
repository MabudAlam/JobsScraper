package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

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

		offset := 0
		limit := 80
		search := c.Query("search")
		company := c.Query("company")
		location := c.Query("location")
		sort := c.DefaultQuery("sort", "newest")

		if o := c.Query("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil {
				offset = parsed
			}
		}
		if l := c.Query("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		jobs, total, err := store.GetActiveJobsPaginated(offset, limit, search, company, location, sort)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"jobs":   jobs,
			"offset": offset,
			"limit":  limit,
			"total":  total,
		})
	})

	r.GET("/companies", func(c *gin.Context) {
		_ = store.InitDB()
		companies, err := store.GetAllCompanies()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"companies": companies})
	})

	r.GET("/locations", func(c *gin.Context) {
		_ = store.InitDB()
		locations, err := store.GetAllLocations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"locations": locations})
	})

	r.GET("/job/:id", func(c *gin.Context) {
		_ = store.InitDB()

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
			return
		}

		job, err := store.GetJobById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}

		c.JSON(http.StatusOK, job)
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
