package notify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ashbyimpl/common"
	"ashbyimpl/scrapers/ashby/config"
	"ashbyimpl/scrapers/ashby/intelligence"
)

func PrintRunSummary(changes []common.Change, scoredJobs []*intelligence.ScoredJob) {
	fmt.Println("\n==================================================")
	fmt.Println("  AshbyHQ Job Scraper - Run Summary")
	fmt.Println("==================================================\n")

	if len(changes) == 0 {
		fmt.Println("  No changes detected.\n")
		return
	}

	byCompany := groupByChanges(changes)
	for company, cs := range byCompany {
		newCount := countByType(cs, string(common.ChangeNew))
		updatedCount := countByType(cs, string(common.ChangeUpdated))
		removedCount := countByType(cs, string(common.ChangeRemoved))

		fmt.Printf("  %s\n", company)
		fmt.Printf("  %d new · %d updated · %d removed\n\n", newCount, updatedCount, removedCount)

		for _, c := range cs {
			if c.Type == common.ChangeNew {
				job := c.Job
				fmt.Printf("    + %s", job.JobName)
				if scoredJobs != nil {
					for _, s := range scoredJobs {
						if s.Meta.ContentHash == job.Meta.ContentHash {
							fmt.Printf(" [score: %d]", s.RelevanceScore)
							break
						}
					}
				}
				fmt.Println()
				fmt.Printf("      %s\n", formatMeta(job))
				if job.Meta.Compensation != "" {
					fmt.Printf("      %s\n", job.Meta.Compensation)
				}
				fmt.Printf("      %s\n\n", job.ApplyLink)
			} else if c.Type == common.ChangeUpdated {
				job := c.Job
				fmt.Printf("    ~ %s\n", job.JobName)
				fmt.Printf("      %s\n\n", formatMeta(job))
			} else if c.Type == common.ChangeRemoved {
				job := c.Job
				fmt.Printf("    - %s\n", job.JobName)
				fmt.Printf("      %s\n\n", formatMeta(job))
			}
		}
	}

	if scoredJobs != nil && len(scoredJobs) > 0 {
		fmt.Println("----------------------------------------------------")
		fmt.Printf("  Top Opportunities (%d above threshold)\n", len(scoredJobs))
		fmt.Println("----------------------------------------------------\n")

		top := scoredJobs
		if len(top) > 10 {
			top = top[:10]
		}
		for _, job := range top {
			fmt.Printf("  %s at %s\n", job.JobName, job.CompanyName)
			fmt.Printf("    Score: %d - %s\n", job.RelevanceScore, strings.Join(job.Signals, ", "))
			fmt.Printf("    %s\n", formatMeta(job.JobPayload))
			if job.Meta.Compensation != "" {
				fmt.Printf("    %s\n", job.Meta.Compensation)
			}
			fmt.Printf("    %s\n\n", job.ApplyLink)
		}
	}

	fmt.Println("==================================================\n")
}

func formatMeta(job *common.JobPayload) string {
	parts := []string{job.Meta.Location}
	if job.Meta.Remote {
		parts = append(parts, "Remote")
	}
	if job.Meta.Department != "" {
		parts = append(parts, job.Meta.Department)
	}
	if job.Meta.Team != "" {
		parts = append(parts, job.Meta.Team)
	}
	if job.Meta.EmploymentType != "" && !contains(parts, job.Meta.EmploymentType) {
		parts = append(parts, job.Meta.EmploymentType)
	}
	return strings.Join(parts, " · ")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func groupByChanges(changes []common.Change) map[string][]common.Change {
	result := make(map[string][]common.Change)
	for _, c := range changes {
		if result[c.Job.CompanyName] == nil {
			result[c.Job.CompanyName] = []common.Change{}
		}
		result[c.Job.CompanyName] = append(result[c.Job.CompanyName], c)
	}
	return result
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

func GenerateReport(changes []common.Change, scoredJobs []*intelligence.ScoredJob) {
	cfg := config.LoadConfig()
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	timeStr := now.Format("15:04")

	dir := cfg.Notify.ReportsDir
	os.MkdirAll(dir, 0755)

	lines := []string{
		fmt.Sprintf("# Job Scraper Report - %s", dateStr),
		"",
		fmt.Sprintf("> Generated at %s %s UTC", dateStr, timeStr),
		"",
	}

	newJobs := filterByType(changes, string(common.ChangeNew))
	updatedJobs := filterByType(changes, string(common.ChangeUpdated))
	removedJobs := filterByType(changes, string(common.ChangeRemoved))

	lines = append(lines, "## Summary", "")
	lines = append(lines, "| Metric | Count |", "|--------|-------|")
	lines = append(lines, fmt.Sprintf("| New Jobs | %d |", len(newJobs)))
	lines = append(lines, fmt.Sprintf("| Updated Jobs | %d |", len(updatedJobs)))
	lines = append(lines, fmt.Sprintf("| Removed Jobs | %d |", len(removedJobs)))
	lines = append(lines, fmt.Sprintf("| Above Relevance Threshold | %d |", len(scoredJobs)))
	lines = append(lines, "")

	if scoredJobs != nil && len(scoredJobs) > 0 {
		lines = append(lines, "## Top Opportunities", "")
		for _, job := range scoredJobs {
			if len(lines) > 50 {
				break
			}
			lines = append(lines, fmt.Sprintf("### %s - %s", job.JobName, job.CompanyName))
			lines = append(lines, "")
			lines = append(lines, fmt.Sprintf("- **Score:** %d", job.RelevanceScore))
			lines = append(lines, fmt.Sprintf("- **Location:** %s", job.Meta.Location))
			if job.Meta.Remote {
				lines = append(lines, "- **Location:** "+job.Meta.Location+" (Remote)")
			}
			if job.Meta.Department != "" {
				lines = append(lines, fmt.Sprintf("- **Department:** %s", job.Meta.Department))
			}
			if job.Meta.Team != "" {
				lines = append(lines, fmt.Sprintf("- **Team:** %s", job.Meta.Team))
			}
			if job.Meta.EmploymentType != "" {
				lines = append(lines, fmt.Sprintf("- **Type:** %s", job.Meta.EmploymentType))
			}
			if job.Meta.Compensation != "" {
				lines = append(lines, fmt.Sprintf("- **Compensation:** %s", job.Meta.Compensation))
			}
			lines = append(lines, fmt.Sprintf("- **Signals:** %s", strings.Join(job.Signals, ", ")))
			lines = append(lines, fmt.Sprintf("- **Apply:** [Link](%s)", job.ApplyLink))
			lines = append(lines, "")
		}
	}

	if len(newJobs) > 0 {
		lines = append(lines, "## New Postings", "")
		for _, c := range newJobs {
			job := c.Job
			remoteStr := ""
			if job.Meta.Remote {
				remoteStr = " (Remote)"
			}
			lines = append(lines, fmt.Sprintf("- **%s** at %s — %s%s — [Apply](%s)", job.JobName, job.CompanyName, job.Meta.Location, remoteStr, job.ApplyLink))
		}
		lines = append(lines, "")
	}

	if len(updatedJobs) > 0 {
		lines = append(lines, "## Updated Postings", "")
		for _, c := range updatedJobs {
			job := c.Job
			lines = append(lines, fmt.Sprintf("- **%s** at %s — %s", job.JobName, job.CompanyName, job.Meta.Location))
		}
		lines = append(lines, "")
	}

	if len(removedJobs) > 0 {
		lines = append(lines, "## Removed Postings", "")
		for _, c := range removedJobs {
			job := c.Job
			lines = append(lines, fmt.Sprintf("- ~~%s~~ at %s — %s", job.JobName, job.CompanyName, job.Meta.Location))
		}
		lines = append(lines, "")
	}

	lines = append(lines, "---")
	lines = append(lines, "*Generated by AshbyHQ Job Scraper*")

	content := strings.Join(lines, "\n")
	filePath := filepath.Join(dir, dateStr+".md")
	existing := ""
	if _, err := os.Stat(filePath); err == nil {
		existing = " (existing replaced)"
	}
	os.WriteFile(filePath, []byte(content), 0644)
	_ = existing
}

func filterByType(changes []common.Change, typ string) []common.Change {
	var result []common.Change
	for _, c := range changes {
		if string(c.Type) == typ {
			result = append(result, c)
		}
	}
	return result
}
