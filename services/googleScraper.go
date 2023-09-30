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

func GoogleScraper() ([]models.Job, error) {
	c := colly.NewCollector()

	var postings []models.Job

	// Loop from page 0 to page 5
	for page := 0; page <= 10; page++ {
		//Sprintf is used to format and return a string without printing it anywhere
		//similar to printf("%d",page")
		url := fmt.Sprintf("https://www.google.com/about/careers/applications/jobs/results/?location=Bengaluru%%2C%%20Karnataka%%2C%%20India&jlo=en-US&hl=en&page=%d", page)

		c.OnError(func(response *colly.Response, err error) {
			log.Fatalf("Error while scraping the URL: %v", err)
		})
		currentTime := time.Now()

		c.OnHTML("div.sMn82b", func(e *colly.HTMLElement) {
			posting := models.Job{
				Title:     strings.TrimSpace(e.ChildText("h3.QJPWVe")),
				Location:  strings.TrimSpace(e.ChildText("span.r0wTof")),
				ApplyURL:  "https://www.google.com/about/careers/applications/" + e.ChildAttr("a", "href"),
				ID:        utility.GenerateRandomID(),
				Company:   "Google",
				CreatedAt: currentTime.Unix(),
				ImageUrl:  "https://www.freepnglogos.com/uploads/google-logo-png/google-logo-png-suite-everything-you-need-know-about-google-newest-0.png",
			}

			// Check if the posting is not already in the 'postings' slice before appending
			if !utility.CheckDuplicates(postings, posting) {
				postings = append(postings, posting)
			}
		})

		err := c.Visit(url)
		if err != nil {
			log.Fatalf("Error visiting the URL: %v", err)
		}
	}

	return postings, nil
}
