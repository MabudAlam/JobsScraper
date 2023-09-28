package services

import (
	"encoding/json"
	"fmt"
	"goscraper/models"
	"net/http"
	"time"
)

func ZohoScraper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://careers.zohocorp.com/recruit/v2/public/Job_Openings?pagename=Careers&source=CareerSite&extra_fields=%5B%22Remote_Job%22%5D"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a slice to hold the item data from the API.
	var zohoPostings models.ZohoJobs

	// Decode the JSON response into the amazonpostings struct.
	if err := json.NewDecoder(resp.Body).Decode(&zohoPostings); err != nil {
		fmt.Print(err)
		return nil, err
	}
	// Create a slice to hold the job postings.
	postings := make([]models.Job, len(zohoPostings.Data))
	currentTime := time.Now()
	// Iterate over the job postings and convert them to models.Job structs.
	for i, posting := range zohoPostings.Data {
		postings[i] = models.Job{
			Title:     posting.PostingTitle,
			ID:        posting.ID,
			Location:  posting.Country1,
			CreatedAt: currentTime.Unix(),
			Company:   "Zoho",
			ApplyURL:  posting.URL,
			ImageUrl:  "https://www.zohowebstatic.com/sites/default/files/zoho-logo-zh-2x.png",
		}
	}

	return postings, nil
}
