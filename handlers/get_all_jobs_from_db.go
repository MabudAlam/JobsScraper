package handlers

import (
	"goscraper/config"
	"goscraper/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetJobsFromDB(c *gin.Context) {
	// Get the GORM database instance from the config package.

	config.Connect()
	db := config.DB

	// Create a slice to store the retrieved jobs.
	var jobs []models.GormDB

	// Query the database to retrieve all jobs.
	if err := db.Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}
