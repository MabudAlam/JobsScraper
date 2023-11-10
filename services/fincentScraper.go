package services

import (
	"fmt"
	"goscraper/models"
	"goscraper/utility"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func FincentScraper() ([]models.Job, error) {
	c := colly.NewCollector()

	url := "https://boards.greenhouse.io/fincent"

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
	c.OnHTML("div.opening", func(e *colly.HTMLElement) {
		// Extracting data from the HTML structure
		title := e.ChildText("a[data-mapped=true]")
		location := strings.TrimSpace(e.ChildText("span.location"))
		applyURL := e.ChildAttr("a[data-mapped=true]", "href")

		posting := models.Job{
			Title:     title,
			Location:  location,
			ApplyURL:  "https://boards.greenhouse.io" + applyURL,
			ID:        utility.GenerateRandomID(),
			Company:   "Fincent",
			CreatedAt: currentTime.Unix(),
			ImageUrl:  "https://res.cloudinary.com/dc0tfxkph/image/upload/v1699633987/Go%20Scraper/btaabv0gjex1erlcnkax.jpg",
		}

		fmt.Print(posting)
		postings = append(postings, posting)
	})

	err := c.Visit(url)
	handleError(err)

	return postings, nil
}
