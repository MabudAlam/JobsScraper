package normalize

import (
	"fmt"
	"time"

	"jobscraper/common"
	"jobscraper/scrapers/atlassian/fetch"
)

func NormalizeAtlassianJobs(apiResponse *fetch.AtlassianResponse) []*common.JobPayload {
	if apiResponse == nil || len(*apiResponse) == 0 {
		return []*common.JobPayload{}
	}

	var jobs []*common.JobPayload
	for _, j := range *apiResponse {
		jobID := fmt.Sprintf("atlassian-%d", j.ID)
		if j.ID == 0 {
			continue
		}

		location := ""
		if len(j.Locations) > 0 {
			location = j.Locations[0]
		}
		if location == "" {
			location = "Unknown"
		}

		title := j.Title
		if title == "" {
			title = "Untitled"
		}

		description := j.Overview + "\n" + j.Responsibilities + "\n" + j.Qualifications

		payload := &common.JobPayload{
			JobName:     title,
			Description: description,
			Date:        time.Now(),
			ApplyLink:   j.ApplyUrl,
			CompanyName: "Atlassian",
			Meta: common.JobMeta{
				Location:    location,
				ContentHash: jobID,
				Source:      "atlassian",
				Department:  j.Category,
			},
		}

		jobs = append(jobs, payload)
	}

	return jobs
}
