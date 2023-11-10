package handlers

import (
	"goscraper/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleScrapingRequest is a generic handler for handling scraping requests.
// It takes a Gin context and a scraper function as parameters, executes the scraper,
// and returns the scraped data as a JSON response.
//
// Parameters:
//
//	c: Gin context, representing the HTTP request and response.
//	scraper: A function that performs scraping and returns a slice of Job structs
//	         and an error if the scraping process encounters any issues.
//
// Usage:
//
//	To use this handler, pass a Gin context and a scraper function when defining your routes.
//	Example:
//	r.GET("/cred", func(c *gin.Context) {
//	    handlers.HandleScrapingRequest(c, services.CredScraper)
//	})
//	In this example, when a GET request is made to "/cred", it will trigger the
//	HandleScrapingRequest function with the CredScraper function as the scraper.
func HandleScrapingRequest(c *gin.Context, scraper func() ([]models.Job, error)) {
	// Perform scraping using the provided scraper function.
	postings, err := scraper()

	// Handle errors during scraping.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of Posting structs as a JSON response.
	c.JSON(http.StatusOK, postings)
}
