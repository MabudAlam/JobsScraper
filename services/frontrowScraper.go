package services

import (
	"encoding/json"
	"goscraper/models"
	"net/http"
	"strconv"
	"time"
)

func FrontRowScrapper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://frontrow.keka.com/careers/api/embedjobs/active/1b7199b8-4d91-4b16-a039-75937a62aa3a"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a slice to hold the item data from the API.
	var amazonpostings models.FrontRow

	// Decode the JSON response into the amazonpostings struct.
	if err := json.NewDecoder(resp.Body).Decode(&amazonpostings); err != nil {
		return nil, err
	}
	// Create a slice to hold the job postings.
	postings := make([]models.Job, len(amazonpostings))
	currentTime := time.Now()
	// Iterate over the job postings and convert them to models.Job structs.
	for i, posting := range amazonpostings {
		postings[i] = models.Job{
			Title:     posting.Title,
			ID:        strconv.Itoa(int(posting.ID)),
			Location:  posting.JobLocations[0].Name,
			CreatedAt: currentTime.Unix(),
			Company:   "FrontRow",
			ApplyURL:  "https://frontrow.kekahire.com/jobdetails/" + strconv.Itoa(int(posting.ID)),
			ImageUrl:  "https://stkekahireprodcin01.blob.core.windows.net/1b7199b8-4d91-4b16-a039-75937a62aa3a/orglogo/7a6e97c8a65b44a69833bbacb30e737b.png",
		}
	}

	return postings, nil
}
