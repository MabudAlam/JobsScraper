package main

import (
	"goscraper/handlers"
	"goscraper/models"
	"goscraper/services"
	"log"
	"net/http"
	"os"

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
	r.GET("/cred", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.CredScraper)
	})
	r.GET(("/atlassian"), func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.AtlassianScrapper)
	})
	r.GET("/amazon", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.AmazonScrapper)
	})
	r.GET("/coursera", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.CourseraScraper)
	})
	r.GET("/freshworks", func(c *gin.Context) {

		handlers.HandleScrapingRequest(c, services.FreshWorksScraper)
	})
	r.GET("/gojek", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.GojekScraper)
	})
	r.GET("/mpl", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.MplScrapper)
	})
	r.GET("/google", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.GoogleScraper)
	})
	r.GET("/fi", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.EpfiScraper)
	})
	r.GET("/frontrow", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.FrontRowScrapper)
	})
	r.GET("/sardine", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.SardineScraper)
	})
	r.GET("/zoho", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.ZohoScraper)

	})
	r.GET("/jar", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.JarScraper)
	})
	r.GET("/paytm", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.PaytmScraper)
	})
	r.GET("/fincent", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.FincentScraper)
	})
	r.GET("/paypal", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.PayPalScraper)
	})
	r.GET("/niyo", func(c *gin.Context) {
		handlers.HandleScrapingRequest(c, services.NiyoSolutionScraper)
	})

	r.POST("/insertCustomJob", func(c *gin.Context) {
		var input models.UserDefinedJob

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		job := models.UserDefinedJob{
			Title:       input.Title,
			ID:          input.ID,
			Location:    input.Location,
			ApplyURL:    input.ApplyURL,
			ImageUrl:    input.ImageUrl,
			CreatedAt:   input.CreatedAt,
			Company:     input.Company,
			Expired:     input.Expired,
			Salary:      input.Salary,
			Skills:      input.Skills,
			Description: input.Description,
		}

		handlers.InsertJob(c, job)

	})
	r.POST("/updateCustomJob", func(c *gin.Context) {
		var input models.UserDefinedJob

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		job := models.UserDefinedJob{
			Title:       input.Title,
			ID:          input.ID,
			Location:    input.Location,
			ApplyURL:    input.ApplyURL,
			ImageUrl:    input.ImageUrl,
			CreatedAt:   input.CreatedAt,
			Company:     input.Company,
			Expired:     input.Expired,
			Salary:      input.Salary,
			Skills:      input.Skills,
			Description: input.Description,
		}

		handlers.UpdateJob(c, job)

	})
	r.POST("/deleteCustomJob", func(c *gin.Context) {
		var input models.UserDefinedJob

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		handlers.DeleteJob(c, input.ID)

	})
	r.GET("/getallCustomJobs", handlers.GetAllJobs)
	r.GET("/syncwithSql", func(c *gin.Context) {
		password := c.Query("password")

		correctPassword := os.Getenv("SYNC_WITH_SQL_PASSWORD")

		if password != correctPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		handlers.AllScrapersHandler(c)
	})

	r.GET("/getallJobsFromSQL", handlers.GetAllJobsFromSqlite)

	r.GET("/", func(c *gin.Context) {
		c.File("static/base.html")
	})

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Host)
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
