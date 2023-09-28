package services

import (
	"encoding/json"
	"fmt"
	"goscraper/models"
	"goscraper/utility"
	"net/http"
	"time"
)

func PayPalScraper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://paypal.eightfold.ai/api/apply/v2/jobs?Country=india&num=10"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a slice to hold the item data from the API.
	var paypalPostings models.PayPalJobs

	// Decode the JSON response into the amazonpostings struct.
	if err := json.NewDecoder(resp.Body).Decode(&paypalPostings); err != nil {
		fmt.Print(err)
		return nil, err
	}
	// Create a slice to hold the job postings.
	postings := make([]models.Job, len(paypalPostings.Positions))
	currentTime := time.Now()
	// Iterate over the job postings and convert them to models.Job structs.
	for i, posting := range paypalPostings.Positions {
		postings[i] = models.Job{
			Title:     posting.Name,
			ID:        utility.GenerateRandomID(),
			Location:  string(posting.Location),
			CreatedAt: currentTime.Unix(),
			Company:   "PayPal",
			ApplyURL:  posting.CanonicalPositionURL,
			ImageUrl:  "https://upload.wikimedia.org/wikipedia/commons/a/a4/Paypal_2014_logo.png",
		}
	}

	return postings, nil
}
