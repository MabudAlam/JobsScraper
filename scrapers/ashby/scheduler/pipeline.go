package scheduler

import (
	"time"

	"ashbyimpl/common"
	"ashbyimpl/scrapers/ashby/config"
	"ashbyimpl/scrapers/ashby/diff"
	"ashbyimpl/scrapers/ashby/fetch"
	"ashbyimpl/scrapers/ashby/intelligence"
	"ashbyimpl/scrapers/ashby/normalize"
	"ashbyimpl/scrapers/ashby/notify"
	"ashbyimpl/scrapers/ashby/store"
	"ashbyimpl/scrapers/ashby/utils"
	"ashbyimpl/target"
)

func RunPipeline() {
	startTime := time.Now().UnixNano() / 1000000
	utils.LoggerInstance.Info("Pipeline run started")

	_ = store.InitDB()

	allCompanies := target.GetEnabledCompanies()
	lastScraped, _ := store.GetAllCompaniesLastScraped()
	companies := target.GetDueCompanies(lastScraped, allCompanies)

	if len(companies) == 0 {
		utils.LoggerInstance.Info("No companies due for scraping")
		return
	}

	utils.LoggerInstance.Info("Processing", len(companies), "companies")

	var allChanges []common.Change
	var allNewJobs []*common.JobPayload

	for i, company := range companies {
		_ = store.UpsertCompany(company.Company, company.AshbySlug)

		rawData, err := fetch.FetchJobBoard(company.AshbySlug)
		if err != nil {
			utils.LoggerInstance.Error("Fetch error for", company.Company+":", err.Error())
			continue
		}

		normalizedJobs := normalize.NormalizeResponse(rawData, company.Company)
		changes := diff.DetectChanges(normalizedJobs, company.Company)
		allChanges = append(allChanges, changes...)

		for _, c := range changes {
			if c.Type == common.ChangeNew {
				allNewJobs = append(allNewJobs, c.Job)
			}
		}

		_ = store.UpdateLastScraped(company.AshbySlug)
		utils.LoggerInstance.Info("Completed", company.Company+":", len(normalizedJobs), "jobs,", len(changes), "changes")

		if i < len(companies)-1 {
			cfg := config.LoadConfig()
			utils.JitteredDelay(cfg.Fetch.DelayBetweenCompaniesMin, cfg.Fetch.DelayBetweenCompaniesMax)
		}
	}

	activeRows, _ := store.GetAllActiveJobs()
	var allActiveJobs []*common.JobPayload
	for i := range activeRows {
		allActiveJobs = append(allActiveJobs, &activeRows[i])
	}

	_, filtered := intelligence.FilterAndRank(allActiveJobs)

	var scoredJobs []*intelligence.ScoredJob
	for _, s := range filtered {
		scoredJobs = append(scoredJobs, s)
	}

	cfg := config.LoadConfig()
	if cfg.Notify.CLI {
		notify.PrintRunSummary(allChanges, scoredJobs)
	}

	if cfg.Notify.Markdown && len(allChanges) > 0 {
		notify.GenerateReport(allChanges, scoredJobs)
	}

	elapsed := (time.Now().UnixNano()/1000000 - startTime) / 1000
	utils.LoggerInstance.Info("Pipeline completed in", elapsed, "s -", len(allChanges), "total changes")
}
