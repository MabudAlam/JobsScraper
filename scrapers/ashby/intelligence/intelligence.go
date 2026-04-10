package intelligence

import (
	"strings"
	"time"

	"ashbyimpl/common"
	"ashbyimpl/scrapers/ashby/config"
	"ashbyimpl/scrapers/ashby/utils"
)

func ScoreJob(job *common.JobPayload) (int, []string) {
	cfg := config.LoadConfig()
	rules := cfg.Intelligence
	score := 0
	signals := []string{}

	score += scoreKeywords(job, rules, &signals)
	score += scoreLocation(job, rules, &signals)
	score += scoreRemote(job, rules, &signals)
	score += scoreDepartment(job, rules, &signals)
	score += scoreFreshness(job, rules, &signals)

	return score, signals
}

func scoreKeywords(job *common.JobPayload, rules config.IntelligenceConfig, signals *[]string) int {
	if rules.Keywords == nil {
		return 0
	}
	total := 0
	searchText := strings.ToLower(job.JobName + " " + job.Description)

	for keyword, weight := range rules.Keywords {
		if strings.Contains(searchText, strings.ToLower(keyword)) {
			total += weight
			*signals = append(*signals, "keyword:\""+keyword+"\" ("+formatSign(weight)+")")
		}
	}
	return total
}

func scoreLocation(job *common.JobPayload, rules config.IntelligenceConfig, signals *[]string) int {
	if rules.Locations == nil {
		return 0
	}
	jobLocation := strings.ToLower(job.Meta.Location)
	for _, loc := range rules.Locations {
		if strings.Contains(jobLocation, strings.ToLower(loc)) {
			*signals = append(*signals, "location:\""+loc+"\" (+"+itoa(rules.LocationBoost)+")")
			return rules.LocationBoost
		}
	}
	return 0
}

func scoreRemote(job *common.JobPayload, rules config.IntelligenceConfig, signals *[]string) int {
	if job.Meta.Remote && rules.RemoteBoost > 0 {
		*signals = append(*signals, "remote (+"+itoa(rules.RemoteBoost)+")")
		return rules.RemoteBoost
	}
	return 0
}

func scoreDepartment(job *common.JobPayload, rules config.IntelligenceConfig, signals *[]string) int {
	if rules.Departments == nil || job.Meta.Department == "" {
		return 0
	}
	dept := strings.ToLower(job.Meta.Department)
	for _, d := range rules.Departments {
		if strings.Contains(dept, strings.ToLower(d)) {
			*signals = append(*signals, "department:\""+d+"\" (+"+itoa(rules.DepartmentBoost)+")")
			return rules.DepartmentBoost
		}
	}
	return 0
}

func scoreFreshness(job *common.JobPayload, rules config.IntelligenceConfig, signals *[]string) int {
	if rules.FreshnessBoostHrs == 0 || rules.FreshnessBoost == 0 {
		return 0
	}
	if job.Date.IsZero() {
		return 0
	}
	hoursAgo := time.Since(job.Date).Hours()
	if hoursAgo <= float64(rules.FreshnessBoostHrs) {
		*signals = append(*signals, "fresh:"+itoa(int(hoursAgo))+"h (+"+itoa(rules.FreshnessBoost)+")")
		return rules.FreshnessBoost
	}
	return 0
}

func FilterAndRank(jobs []*common.JobPayload) (all []*ScoredJob, filtered []*ScoredJob) {
	cfg := config.LoadConfig()
	minScore := cfg.Intelligence.MinScore

	var scored []*ScoredJob
	for _, job := range jobs {
		score, signals := ScoreJob(job)
		scored = append(scored, &ScoredJob{
			JobPayload:     job,
			RelevanceScore: score,
			Signals:        signals,
		})
	}

	for i := range scored {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].RelevanceScore > scored[i].RelevanceScore {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	var above []*ScoredJob
	for _, s := range scored {
		if s.RelevanceScore >= minScore {
			above = append(above, s)
		}
	}

	utils.LoggerInstance.Info("Intelligence:", len(above), "/", len(scored), "jobs above threshold (min:", minScore, ")")
	return scored, above
}

type ScoredJob struct {
	*common.JobPayload
	RelevanceScore int
	Signals        []string
}

func formatSign(n int) string {
	if n > 0 {
		return "+" + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}
