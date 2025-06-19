package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"goscraper/models"
	"goscraper/services"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error

	// Use file: prefix with proper permissions and WAL mode for better concurrency
	db, err = sql.Open("sqlite", "file:jobs.db?cache=shared&mode=rwc&_journal_mode=WAL")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite works better with limited connections
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize the database schema
	if err = initializeDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
}

func initializeDatabase() error {
	// Create the jobs table if it doesn't exist
	_, err := db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("failed to create jobs table: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

type ScraperInfo struct {
	name    string
	scraper func() ([]models.Job, error)
}

var scrapers = []ScraperInfo{
	{
		name:    "Amazon",
		scraper: services.AmazonScrapper,
	},
	{
		name:    "Atlassian",
		scraper: services.AtlassianScrapper,
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
	},
	{
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
		name:    "Google",
		scraper: services.GoogleScraper,
	},
	{
		name:    "Zoho",
		scraper: services.ZohoScraper,
	},
	{
		name:    "Paytm",
		scraper: services.PaytmScraper,
	},
	{
		name:    "Atlan",
		scraper: services.AtlanScrapper,
	},
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

type ScrapingResult struct {
	scraperName string
	jobs        []models.Job
	err         error
}

func checkDatabaseWritePermissions() error {
	// Try a simple write operation to test permissions
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS _test_write (id INTEGER);")
	if err != nil {
		return fmt.Errorf("database is not writable: %w", err)
	}

	// Clean up test table
	_, err = db.Exec("DROP TABLE IF EXISTS _test_write;")
	if err != nil {
		log.Printf("Warning: failed to clean up test table: %v", err)
	}

	return nil
}

func handleError(c *gin.Context, err error) {
	log.Printf("Handler error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func AllScrapersHandler(c *gin.Context) {
	startTime := time.Now()
	log.Println("Starting job scraping process...")

	// Check database write permissions first
	if err := checkDatabaseWritePermissions(); err != nil {
		handleError(c, fmt.Errorf("database write permission error: %w", err))
		return
	}

	// Create transaction
	tx, err := db.Begin()
	if err != nil {
		handleError(c, fmt.Errorf("failed to begin transaction: %w", err))
		return
	}

	// Ensure rollback on error
	defer func() {
		if tx != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
		}
	}()

	// Channel to collect results from all scrapers
	resultChannel := make(chan ScrapingResult, len(scrapers))
	var wg sync.WaitGroup

	// Launch all scrapers concurrently
	for _, scraperInfo := range scrapers {
		wg.Add(1)
		go func(si ScraperInfo) {
			defer wg.Done()

			scraperStart := time.Now()
			log.Printf("Starting scraper: %s", si.name)

			jobs, err := si.scraper()
			scraperDuration := time.Since(scraperStart)

			if err != nil {
				log.Printf("Scraper '%s' failed after %v: %v", si.name, scraperDuration, err)
				resultChannel <- ScrapingResult{
					scraperName: si.name,
					jobs:        nil,
					err:         err,
				}
				return
			}

			log.Printf("Scraper '%s' completed successfully in %v, found %d jobs",
				si.name, scraperDuration, len(jobs))

			resultChannel <- ScrapingResult{
				scraperName: si.name,
				jobs:        jobs,
				err:         nil,
			}
		}(scraperInfo)
	}

	// Wait for all scrapers to complete
	wg.Wait()
	close(resultChannel)

	// Collect results
	var allJobs []models.Job
	var successfulScrapers []string
	var failedScrapers []string

	for result := range resultChannel {
		if result.err != nil {
			failedScrapers = append(failedScrapers, result.scraperName)
		} else {
			successfulScrapers = append(successfulScrapers, result.scraperName)
			allJobs = append(allJobs, result.jobs...)
		}
	}

	// Log summary
	log.Printf("Scraping completed in %v", time.Since(startTime))
	log.Printf("Successful scrapers (%d): %v", len(successfulScrapers), successfulScrapers)
	if len(failedScrapers) > 0 {
		log.Printf("Failed scrapers (%d): %v", len(failedScrapers), failedScrapers)
	}
	log.Printf("Total jobs collected: %d", len(allJobs))

	// If no jobs were collected, return error
	if len(allJobs) == 0 {
		log.Println("No jobs collected from any scraper")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":               "No data found from any scraper",
			"failed_scrapers":     failedScrapers,
			"successful_scrapers": successfulScrapers,
		})
		return
	}

	// Database operations
	log.Println("Starting database operations...")

	// Clear existing data
	if _, err := tx.Exec(clearTable); err != nil {
		handleError(c, fmt.Errorf("failed to clear table: %w", err))
		return
	}

	// Ensure table exists
	if _, err := tx.Exec(createTable); err != nil {
		handleError(c, fmt.Errorf("failed to create table: %w", err))
		return
	}

	// Prepare statement for batch insert
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		handleError(c, fmt.Errorf("failed to prepare insert statement: %w", err))
		return
	}
	defer stmt.Close()

	// Insert all jobs
	insertStart := time.Now()
	insertedCount := 0

	for _, job := range allJobs {
		_, err := stmt.Exec(job.Title, job.Location, job.CreatedAt, job.Company, job.ApplyURL, job.ImageUrl)
		if err != nil {
			log.Printf("Failed to insert job (Title: %s, Company: %s): %v", job.Title, job.Company, err)
			// Continue with other jobs instead of failing completely
			continue
		}
		insertedCount++
	}

	log.Printf("Database insert completed in %v, inserted %d/%d jobs",
		time.Since(insertStart), insertedCount, len(allJobs))

	// Commit transaction
	if err := tx.Commit(); err != nil {
		handleError(c, fmt.Errorf("failed to commit transaction: %w", err))
		return
	}

	// Clear tx to prevent rollback in defer
	tx = nil

	// Return success response
	response := gin.H{
		"success":             true,
		"total_jobs":          len(allJobs),
		"inserted_jobs":       insertedCount,
		"successful_scrapers": successfulScrapers,
		"failed_scrapers":     failedScrapers,
		"execution_time_ms":   time.Since(startTime).Milliseconds(),
		"jobs":                allJobs,
	}

	c.JSON(http.StatusOK, response)
	log.Printf("Request completed successfully in %v", time.Since(startTime))
}
