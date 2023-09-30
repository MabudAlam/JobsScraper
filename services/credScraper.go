package services

import (
	"encoding/json"
	"goscraper/models"
	"net/http"
)

// Function to fetch postings from the API and convert them to Posting structs.
func CredScraper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://api.lever.co/v0/postings/cred?mode=json"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a slice to hold the CRED data from the API.
	var credData []models.CRED

	// Decode the JSON response into the credData slice.
	if err := json.NewDecoder(resp.Body).Decode(&credData); err != nil {
		return nil, err
	}

	// Convert CRED data to Posting structs.
	postings := make([]models.Job, len(credData))
	for i, cred := range credData {
		postings[i] = models.Job{
			Title:     cred.Text,
			ID:        cred.ID,
			Location:  cred.Categories.Location,
			CreatedAt: cred.CreatedAt,
			Company:   "CRED",
			ApplyURL:  cred.ApplyURL,

			ImageUrl: "https://play-lh.googleusercontent.com/r2ZbsIr5sQ7Wtl1T6eevyWj4KS7QbezF7JYB9gxQnLWbf0K4kU7qaLNcJLLUh0WG-3pK",
		}
	}

	return postings, nil
}
