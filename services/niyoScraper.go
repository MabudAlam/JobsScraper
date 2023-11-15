package services

import (
	"encoding/json"
	"fmt"
	"goscraper/models"
	"net/http"
	"time"
)

func NiyoSolutionScraper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://niyo.darwinbox.in/ms/candidateapi/job?page=1&limit=50"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a slice to hold the item data from the API.
	var niyo models.NiyoSolution

	// Decode the JSON response into the amazonpostings struct.
	if err := json.NewDecoder(resp.Body).Decode(&niyo); err != nil {
		fmt.Print(err)
		return nil, err
	}
	// Create a slice to hold the job postings.
	postings := make([]models.Job, len(niyo.Message.Jobs))
	currentTime := time.Now()
	// Iterate over the job postings and convert them to models.Job structs.
	for i, posting := range niyo.Message.Jobs {
		postings[i] = models.Job{
			Title:     posting.Title,
			ID:        posting.ID,
			Location:  string(posting.OfficelocationShowArr),
			CreatedAt: currentTime.Unix(),
			Company:   "Niyo Solutions",
			ApplyURL:  "https://niyo.darwinbox.in/ms/candidate/careers/" + posting.ID,

			ImageUrl: "https://s3-ap-southeast-1.amazonaws.com/darwinbox-data/INSTANCE5_a63309506e450f_234/logo/a635bf69b93798__tenant-avatar-234_350882374.JPG",
		}
	}

	return postings, nil
}
