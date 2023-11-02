package services

import (
	"goscraper/models"
	"goscraper/utility"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func PaytmScraper() ([]models.Job, error) {
	c := colly.NewCollector()

	url := "https://jobs.lever.co/paytm"

	var postings []models.Job

	c.OnError(func(response *colly.Response, err error) {
		log.Fatalf("Error while scraping the URL: %v", err)
	})
	currentTime := time.Now()

	c.OnHTML("div.posting", func(e *colly.HTMLElement) {
		location := e.ChildText("span.sort-by-location")
		applyURL := e.ChildAttr("a.posting-btn-submit", "href")
		jobTitle := e.ChildText("h5[data-qa=posting-name]")
		jobTitle = strings.TrimSpace(jobTitle)
		// Create a new job posting
		posting := models.Job{
			Title:     jobTitle,
			Location:  location,
			ApplyURL:  applyURL,
			ID:        utility.GenerateRandomID(),
			Company:   "Paytm",
			CreatedAt: currentTime.Unix(),
			ImageUrl:  "https://1000logos.net/wp-content/uploads/2021/03/Paytm_Logo.png",
		}
		postings = append(postings, posting)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error visiting the URL: %v", err)
	}

	return postings, nil
}
