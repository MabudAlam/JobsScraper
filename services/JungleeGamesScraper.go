// services/jungleegames.go

package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func JungleeGamesScraper() ([]models.Job, error) {
	// Perform HTTP request and get HTML content
	htmlContent, err := makeHTTPRequest("https://apply.workable.com/junglee-games/?lng=en")
	if err != nil {
		log.Println("Error making HTTP request:", err)
		return nil, err
	}

	// Use goquery to parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return nil, err
	}

	var jobs []models.Job

	// Example: Extract job information
	doc.Find(".job-listing").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".job-title").Text()
		location := s.Find(".job-location").Text()
		applyURL, exists := s.Find(".apply-button").Attr("href")

		if exists {
			// Create a Job instance and append to the jobs slice
			job := models.Job{
				Title:    title,
				Location: location,
				ApplyURL: applyURL,
			}
			jobs = append(jobs, job)
		}
	})

	// Print job information (for testing purposes)
	for _, job := range jobs {
		fmt.Printf("Title: %s, Location: %s, ApplyURL: %s\n", job.Title, job.Location, job.ApplyURL)
	}

	return jobs, nil
}
