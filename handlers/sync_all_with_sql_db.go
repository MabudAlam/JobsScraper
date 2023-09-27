package handlers

import (
	"goscraper/config"
	"goscraper/models"
	"goscraper/services"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Defined a list of scraper functions
var allscrapers = []struct {
	name    string
	scraper func() ([]models.Job, error)
}{
	{
		name:    "Amazon",
		scraper: services.AmazonScrapper,
	},
	{
		name:    "Atlassian",
		scraper: services.AtlassianScrapper,
	},
	{
		name:    "Coursera",
		scraper: services.CourseraScraper,
	},
	{
		name:    "CRED",
		scraper: services.CredScraper,
	}, {
		name:    "FreshWorks",
		scraper: services.FreshWorksScraper,
	},
	{
		name:    "Gojek",
		scraper: services.GojekScraper,
	},
	{
		name:    "MPL",
		scraper: services.MplScrapper,
	},
}

func SyncAll(c *gin.Context) {
	// Get the GORM database instance from the config package.
	config.Connect()
	db := config.DB

	// Clear existing job postings from the database.
	clearTable := `
        DELETE FROM gorm_dbs;
    `
	if err := db.Exec(clearTable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create an empty slice to hold all the job postings.
	var allPostings []models.Job

	// Loop through the list of scrapers and call each one.
	for _, scraper := range scrapers {
		postings, err := scraper.scraper()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		allPostings = append(allPostings, postings...)
	}

	// Convert the Job slice to GormDB slice.
	var gormPostings []models.GormDB
	for _, posting := range allPostings {
		gormPosting := models.ConvertToGormDB(posting)
		gormPostings = append(gormPostings, gormPosting)
	}

	// Insert all GormDB records into the database in a single batch.
	if err := db.Create(&gormPostings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the aggregated job postings as a JSON response.
	c.JSON(http.StatusOK, allPostings)
}
