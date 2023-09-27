package services

import (
	"encoding/json"
	"fmt"
	"net/http"
"goscraper/models")

func AtlassianScrapper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://www.atlassian.com/.rest/postings"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a slice to hold the item data from the API.
	var atlassianpostins models.Temperatures

	// Decode the JSON response into the atlassianpostins struct.
	if err := json.NewDecoder(resp.Body).Decode(&atlassianpostins); err != nil {
		fmt.Print(err)
		return nil, err
	}

	// Create a slice to hold the job postings.
	postings := make([]models.Job, len(atlassianpostins.Postings))

	// Iterate over the job postings and convert them to models.Job structs.
	for i, posting := range atlassianpostins.Postings {
		postings[i] = models.Job{
			Title:     posting.Text,
			ID:        posting.ID,
			Location:  posting.Categories.Location,
			CreatedAt: posting.CreatedAt,
			Company:   "Atlassian",
			ApplyURL:  posting.Urls.ApplyURL,
			ImageUrl:  "https://seeklogo.com/images/A/atlassian-logo-DF2FCF6E4D-seeklogo.com.png",
		}
	}

	return postings, nil
}
