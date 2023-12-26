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

		var locationValue = obj.Location[0]
		if len(obj.Location) > 1 {
			locationValue = "Many Locations"
		}

		postings[i] = models.Job{
			Title:     obj.Title,
			ID:        obj.LeverId,
			Location:  locationValue,
			CreatedAt: time.Now().Unix(),
			Company:   "Atlassian",
			ApplyURL:  obj.ApplyUrl,
			ImageUrl:  "https://seeklogo.com/images/A/atlassian-logo-DF2FCF6E4D-seeklogo.com.png",
		}
	}

	return postings, nil
}
