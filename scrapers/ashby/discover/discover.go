package discover

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

var cachedClient *http.Client

const maxBodyBytes = 10 * 1024 * 1024

func getClient() *http.Client {
	if cachedClient == nil {
		cachedClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return cachedClient
}

type Result struct {
	Title string
	URL   string
}

type SearXNGResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type SearXNGResponse struct {
	Results []SearXNGResult `json:"results"`
}

func Search(query string) ([]Result, error) {
	searxngURL := os.Getenv("SEARXNG_URL")

	if searxngURL != "" {
		return searxngSearch(query, searxngURL)
	}

	apiKey := os.Getenv("GOOGLE_SEARCH_API_KEY")
	cx := os.Getenv("GOOGLE_SEARCH_CX")

	if apiKey != "" && cx != "" {
		return googleSearch(query, apiKey, cx)
	}

	return nil, fmt.Errorf("no search provider configured: set SEARXNG_URL or GOOGLE_SEARCH_API_KEY + GOOGLE_SEARCH_CX")
}

func searxngSearch(query string, baseURL string) ([]Result, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s&format=json&engines=google", baseURL, url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := getClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	limitedBody := http.MaxBytesReader(nil, resp.Body, maxBodyBytes)
	body, err := io.ReadAll(limitedBody)
	if err != nil {
		return nil, fmt.Errorf("response body too large or read error: %v", err)
	}

	var result SearXNGResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("SearXNG parse error: %v", err)
	}

	var results []Result
	for _, item := range result.Results {
		results = append(results, Result{
			Title: item.Title,
			URL:   item.URL,
		})
	}

	return results, nil
}

func googleSearch(query string, apiKey string, cx string) ([]Result, error) {
	baseURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?cx=%s&q=%s", cx, url.QueryEscape(query))

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Goog-Api-Key", apiKey)

	resp, err := getClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	limitedBody := http.MaxBytesReader(nil, resp.Body, maxBodyBytes)
	body, err := io.ReadAll(limitedBody)
	if err != nil {
		return nil, fmt.Errorf("response body too large or read error: %v", err)
	}
	if err != nil {
		return nil, err
	}

	var result GoogleSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var results []Result
	for _, item := range result.Items {
		results = append(results, Result{
			Title: item.Title,
			URL:   item.Link,
		})
	}

	return results, nil
}

type GoogleSearchResponse struct {
	Items []GoogleSearchItem `json:"items"`
}

type GoogleSearchItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func DiscoverCompanies() ([]string, error) {
	queries := []string{
		"site:jobs.ashbyhq.com",
		"jobs.ashbyhq.com careers",
	}

	var allResults []Result
	for _, query := range queries {
		results, err := Search(query)
		if err != nil {
			fmt.Printf("Query '%s' failed: %v\n", query, err)
			continue
		}
		allResults = append(allResults, results...)
		time.Sleep(2 * time.Second)
	}

	seen := make(map[string]bool)
	var companies []string

	for _, r := range allResults {
		slug := extractAshbySlug(r.URL)
		if slug != "" && !seen[slug] {
			seen[slug] = true
			companies = append(companies, slug)
		}
	}

	return companies, nil
}

func extractAshbySlug(rawURL string) string {
	re := regexp.MustCompile(`jobs\.ashbyhq\.com/([a-zA-Z0-9_-]+)`)
	matches := re.FindStringSubmatch(rawURL)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
