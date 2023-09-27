package main

import (

	"goscraper/handlers"
	"log"
	"net/http"
	

	"github.com/gin-gonic/gin"
)

type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniversalDomain         string `json:"universal_domain"`
}

func main() {

	
	r := gin.Default()

	r.GET("/cred", handlers.GetPostingsHandler)
	r.GET("/atlassian", handlers.AtlassianHandler)
	r.GET("/amazon", handlers.Amazonhandler)
	r.GET("/coursera", handlers.CourseraHandler)
	r.GET("/freshworks", handlers.FreshWorksHandler)
	r.GET("/gojek", handlers.Gojekhandler)
	r.GET("/mpl", handlers.MplHandler)

	r.GET("/syncwithSql", handlers.AllScrapersHandler)
	r.GET("/getallJobsFromSQL", handlers.GetAllJobsFromSqlite)

	
	r.GET("/", func(c *gin.Context) {
		c.File("static/base.html")
	})
	r.Use(CORSMiddleware())

	if err := r.Run(":" + "8080"); err != nil {
		log.Panicf("error: %s", err)
	}
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

// RootHandler handles the root route (/).
func RootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"hello": "Mabud"})
}
