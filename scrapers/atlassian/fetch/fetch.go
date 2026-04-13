package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RawJob struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Locations        []string `json:"locations"`
	Category         string   `json:"category"`
	Overview         string   `json:"overview"`
	Responsibilities string   `json:"responsibilities"`
	Qualifications   string   `json:"qualifications"`
	ApplyUrl         string   `json:"applyUrl"`
}

type AtlassianResponse []RawJob

func FetchAtlassianJobs() (*AtlassianResponse, error) {
	url := "https://www.atlassian.com/endpoint/careers/listings"

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result AtlassianResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	return &result, nil
}
