package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxBodyBytes = 10 * 1024 * 1024

type RawJob struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Location           string `json:"location"`
	NormalizedLocation string `json:"normalized_location"`
	Description        string `json:"description"`
	DescriptionShort   string `json:"description_short"`
	JobPath            string `json:"job_path"`
}

type AmazonResponse struct {
	Jobs []RawJob `json:"jobs"`
}

func FetchAmazonJobs() (*AmazonResponse, error) {
	url := "https://www.amazon.jobs/en/search.json?country=IND&result_limit=100"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	limitedBody := http.MaxBytesReader(nil, resp.Body, maxBodyBytes)
	body, err := io.ReadAll(limitedBody)
	if err != nil {
		return nil, fmt.Errorf("response body too large or read error: %v", err)
	}
	if err != nil {
		return nil, err
	}

	var result AmazonResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	return &result, nil
}
