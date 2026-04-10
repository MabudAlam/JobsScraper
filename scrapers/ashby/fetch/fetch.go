package fetch

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RawJob struct {
	ID                        string              `json:"id"`
	Title                     string              `json:"title"`
	Department                string              `json:"department"`
	Team                      string              `json:"team"`
	EmploymentType            string              `json:"employmentType"`
	Location                  string              `json:"location"`
	ShouldDisplayCompensation bool                `json:"shouldDisplayCompensationOnJobPostings"`
	SecondaryLocations        []SecondaryLocation `json:"secondaryLocations"`
	PublishedAt               string              `json:"publishedAt"`
	IsListed                  bool                `json:"isListed"`
	IsRemote                  bool                `json:"isRemote"`
	WorkplaceType             string              `json:"workplaceType"`
	Address                   *Address            `json:"address"`
	JobURL                    string              `json:"jobUrl"`
	ApplyURL                  string              `json:"applyUrl"`
	DescriptionHTML           string              `json:"descriptionHtml"`
	DescriptionPlain          string              `json:"descriptionPlain"`
	Compensation              *Compensation       `json:"compensation"`
}

type SecondaryLocation struct {
	Name string `json:"name"`
}

type Address struct {
	PostalAddress PostalAddress `json:"postalAddress"`
}

type PostalAddress struct {
	AddressCountry string `json:"addressCountry"`
}

type Compensation struct {
	CompensationTierSummary             string `json:"compensationTierSummary"`
	ScrapeableCompensationSalarySummary string `json:"scrapeableCompensationSalarySummary"`
}

type NormalizedResponse struct {
	Jobs []RawJob `json:"jobs"`
}

func FetchJobBoard(ashbySlug string) (*NormalizedResponse, error) {
	baseURL := fmt.Sprintf("https://api.ashbyhq.com/posting-api/job-board/%s", ashbySlug)

	req, err := http.NewRequest("GET", baseURL, nil)
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

	var reader io.ReadCloser
	switch strings.ToLower(resp.Header.Get("Content-Encoding")) {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var result NormalizedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	return &result, nil
}
