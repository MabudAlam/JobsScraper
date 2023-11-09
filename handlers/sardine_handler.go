package handlers

import (
	"goscraper/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SardineScraper(c *gin.Context) {
	postings, err := services.SardineScraper()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, postings)
}
