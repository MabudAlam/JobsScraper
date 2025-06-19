package services

import (
	"goscraper/models"
	"goscraper/utility"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func EpfiScraper() ([]models.Job, error) {
	c := colly.NewCollector()

	url := "https://jobs.lever.co/epifi"

	var postings []models.Job

	handleError := func(err error) {
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	c.OnError(func(response *colly.Response, err error) {
		handleError(err)
	})

	currentTime := time.Now()
	c.OnHTML(".posting", func(e *colly.HTMLElement) {
		posting := models.Job{
			Title:     strings.TrimSpace(e.ChildText("h5[data-qa='posting-name']")),
			Location:  strings.TrimSpace(e.ChildText("span.location")),
			ApplyURL:  e.ChildAttr("a.posting-btn-submit", "href"),
			ID:        utility.GenerateRandomID(),
			Company:   "Fi Money",
			CreatedAt: currentTime.Unix(),
			ImageUrl:  "https://lever-client-logos.s3.us-west-2.amazonaws.com/b1a67862-9e17-4f38-896f-347e051e94b3-1612511026673.png",
		}
		postings = append(postings, posting)
	})

	err := c.Visit(url)
	handleError(err)

	return postings, nil
}
