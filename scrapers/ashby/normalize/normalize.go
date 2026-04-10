package normalize

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"ashbyimpl/common"
	"ashbyimpl/scrapers/ashby/fetch"
	"ashbyimpl/scrapers/ashby/utils"
)

var tagRegex = regexp.MustCompile(`<[^>]*>`)

func ExtractJobID(jobURL string) string {
	if jobURL == "" {
		return ""
	}
	u, err := url.Parse(jobURL)
	if err != nil {
		return jobURL
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func SanitizeDescription(html string) string {
	if html == "" {
		return ""
	}
	sanitized := tagRegex.ReplaceAllString(html, "")
	return strings.TrimSpace(sanitized)
}

func NormalizeJob(raw *fetch.RawJob, company string) *common.JobPayload {
	jobID := raw.ID
	if jobID == "" {
		utils.LoggerInstance.Warn("Skipping job with no extractable ID:", raw.Title)
		return nil
	}

	description := raw.DescriptionPlain
	if description == "" {
		description = SanitizeDescription(raw.DescriptionHTML)
	}

	compensation := ""
	if raw.Compensation != nil {
		if raw.Compensation.CompensationTierSummary != "" {
			compensation = raw.Compensation.CompensationTierSummary
		} else if raw.Compensation.ScrapeableCompensationSalarySummary != "" {
			compensation = raw.Compensation.ScrapeableCompensationSalarySummary
		}
	}

	publishedAt := time.Now()
	if raw.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, raw.PublishedAt); err == nil {
			publishedAt = t
		}
	}

	location := raw.Location
	if location == "" {
		location = "Unknown"
	}

	title := raw.Title
	if title == "" {
		title = "Untitled"
	}

	payload := &common.JobPayload{
		JobName:     title,
		Description: description,
		Date:        publishedAt,
		ApplyLink:   raw.ApplyURL,
		CompanyName: company,
		Meta: common.JobMeta{
			Location:       location,
			Remote:         raw.IsRemote,
			Department:     raw.Department,
			Team:           raw.Team,
			EmploymentType: raw.EmploymentType,
			Compensation:   compensation,
			JobURL:         raw.JobURL,
			ContentHash:    jobID,
			Source:         "ashby",
		},
	}

	payload.Meta.ContentHash = utils.ContentHash(
		title,
		location,
		description,
		raw.EmploymentType,
		boolStr(raw.IsRemote),
		raw.Team,
		raw.Department,
	)

	return payload
}

func NormalizeResponse(apiResponse *fetch.NormalizedResponse, company string) []*common.JobPayload {
	if apiResponse == nil || len(apiResponse.Jobs) == 0 {
		utils.LoggerInstance.Warn("No jobs array in response for", company)
		return []*common.JobPayload{}
	}

	var jobs []*common.JobPayload
	for _, j := range apiResponse.Jobs {
		if !j.IsListed {
			continue
		}
		normalized := NormalizeJob(&j, company)
		if normalized != nil {
			jobs = append(jobs, normalized)
		}
	}

	utils.LoggerInstance.Debug("Normalized", len(jobs), "jobs for", company)
	return jobs
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
