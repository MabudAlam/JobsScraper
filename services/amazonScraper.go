package services

import (
	"encoding/json"
	"fmt"
	"goscraper/models"
	"net/http"
	"time"
)

func AmazonScrapper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://www.amazon.jobs/en/search.json?country=IND&result_limit=100"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a slice to hold the item data from the API.
	var amazonpostings models.Amazon

	// Decode the JSON response into the amazonpostings struct.
	if err := json.NewDecoder(resp.Body).Decode(&amazonpostings); err != nil {
		fmt.Print(err)
		return nil, err
	}
	// Create a slice to hold the job postings.
	postings := make([]models.Job, len(amazonpostings.Jobs))
	currentTime := time.Now()
	// Iterate over the job postings and convert them to models.Job structs.
	for i, posting := range amazonpostings.Jobs {
		postings[i] = models.Job{
			Title:     posting.Title,
			ID:        posting.ID,
			Location:  posting.Location,
			CreatedAt: currentTime.Unix(),
			Company:   "Amazon",
			ApplyURL:  "https://www.amazon.jobs" + posting.JobPath,
			ImageUrl:  "https://upload.wikimedia.org/wikipedia/commons/thumb/4/4a/Amazon_icon.svg/2500px-Amazon_icon.svg.png",
		}
	}

	return postings, nil
}
