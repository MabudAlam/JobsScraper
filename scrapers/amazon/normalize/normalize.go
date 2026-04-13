package normalize

import (
	"strings"
	"time"

	"jobscraper/common"
	"jobscraper/scrapers/amazon/fetch"
)

func NormalizeAmazonJobs(apiResponse *fetch.AmazonResponse) []*common.JobPayload {
	if apiResponse == nil || len(apiResponse.Jobs) == 0 {
		return []*common.JobPayload{}
	}

	var jobs []*common.JobPayload
	for _, j := range apiResponse.Jobs {
		jobID := j.ID
		if jobID == "" {
			continue
		}

		location := j.NormalizedLocation
		if location == "" {
			location = j.Location
		}
		if location == "" {
			location = "Unknown"
		}

		title := j.Title
		if title == "" {
			title = "Untitled"
		}

		description := j.DescriptionShort
		if description == "" {
			description = j.Description
		}

		jobPath := j.JobPath
		if !strings.HasPrefix(jobPath, "/") || strings.HasPrefix(jobPath, "//") {
			jobPath = ""
		}

		applyLink := ""
		if jobPath != "" {
			applyLink = "https://www.amazon.jobs" + jobPath
		}

		payload := &common.JobPayload{
			JobName:     title,
			Description: description,
			Date:        time.Now(),
			ApplyLink:   applyLink,
			CompanyName: "Amazon",
			Meta: common.JobMeta{
				Location:    location,
				ContentHash: jobID,
				Source:      "amazon",
			},
		}

		jobs = append(jobs, payload)
	}

	return jobs
}
