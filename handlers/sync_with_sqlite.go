package handlers

import (
	"database/sql"
	"goscraper/models"
	"goscraper/services"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

// Defined a list of scraper functions
var scrapers = []struct {
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
	},
	{
		name:    "FreshWorks",
		scraper: services.FreshWorksScraper,
	},
	{
		name:    "Fi Money",
		scraper: services.EpfiScraper,
	},
	{
		name:    "FrontRow",
		scraper: services.FrontRowScrapper,
	},
	{
		name:    "Gojek",
		scraper: services.GojekScraper,
	},
	{
		name:    "MPL",
		scraper: services.MplScrapper,
	},
	{
		name:    "Google",
		scraper: services.GoogleScraper,
	},
	{
		name:    "Zoho",
		scraper: services.ZohoScraper,
	},
	{
		name:    "Jar App",
		scraper: services.JarScraper,
	},
	{
		name:    "Paytm",
		scraper: services.PaytmScraper,
	},
	{
		name:    "Sardine",
		scraper: services.SardineScraper,
	},
	// {
	// 	name:    "PayPal",
	// 	scraper: services.PayPalScraper,
	// },
}

func AllScrapersHandler(c *gin.Context) {
	// Open or create the SQLite database file.
	db, err := sql.Open("sqlite", "jobs.db") // Use "sqlite" as the driver name
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

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

	if allPostings != nil {
		clearTable := `
		DELETE FROM jobs;
	`
		_, err = db.Exec(clearTable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		createTable := `
			CREATE TABLE IF NOT EXISTS jobs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT,
				location TEXT,
				created_at INTEGER,
				company TEXT,
				apply_url TEXT,
				image_url TEXT
			);
		`
		_, err = db.Exec(createTable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Loop through the job postings and insert them into the database.
		for _, posting := range allPostings {
			insertSQL := `
				INSERT INTO jobs (title, location, created_at, company, apply_url, image_url)
				VALUES (?, ?, ?, ?, ?, ?);
			`
			_, err := db.Exec(insertSQL, posting.Title, posting.Location, posting.CreatedAt, posting.Company, posting.ApplyURL, posting.ImageUrl)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, allPostings)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No data found"})
		return
	}

	// Return the aggregated job postings as a JSON response.

}
