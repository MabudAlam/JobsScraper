package services

import (
	"bytes"
	"encoding/json"
	"goscraper/models"
	"io"
	"net/http"
	"time"
)

func AtlanScrapper() ([]models.Job, error) {
	// URL of the API to fetch data from.
	url := "https://jobs.ashbyhq.com/api/non-user-graphql?op=ApiJobBoardWithTeams"

	body := []byte(`{
    "operationName": "ApiJobBoardWithTeams",
    "variables": {
        "organizationHostedJobsPageName": "atlan"
    },
    "query": "query ApiJobBoardWithTeams($organizationHostedJobsPageName: String!) {\n  jobBoard: jobBoardWithTeams(\n    organizationHostedJobsPageName: $organizationHostedJobsPageName\n  ) {\n    teams {\n      id\n      name\n      parentTeamId\n      __typename\n    }\n    jobPostings {\n      id\n      title\n      teamId\n      locationId\n      locationName\n      workplaceType\n      employmentType\n      secondaryLocations {\n        ...JobPostingSecondaryLocationParts\n        __typename\n      }\n      compensationTierSummary\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment JobPostingSecondaryLocationParts on JobPostingSecondaryLocation {\n  locationId\n  locationName\n  __typename\n}"
}`) // Your request body

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	// Add more headers as needed

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var atlanJobData models.AtlanJobData
	if err := json.Unmarshal(respBody, &atlanJobData); err != nil {
		return nil, err
	}

	postings := make([]models.Job, len(atlanJobData.Data.JobBoard.JobPostings))
	for i, posting := range atlanJobData.Data.JobBoard.JobPostings {
		location := string(posting.LocationName)
		postings[i] = models.Job{
			ID:          posting.ID,
			Title:       posting.Title,
			Location:    location,
			Description: "",
			CreatedAt:   time.Now().Unix(),
			Company:     "Atlan",
			ApplyURL:    "https://atlan.com/careers/?ashby_jid=" + posting.ID + "#job-list", // Update if you have a direct apply URL
			ImageUrl:    "https://res.cloudinary.com/dc0tfxkph/image/upload/v1750322245/Go%20Scraper/xfw60yza1pos3obakpsj.png",
		}
	}

	return postings, nil
}
