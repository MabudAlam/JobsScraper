package services

import (
	"encoding/json"
	"fmt"
	"goscraper/models"
	"io"
	"net/http"
	"time"
)

func AtlassianScrapper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://www.atlassian.com/endpoint/careers/listings"

	// Make an HTTP GET request to the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resp.Body)

	// Create a slice to hold the item data from the API.
	var atlassianpostings []models.Posting

	// Decode the JSON response into the atlassianpostings struct.
	jsonByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonByte, &atlassianpostings)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold the atlassianpostings postings.
	postings := make([]models.Job, len(atlassianpostings))

	// Iterate over the atlassianpostings and convert them to models.Job structs.
	for i, obj := range atlassianpostings {
		locationValue := ""
		if len(obj.Location) > 0 {
			locationValue = obj.Location[0]
		}
		postings[i] = models.Job{
			Title:       obj.Title,
			ID:          obj.LeverId,
			Location:    locationValue,
			Description: obj.Overview + "\n" + obj.Responsibilities + "\n" + obj.Qualifications,
			CreatedAt:   time.Now().Unix(),
			Company:     "Atlassian",
			ApplyURL:    obj.ApplyUrl,
			ImageUrl:    "https://res.cloudinary.com/dc0tfxkph/image/upload/v1704286280/Go%20Scraper/u4eaxsknspamducttdh3.png",
		}
	}

	return postings, nil
}
