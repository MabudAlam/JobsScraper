package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"jobscraper/common"
	amazonfetch "jobscraper/scrapers/amazon/fetch"
	amazonnormalize "jobscraper/scrapers/amazon/normalize"
	ashbyfetch "jobscraper/scrapers/ashby/fetch"
	ashbynormalize "jobscraper/scrapers/ashby/normalize"
	atlassianfetch "jobscraper/scrapers/atlassian/fetch"
	atlassiannormalize "jobscraper/scrapers/atlassian/normalize"
	leverfetch "jobscraper/scrapers/lever/fetch"
	levernormalize "jobscraper/scrapers/lever/normalize"
	"jobscraper/target"

	embeddings "jobscraper/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type LeverCompany struct {
	Company        string `json:"Company"`
	LeverSlug      string `json:"LeverSlug"`
	Enabled        bool   `json:"Enabled"`
	FrequencyHours int    `json:"FrequencyHours"`
}

func getLeverCompanies() []LeverCompany {
	data, err := os.ReadFile("levercompanies.json")
	if err != nil {
		return []LeverCompany{}
	}
	var companies []LeverCompany
	if err := json.Unmarshal(data, &companies); err != nil {
		return []LeverCompany{}
	}
	return companies
}

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

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Job Scraper API"})
	})

	r.GET("/sync", func(c *gin.Context) {
		raw, err := amazonfetch.FetchAmazonJobs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		jobs := amazonnormalize.NormalizeAmazonJobs(raw)
		if len(jobs) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "no jobs found"})
			return
		}

		var payloads []*common.JobPayload
		for _, j := range jobs {
			payloads = append(payloads, j)
		}

		if err := embeddings.InsertJobs(payloads); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "synced successfully",
			"count":   len(jobs),
		})
	})

	r.GET("/syncall", func(c *gin.Context) {
		password := c.Query("password")
		if password != "password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			return
		}

		type SyncResult struct {
			Company string `json:"company"`
			Status  string `json:"status"`
			Count   int    `json:"count,omitempty"`
			Error   string `json:"error,omitempty"`
		}

		var results []SyncResult
		totalCount := 0

		log.Println("Starting sync - Amazon")
		amazonResult := SyncResult{Company: "Amazon"}
		raw, err := amazonfetch.FetchAmazonJobs()
		if err != nil {
			amazonResult.Status = "failed"
			amazonResult.Error = err.Error()
			log.Printf("Amazon: FAILED - %v", err)
		} else if len(raw.Jobs) > 0 {
			jobs := amazonnormalize.NormalizeAmazonJobs(raw)
			if len(jobs) > 0 {
				var payloads []*common.JobPayload
				for _, j := range jobs {
					payloads = append(payloads, j)
				}
				if err := embeddings.InsertJobs(payloads); err != nil {
					amazonResult.Status = "failed"
					amazonResult.Error = err.Error()
					log.Printf("Amazon: FAILED - insert error: %v", err)
				} else {
					amazonResult.Status = "success"
					amazonResult.Count = len(jobs)
					totalCount += len(jobs)
					log.Printf("Amazon: SUCCESS - %d jobs", len(jobs))
				}
			} else {
				amazonResult.Status = "success"
				amazonResult.Count = 0
				log.Println("Amazon: SUCCESS - 0 jobs")
			}
		}
		results = append(results, amazonResult)

		log.Println("Starting sync - Atlassian")
		atlassianResult := SyncResult{Company: "Atlassian"}
		atlassianRaw, err := atlassianfetch.FetchAtlassianJobs()
		if err != nil {
			atlassianResult.Status = "failed"
			atlassianResult.Error = err.Error()
			log.Printf("Atlassian: FAILED - %v", err)
		} else if len(*atlassianRaw) > 0 {
			atlassianJobs := atlassiannormalize.NormalizeAtlassianJobs(atlassianRaw)
			if len(atlassianJobs) > 0 {
				var payloads []*common.JobPayload
				for _, j := range atlassianJobs {
					payloads = append(payloads, j)
				}
				if err := embeddings.InsertJobs(payloads); err != nil {
					atlassianResult.Status = "failed"
					atlassianResult.Error = err.Error()
					log.Printf("Atlassian: FAILED - insert error: %v", err)
				} else {
					atlassianResult.Status = "success"
					atlassianResult.Count = len(atlassianJobs)
					totalCount += len(atlassianJobs)
					log.Printf("Atlassian: SUCCESS - %d jobs", len(atlassianJobs))
				}
			} else {
				atlassianResult.Status = "success"
				atlassianResult.Count = 0
				log.Println("Atlassian: SUCCESS - 0 jobs")
			}
		}
		results = append(results, atlassianResult)

		leverCompanies := getLeverCompanies()
		for _, company := range leverCompanies {
			if !company.Enabled {
				continue
			}
			log.Printf("Starting sync - Lever/%s", company.Company)
			leverResult := SyncResult{Company: company.Company}
			leverRaw, err := leverfetch.FetchLeverJobs(company.LeverSlug)
			if err != nil {
				leverResult.Status = "failed"
				leverResult.Error = err.Error()
				results = append(results, leverResult)
				log.Printf("Lever/%s: FAILED - %v", company.Company, err)
				continue
			}
			leverJobs := levernormalize.NormalizeLeverJobs(leverRaw, company.Company)
			if len(leverJobs) > 0 {
				if err := embeddings.InsertJobs(leverJobs); err != nil {
					leverResult.Status = "failed"
					leverResult.Error = err.Error()
					log.Printf("Lever/%s: FAILED - insert error: %v", company.Company, err)
				} else {
					leverResult.Status = "success"
					leverResult.Count = len(leverJobs)
					totalCount += len(leverJobs)
					log.Printf("Lever/%s: SUCCESS - %d jobs", company.Company, len(leverJobs))
				}
			} else {
				leverResult.Status = "success"
				leverResult.Count = 0
				log.Printf("Lever/%s: SUCCESS - 0 jobs", company.Company)
			}
			results = append(results, leverResult)
		}

		ashbyCompanies := target.GetEnabledCompanies()
		for _, company := range ashbyCompanies {
			log.Printf("Starting sync - Ashby/%s", company.Company)
			ashbyResult := SyncResult{Company: company.Company}
			ashbyRaw, err := ashbyfetch.FetchJobBoard(company.AshbySlug)
			if err != nil {
				ashbyResult.Status = "failed"
				ashbyResult.Error = err.Error()
				results = append(results, ashbyResult)
				log.Printf("Ashby/%s: FAILED - %v", company.Company, err)
				continue
			}
			ashbyJobs := ashbynormalize.NormalizeResponse(ashbyRaw, company.Company)
			if len(ashbyJobs) > 0 {
				if err := embeddings.InsertJobs(ashbyJobs); err != nil {
					ashbyResult.Status = "failed"
					ashbyResult.Error = err.Error()
					log.Printf("Ashby/%s: FAILED - insert error: %v", company.Company, err)
				} else {
					ashbyResult.Status = "success"
					ashbyResult.Count = len(ashbyJobs)
					totalCount += len(ashbyJobs)
					log.Printf("Ashby/%s: SUCCESS - %d jobs", company.Company, len(ashbyJobs))
				}
			} else {
				ashbyResult.Status = "success"
				ashbyResult.Count = 0
				log.Printf("Ashby/%s: SUCCESS - 0 jobs", company.Company)
			}
			results = append(results, ashbyResult)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "synced successfully",
			"count":   totalCount,
			"results": results,
		})
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
