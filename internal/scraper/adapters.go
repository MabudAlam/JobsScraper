package scraper

import (
	"context"

	"jobscraper/common"
	db "jobscraper/db"
	amazonfetch "jobscraper/scrapers/amazon/fetch"
	amazonnormalize "jobscraper/scrapers/amazon/normalize"
	ashbyfetch "jobscraper/scrapers/ashby/fetch"
	ashbynormalize "jobscraper/scrapers/ashby/normalize"
	atlassianfetch "jobscraper/scrapers/atlassian/fetch"
	atlassiannormalize "jobscraper/scrapers/atlassian/normalize"
	leverfetch "jobscraper/scrapers/lever/fetch"
	levernormalize "jobscraper/scrapers/lever/normalize"
	"jobscraper/target"
)

type amazonScraper struct{}

func (s *amazonScraper) Name() string { return "Amazon" }

func (s *amazonScraper) Fetch(ctx context.Context) ([]*common.JobPayload, error) {
	raw, err := amazonfetch.FetchAmazonJobs()
	if err != nil || len(raw.Jobs) == 0 {
		return nil, err
	}
	jobs := amazonnormalize.NormalizeAmazonJobs(raw)
	if len(jobs) > 0 {
		db.InsertJobs(jobs)
	}
	return jobs, nil
}

type atlassianScraper struct{}

func (s *atlassianScraper) Name() string { return "Atlassian" }

func (s *atlassianScraper) Fetch(ctx context.Context) ([]*common.JobPayload, error) {
	raw, err := atlassianfetch.FetchAtlassianJobs()
	if err != nil || len(*raw) == 0 {
		return nil, err
	}
	jobs := atlassiannormalize.NormalizeAtlassianJobs(raw)
	if len(jobs) > 0 {
		db.InsertJobs(jobs)
	}
	return jobs, nil
}

type leverScraper struct {
	company target.CompanyInfo
}

func (s *leverScraper) Name() string { return s.company.Company }

func (s *leverScraper) Fetch(ctx context.Context) ([]*common.JobPayload, error) {
	raw, err := leverfetch.FetchLeverJobs(s.company.LeverSlug)
	if err != nil || len(*raw) == 0 {
		return nil, err
	}
	jobs := levernormalize.NormalizeLeverJobs(raw, s.company.Company)
	if len(jobs) > 0 {
		db.InsertJobs(jobs)
	}
	return jobs, nil
}

type ashbyScraper struct {
	company target.CompanyInfo
}

func (s *ashbyScraper) Name() string { return s.company.Company }

func (s *ashbyScraper) Fetch(ctx context.Context) ([]*common.JobPayload, error) {
	raw, err := ashbyfetch.FetchJobBoard(s.company.AshbySlug)
	if err != nil || len(raw.Jobs) == 0 {
		return nil, err
	}
	jobs := ashbynormalize.NormalizeResponse(raw, s.company.Company)
	if len(jobs) > 0 {
		db.InsertJobs(jobs)
	}
	return jobs, nil
}

func GetAllScrapers() []Scraper {
	var scrapers []Scraper

	scrapers = append(scrapers, &amazonScraper{})
	scrapers = append(scrapers, &atlassianScraper{})

	for _, c := range target.GetEnabledLeverCompanies() {
		if c.Enabled {
			scrapers = append(scrapers, &leverScraper{company: c})
		}
	}

	for _, c := range target.GetEnabledCompanies() {
		if c.Enabled {
			scrapers = append(scrapers, &ashbyScraper{company: c})
		}
	}

	return scrapers
}
