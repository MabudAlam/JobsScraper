package services

import (
	"goscraper/models"
	"goscraper/utility"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func SardineScraper() ([]models.Job, error) {
	c := colly.NewCollector()

	url := "https://www.sardine.ai/careers#open-positions"

	var postings []models.Job

	c.OnError(func(response *colly.Response, err error) {
		log.Fatalf("Error while scraping the URL: %v", err)
	})
	currentTime := time.Now()

	c.OnHTML("div.job-lists-layout", func(e *colly.HTMLElement) {
		location := e.ChildText("span.text-size-small")
		applyURL := e.ChildAttr("a", "href")
		jobTitle := e.ChildText("div.text-size-regular")
		jobTitle = strings.TrimSpace(jobTitle)
		posting := models.Job{
			Title:     jobTitle,
			Location:  location,
			ApplyURL:  applyURL,
			ID:        utility.GenerateRandomID(),
			Company:   "Sardine",
			CreatedAt: currentTime.Unix(),
			ImageUrl:  "https://res.cloudinary.com/dc0tfxkph/image/upload/v1699560212/Go%20Scraper/yellbefk4skxwlouoznd.png",
		}
		postings = append(postings, posting)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error visiting the URL: %v", err)
	}

	return postings, nil
}
