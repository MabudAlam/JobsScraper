package handlers

import (
	"goscraper/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JarHandler(c *gin.Context) {

	postings, err := services.JarScraper()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of Posting structs as a JSON response.
	c.JSON(http.StatusOK, postings)
}
