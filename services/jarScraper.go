package services

import (
	"encoding/json"
	"errors"
	"goscraper/models"
	"goscraper/utility"
	"io/ioutil"
	"net/http"
	"time"
)

// Helper functions moved to a separate file for clarity and maintainability.

func JarScraper() ([]models.Job, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.pyjamahr.com/api/public/jobs/", nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Authorization", "Token b28b5dbc825124022c8638b6a28ebaea9d9b9a28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch data from API")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse models.JarAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	var postings []models.Job
	currentTime := time.Now()

	for _, item := range apiResponse.Results {
		desc := buildJarDescription(item)
		posting := models.Job{
			Title:       item.Title,
			Location:    item.Location,
			ApplyURL:    item.ApplyLink,
			Description: desc,
			ID:          utility.GenerateRandomID(),
			Company:     item.CompanyName,
			CreatedAt:   currentTime.Unix(),
			ImageUrl:    "https://lever-client-logos.s3.us-west-2.amazonaws.com/2a81954c-124c-4e43-a34f-916adc6d3e9a-1694765296054.png",
		}
		if !utility.CheckDuplicates(postings, posting) {
			postings = append(postings, posting)
		}
	}

	return postings, nil
}
