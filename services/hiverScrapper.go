package services

import (
	"fmt"
	"goscraper/models"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// HiverScraper scrapes job postings from the provided Hiver link.
func HiverScraper() ([]models.Job, error) {
	// URL of the Hiver job listings page.
	url := "https://jobs.lever.co/hiver"

	// Create a new collector
	c := colly.NewCollector()

	// Create a slice to hold job postings
	var hiverPostings []models.Job

	// Set up a callback to be executed when a job posting is found
	c.OnHTML(".posting", func(e *colly.HTMLElement) {
		// Extract relevant information
		title := e.ChildText(".posting-title")
		location := e.ChildText(".posting-categories span") // Adjust this based on HTML structure
		description := e.ChildText(".posting-about")

		// Print or store the information as needed
		fmt.Printf("Title: %s\nLocation: %s\nDescription: %s\n", title, location, description)

		// Create a new job posting
		hiverPostings = append(hiverPostings, models.Job{
			Title:     title,
			Location:  location,
			CreatedAt: time.Now().Unix(),
			Company:   "Hiver",
			// Add other fields as needed
		})
	})

	// Visit the Hiver job listings page
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return hiverPostings, nil
}
