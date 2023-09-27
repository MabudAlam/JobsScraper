package services

import (
	"fmt"
	"goscraper/models"
	"log"
	"goscraper/utility"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func GojekScraper() ([]models.Job, error) {
	c := colly.NewCollector()

	url := "https://jobs.lever.co/GoToGroup"

	var postings []models.Job

	c.OnError(func(response *colly.Response, err error) {
		log.Fatalf("Error while scraping the URL: %v", err)
	})
	currentTime := time.Now()

	c.OnHTML("div.posting", func(e *colly.HTMLElement) {
		location := e.ChildText("span.sort-by-location")
		applyURL := e.ChildAttr("a", "href")
		jobTitle := e.DOM.Find("h5[data-qa='posting-name']").Text()
		jobTitle = strings.TrimSpace(jobTitle)
		// Create a new job posting
		posting := models.Job{
			Title:     jobTitle,
			Location:  location,
			ApplyURL:  applyURL,
			ID:        utility.GenerateRandomID(),
			Company:   "GoJek",
			CreatedAt: currentTime.Unix(),
			ImageUrl:  "https://lever-client-logos.s3.us-west-2.amazonaws.com/ea3dc26a-36d6-4989-8596-ac871e4a9e82-1682066441676.png",
		}
		fmt.Print(posting)
		postings = append(postings, posting)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error visiting the URL: %v", err)
	}

	// // Print the scraped job postings
	// for i, posting := range postings {
	// 	fmt.Printf("Job #%d\n", i+1)
	// 	fmt.Printf("Title: %s\n", posting.Title)
	// 	fmt.Printf("Location: %s\n", posting.Location)
	// 	fmt.Printf("URL: %s\n", posting.ApplyURL)
	// 	fmt.Println()
	// }

	return postings, nil
}
