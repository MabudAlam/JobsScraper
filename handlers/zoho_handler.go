package handlers

import (
	"goscraper/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ZohoHandler(c *gin.Context) {

	postings, err := services.ZohoScraper()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the list of Posting structs as a JSON response.
	c.JSON(http.StatusOK, postings)
}
