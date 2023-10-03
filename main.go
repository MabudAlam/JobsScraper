package main

import (
	"goscraper/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/cred", handlers.GetPostingsHandler)
	r.GET("/atlassian", handlers.AtlassianHandler)
	r.GET("/amazon", handlers.Amazonhandler)
	r.GET("/coursera", handlers.CourseraHandler)
	r.GET("/freshworks", handlers.FreshWorksHandler)
	r.GET("/gojek", handlers.Gojekhandler)
	r.GET("/mpl", handlers.MplHandler)
	r.GET("/google", handlers.GoogleHandler)
	r.GET("/fi", handlers.EpifiHandler)

	r.GET("/zoho", handlers.ZohoHandler)
	r.GET("/jar", handlers.JarHandler)

	r.GET("/syncwithSql", handlers.AllScrapersHandler)
	r.GET("/getallJobsFromSQL", handlers.GetAllJobsFromSqlite)

	r.GET("/", func(c *gin.Context) {
		c.File("static/base.html")
	})

	if err := r.Run(":" + "8080"); err != nil {
		log.Panicf("error: %s", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from these origins
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080, https://jobs-scraper-production.up.railway.app")
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

// RootHandler handles the root route (/).
func RootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"hello": "Mabud"})
}
