package services

import (
	"encoding/json"
	"fmt"
	"goscraper/models"
	"net/http"
	"time"
)

func MplScrapper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://mpl.darwinbox.in/ms/candidateapi/job?page=1&limit=100"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a slice to hold the item data from the API.
	var amazonpostings models.MPL

	// Decode the JSON response into the amazonpostings struct.
	if err := json.NewDecoder(resp.Body).Decode(&amazonpostings); err != nil {
		fmt.Print(err)
		return nil, err
	}
	// Create a slice to hold the job postings.
	postings := make([]models.Job, len(amazonpostings.Message.Jobs))
	currentTime := time.Now()
	// Iterate over the job postings and convert them to models.Job structs.
	for i, posting := range amazonpostings.Message.Jobs {
		postings[i] = models.Job{
			Title:     posting.Title,
			ID:        posting.ID,
			Location:  string(posting.OfficelocationShowArr),
			CreatedAt: currentTime.Unix(),
			Company:   "MPL",
			ApplyURL:  fmt.Sprintf("https://mpl.darwinbox.in/ms/candidate/careers/%s", posting.ID),

			ImageUrl: "https://s3-ap-southeast-1.amazonaws.com/darwinbox-data/INSTANCE1_5ec3db08d9b48_318/logo/a60dc6d85d1baa__tenant-avatar-318_1584778816.png",
		}
	}

	return postings, nil
}
