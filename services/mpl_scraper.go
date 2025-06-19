package services

import (
	"encoding/json"
	"fmt"
	"goscraper/models"
	"io"
	"net/http"
	"time"
)

func MplScrapper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://mpl.darwinbox.in/ms/candidateapi/job?page=1&limit=100"

	client := &http.Client{}
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Set headers to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://mpl.darwinbox.in/")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://mpl.darwinbox.in")

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mplJobPostings models.MPLJobData
	if err := json.Unmarshal(bodyBytes, &mplJobPostings); err != nil {
		return nil, err
	}

	postings := make([]models.Job, len(mplJobPostings.Message.Jobs))
	currentTime := time.Now()
	for i, posting := range mplJobPostings.Message.Jobs {
		postings[i] = models.Job{
			Title:     posting.Title,
			ID:        posting.ID,
			Location:  string(posting.OfficelocationArr),
			CreatedAt: currentTime.Unix(),
			Company:   "MPL",
			ApplyURL:  fmt.Sprintf("https://mpl.darwinbox.in/ms/candidate/careers/%s", posting.ID),
			ImageUrl:  "https://s3-ap-southeast-1.amazonaws.com/darwinbox-data/INSTANCE1_5ec3db08d9b48_318/logo/a60dc6d85d1baa__tenant-avatar-318_1584778816.png",
		}
	}

	return postings, nil
}
