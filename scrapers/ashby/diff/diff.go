package diff

import (
	"fmt"

	"jobscraper/common"
	embeddings "jobscraper/db"
	"jobscraper/scrapers/ashby/utils"
)

func DetectChanges(normalizedJobs []*common.JobPayload, company string) []common.Change {
	var changes []common.Change
	var currentActiveJobIds []string

	for _, job := range normalizedJobs {
		currentActiveJobIds = append(currentActiveJobIds, job.Meta.ContentHash)
		result, err := embeddings.UpsertJob(job)
		if err != nil {
			utils.LoggerInstance.Error(fmt.Sprintf("Upsert error for %s: %v", job.JobName, err))
			continue
		}

		if result == embeddings.Inserted {
			changes = append(changes, common.Change{Type: common.ChangeNew, Job: job})
			_ = embeddings.SaveSnapshot(job)
			utils.LoggerInstance.Debug("NEW:", job.JobName, "at", company)
		} else if result == embeddings.Updated {
			changes = append(changes, common.Change{Type: common.ChangeUpdated, Job: job})
			_ = embeddings.SaveSnapshot(job)
			utils.LoggerInstance.Debug("UPDATED:", job.JobName, "at", company)
		}
	}

	removed, _ := embeddings.MarkRemovedJobs(company, currentActiveJobIds)
	for _, r := range removed {
		job := &common.JobPayload{
			JobName: r,
			Meta: common.JobMeta{
				ContentHash: r,
			},
		}
		changes = append(changes, common.Change{Type: common.ChangeRemoved, Job: job})
		utils.LoggerInstance.Debug("REMOVED:", r, "at", company)
	}

	newCount := countByType(changes, string(common.ChangeNew))
	updatedCount := countByType(changes, string(common.ChangeUpdated))
	removedCount := countByType(changes, string(common.ChangeRemoved))
	utils.LoggerInstance.Info("Diff for", company+":", newCount, "new,", updatedCount, "updated,", removedCount, "removed")

	return changes
}

func countByType(changes []common.Change, typ string) int {
	count := 0
	for _, c := range changes {
		if string(c.Type) == typ {
			count++
		}
	}
	return count
}
