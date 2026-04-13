package normalize

import (
	"fmt"
	"time"

	"jobscraper/common"
	"jobscraper/scrapers/lever/fetch"
)

func NormalizeLeverJobs(apiResponse *fetch.LeverResponse, company string) []*common.JobPayload {
	if apiResponse == nil || len(*apiResponse) == 0 {
		return []*common.JobPayload{}
	}

	var jobs []*common.JobPayload
	for _, j := range *apiResponse {
		jobID := fmt.Sprintf("lever-%s-%s", company, j.ID)
		if j.ID == "" {
			continue
		}

		location := j.Categories.Location
		if location == "" {
			location = "Unknown"
		}

		title := j.Text
		if title == "" {
			title = "Untitled"
		}

		description := j.DescriptionPlain
		if description == "" {
			description = j.Description
		}
		if j.AdditionalPlain != "" {
			description += "\n\n" + j.AdditionalPlain
		}

		postedAt := time.Now()
		if j.CreatedAt > 0 {
			postedAt = time.Unix(j.CreatedAt/1000, 0)
		}

		payload := &common.JobPayload{
			JobName:     title,
			Description: description,
			Date:        postedAt,
			ApplyLink:   j.ApplyUrl,
			CompanyName: company,
			Meta: common.JobMeta{
				Location:       location,
				ContentHash:    jobID,
				Source:         "lever",
				Department:     j.Categories.Department,
				Team:           j.Categories.Team,
				EmploymentType: j.Categories.Commitment,
			},
		}

		jobs = append(jobs, payload)
	}

	return jobs
}
