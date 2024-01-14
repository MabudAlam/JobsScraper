package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"goscraper/models"
	"goscraper/services"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", "jobs.db")
	if err != nil {
		log.Fatal(err)
	}
}

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
		name:    "Hiver",
		scraper: services.HiverScraper,
	},
	{
		name:    "Niyo Solutions",
		scraper: services.NiyoSolutionScraper,
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
	}, {
		name:    "Fincent",
		scraper: services.FincentScraper,
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

const clearTable = `DELETE FROM jobs;`

const createTable = `
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

const insertSQL = `
	INSERT INTO jobs (title, location, created_at, company, apply_url, image_url)
	VALUES (?, ?, ?, ?, ?, ?);
`

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func AllScrapersHandler(c *gin.Context) {
	tx, err := db.Begin()
	if err != nil {
		handleError(c, err)
		return
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println("Rollback failed:", err)
		}
	}()

	var allPostings []models.Job

	for _, scraper := range scrapers {
		postings, err := scraper.scraper()
		if err != nil {
			handleError(c, err)
			return
		}
		allPostings = append(allPostings, postings...)
	}

	if len(allPostings) > 0 {
		_, err := tx.Exec(clearTable)
		if err != nil {
			handleError(c, err)
			return
		}

		_, err = tx.Exec(createTable)
		if err != nil {
			handleError(c, err)
			return
		}

		for _, posting := range allPostings {
			_, err := tx.Exec(insertSQL, posting.Title, posting.Location, posting.CreatedAt, posting.Company, posting.ApplyURL, posting.ImageUrl)
			if err != nil {
				handleError(c, err)
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, allPostings)
		log.Printf("Inserted %d jobs into the database", len(allPostings))
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No data found"})
	}
}
