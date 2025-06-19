package services

import (
	"goscraper/models"
	"goscraper/utility"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func FreshWorksScraper() ([]models.Job, error) {
	c := colly.NewCollector()

	url := "https://careers.smartrecruiters.com/Freshworks"

	var postings []models.Job
	// c.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*",
	// 	Parallelism: 2,
	// 	Delay:       2 * time.Second,
	// })

	c.OnError(func(response *colly.Response, err error) {
		log.Fatalf("Error while scraping the URL: %v", err)
	})
	currentTime := time.Now()

	c.OnHTML("section.openings-section", func(e *colly.HTMLElement) {
		posting := models.Job{
			Title: strings.TrimSpace(e.ChildText("h4.details-title")),

			Location: strings.TrimSpace(e.ChildText("h3.opening-title")),
			ApplyURL: e.ChildAttr("a", "href"),

			ID:        utility.GenerateRandomID(),
			Company:   "Freshworks",
			CreatedAt: currentTime.Unix(),
			ImageUrl:  "https://res.cloudinary.com/dc0tfxkph/image/upload/v1750323820/Go%20Scraper/ablsfwdo8uyqmb1a46tl.png",
		}
		postings = append(postings, posting)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error visiting the URL: %v", err)
	}

	return postings, nil
}
