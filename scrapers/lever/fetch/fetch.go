package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

var slugPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

const maxBodyBytes = 10 * 1024 * 1024

type Category struct {
	Commitment   string   `json:"commitment"`
	Location     string   `json:"location"`
	Department   string   `json:"department"`
	Team         string   `json:"team"`
	AllLocations []string `json:"allLocations"`
}

type RawJob struct {
	ID               string   `json:"id"`
	Text             string   `json:"text"`
	DescriptionPlain string   `json:"descriptionPlain"`
	Description      string   `json:"description"`
	AdditionalPlain  string   `json:"additionalPlain"`
	Additional       string   `json:"additional"`
	Categories       Category `json:"categories"`
	CreatedAt        int64    `json:"createdAt"`
	HostedUrl        string   `json:"hostedUrl"`
	ApplyUrl         string   `json:"applyUrl"`
	WorkplaceType    string   `json:"workplaceType"`
}

type LeverResponse []RawJob

func FetchLeverJobs(companySlug string) (*LeverResponse, error) {
	if !slugPattern.MatchString(companySlug) {
		return nil, fmt.Errorf("invalid company slug: %q", companySlug)
	}

	url := fmt.Sprintf("https://api.lever.co/v0/postings/%s?mode=json", companySlug)

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

	var result LeverResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	return &result, nil
}
