package handlers

import (
	"github.com/gin-gonic/gin"
	"goscraper/services"
	"net/http"
)

// GetPostingsHandler handles the /get-postings route.
func Amazonhandler(c *gin.Context) {
	// Call fetchPostings to retrieve the data.
	postings, err := services.AmazonScrapper()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the list of Posting structs as a JSON response.
	c.JSON(http.StatusOK, postings)
}
