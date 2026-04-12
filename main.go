package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"ashbyimpl/embeddings"
	"ashbyimpl/scrapers/ashby/scheduler"

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

	embeddings.InitDB()

	r.GET("/syncall", func(c *gin.Context) {
		password := c.Query("password")

		correctPassword := os.Getenv("SYNC_PASSWORD")

		if password != correctPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		scheduler.RunPipeline()
		c.JSON(http.StatusOK, gin.H{"status": "sync completed"})
	})

	r.POST("/reembed", func(c *gin.Context) {
		password := c.Query("password")
		correctPassword := os.Getenv("SYNC_PASSWORD")

		if password != correctPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		apiKey := os.Getenv("OPENROUTER_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OPENROUTER_API_KEY not set"})
			return
		}

		embedSvc := embeddings.NewEmbeddingService(apiKey)

		jobs, err := embeddings.GetAllActiveJobs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Re-embedding %d jobs", len(jobs))

		texts := make([]string, len(jobs))
		for i, job := range jobs {
			texts[i] = embedSvc.GenerateJobText(&job)
		}

		vectors, err := embedSvc.EmbedTexts(texts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		success := 0
		for i, job := range jobs {
			if i < len(vectors) {
				job.Embedding = vectors[i]
				if _, err := embeddings.UpsertJob(&job); err == nil {
					success++
				}
			}
		}

		log.Printf("Successfully re-embedded %d/%d jobs", success, len(jobs))
		c.JSON(http.StatusOK, gin.H{
			"status":  "completed",
			"total":   len(jobs),
			"success": success,
		})
	})

	r.GET("/getallJobsFromSQL", func(c *gin.Context) {
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

		jobs, total, err := embeddings.GetActiveJobsPaginated(offset, limit, search, company, location, sort)
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
		companies, err := embeddings.GetAllCompanies()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"companies": companies})
	})

	r.GET("/locations", func(c *gin.Context) {
		locations, err := embeddings.GetAllLocations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"locations": locations})
	})

	r.GET("/job/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
			return
		}

		job, err := embeddings.GetJobById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}

		c.JSON(http.StatusOK, job)
	})

	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
			return
		}

		limit := 20
		if l := c.Query("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		company := c.Query("company")
		location := c.Query("location")

		vectors, err := embeddings.EmbedJobTexts([]string{query})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		results, err := embeddings.SearchJobs(vectors[0], limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var filteredResults []map[string]interface{}
		for _, r := range results {
			if company != "" {
				if comp, ok := r["company"].(string); ok && comp != company {
					continue
				}
			}
			if location != "" {
				if loc, ok := r["location"].(string); ok && loc != location {
					continue
				}
			}
			filteredResults = append(filteredResults, r)
		}

		c.JSON(http.StatusOK, gin.H{
			"results": filteredResults,
			"query":   query,
			"total":   len(filteredResults),
		})
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
