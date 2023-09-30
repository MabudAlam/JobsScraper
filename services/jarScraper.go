package services

import (
	"goscraper/models"
	"goscraper/utility"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func JarScraper() ([]models.Job, error) {
	c := colly.NewCollector()

	var postings []models.Job

	url := "https://jobs.lever.co/jar-app"

	c.OnError(func(response *colly.Response, err error) {
		log.Fatalf("Error while scraping the URL: %v", err)
	})

	currentTime := time.Now()

	c.OnHTML(".posting", func(e *colly.HTMLElement) {
		posting := models.Job{
			Title:     strings.TrimSpace(e.ChildText("h5[data-qa='posting-name']")),
			Location:  strings.TrimSpace(e.ChildText(".location")),
			ApplyURL:  e.ChildAttr("a.posting-btn-submit", "href"),
			ID:        utility.GenerateRandomID(),
			Company:   "Jar App",
			CreatedAt: currentTime.Unix(),
			ImageUrl:  "https://lever-client-logos.s3.us-west-2.amazonaws.com/2a81954c-124c-4e43-a34f-916adc6d3e9a-1694765296054.png",
		}

		if !utility.CheckDuplicates(postings, posting) {
			postings = append(postings, posting)
		}
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error visiting the URL: %v", err)
	}

	return postings, nil
}
