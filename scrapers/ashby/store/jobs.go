package store

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"ashbyimpl/common"
)

const MaxSnapshotsPerJob = 3

type UpsertResult string

const (
	Inserted  UpsertResult = "inserted"
	Updated   UpsertResult = "updated"
	Unchanged UpsertResult = "unchanged"
)

func UpsertJob(job *common.JobPayload) (UpsertResult, error) {
	var existingHash sql.NullString
	row := db.QueryRow(`SELECT content_hash FROM jobs WHERE company = ? AND job_id = ?`,
		job.CompanyName, job.Meta.ContentHash)
	_ = row.Scan(&existingHash)

	isRemote := 0
	if job.Meta.Remote {
		isRemote = 1
	}

	publishedAt := ""
	if !job.Date.IsZero() {
		publishedAt = job.Date.Format(time.RFC3339)
	}

	scrapedAt := time.Now().Format(time.RFC3339)

	result, err := db.Exec(`
		INSERT INTO jobs (
			job_id, company, title, location, team, department,
			employment_type, remote, description,
			apply_url, job_url, published_at, scraped_at,
			compensation_summary, content_hash, is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)
		ON CONFLICT(company, job_id) DO UPDATE SET
			title = EXCLUDED.title,
			location = EXCLUDED.location,
			team = EXCLUDED.team,
			department = EXCLUDED.department,
			employment_type = EXCLUDED.employment_type,
			remote = EXCLUDED.remote,
			description = EXCLUDED.description,
			apply_url = EXCLUDED.apply_url,
			job_url = EXCLUDED.job_url,
			published_at = EXCLUDED.published_at,
			scraped_at = EXCLUDED.scraped_at,
			compensation_summary = EXCLUDED.compensation_summary,
			content_hash = CASE WHEN jobs.content_hash = EXCLUDED.content_hash THEN jobs.content_hash ELSE EXCLUDED.content_hash END,
			is_active = 1,
			updated_at = CURRENT_TIMESTAMP`,
		job.Meta.ContentHash, job.CompanyName, job.JobName, job.Meta.Location, job.Meta.Team,
		job.Meta.Department, job.Meta.EmploymentType, isRemote, job.Description,
		job.ApplyLink, job.Meta.JobURL, publishedAt, scrapedAt,
		job.Meta.Compensation, job.Meta.ContentHash)
	if err != nil {
		return "", err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		return Inserted, nil
	} else if existingHash.Valid && existingHash.String == job.Meta.ContentHash {
		return Unchanged, nil
	}
	return Updated, nil
}

func MarkRemovedJobs(company string, activeJobIds []string) ([]string, error) {
	var currentActive []string
	rows, err := db.Query(`SELECT job_id FROM jobs WHERE company = ? AND is_active = 1`, company)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var jobID string
		if err := rows.Scan(&jobID); err != nil {
			continue
		}
		currentActive = append(currentActive, jobID)
	}

	activeSet := make(map[string]bool)
	for _, id := range activeJobIds {
		activeSet[id] = true
	}

	var removed []string
	for _, j := range currentActive {
		if !activeSet[j] {
			removed = append(removed, j)
		}
	}

	if len(removed) == 0 {
		return removed, nil
	}

	placeholders := make([]string, len(removed))
	args := make([]interface{}, len(removed)+1)
	args[0] = company
	for i, r := range removed {
		placeholders[i] = "?"
		args[i+1] = r
	}

	query := `UPDATE jobs SET is_active = 0, updated_at = CURRENT_TIMESTAMP
		WHERE company = ? AND job_id IN (` + strings.Join(placeholders, ",") + `)`
	_, err = db.Exec(query, args...)
	return removed, err
}

func SaveSnapshot(job *common.JobPayload) error {
	var count int
	row := db.QueryRow(`
		SELECT COUNT(*) FROM job_snapshots
		WHERE company = ? AND job_id = ? AND content_hash = ?`,
		job.CompanyName, job.Meta.ContentHash, job.Meta.ContentHash)
	row.Scan(&count)
	if count > 0 {
		return nil
	}

	snapshotJSON, _ := json.Marshal(job)

	_, err := db.Exec(`
		INSERT INTO job_snapshots (job_id, company, content_hash, snapshot_data)
		VALUES (?, ?, ?, ?)`,
		job.Meta.ContentHash, job.CompanyName, job.Meta.ContentHash, string(snapshotJSON))
	if err != nil {
		return err
	}

	_, _ = db.Exec(`
		DELETE FROM job_snapshots WHERE company = ? AND job_id = ? AND id NOT IN (
			SELECT id FROM job_snapshots WHERE company = ? AND job_id = ?
			ORDER BY captured_at DESC LIMIT ?
		)`,
		job.CompanyName, job.Meta.ContentHash, job.CompanyName, job.Meta.ContentHash, MaxSnapshotsPerJob)

	return nil
}

func GetAllActiveJobs() ([]common.JobPayload, error) {
	rows, err := db.Query(`
		SELECT job_id, company, title, location, team, department,
			employment_type, remote, description,
			apply_url, job_url, published_at, compensation_summary
		FROM jobs WHERE is_active = 1
		ORDER BY company, published_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []common.JobPayload
	for rows.Next() {
		var j common.JobPayload
		var jobID, location, team, department, employmentType, description, applyURL, jobURL, publishedAt, compSummary *string
		var remote int

		if err := rows.Scan(&jobID, &j.CompanyName, &j.JobName, &location, &team, &department,
			&employmentType, &remote, &description,
			&applyURL, &jobURL, &publishedAt, &compSummary); err != nil {
			continue
		}

		j.ApplyLink = derefString(applyURL)
		j.Meta.JobURL = derefString(jobURL)
		j.Meta.Location = derefString(location)
		j.Meta.Team = derefString(team)
		j.Meta.Department = derefString(department)
		j.Meta.EmploymentType = derefString(employmentType)
		j.Meta.Compensation = derefString(compSummary)
		j.Meta.Remote = remote == 1
		j.Description = derefString(description)
		j.Meta.ContentHash = derefString(jobID)
		j.Meta.Source = "ashby"

		if publishedAt != nil && *publishedAt != "" {
			if t, err := time.Parse(time.RFC3339, *publishedAt); err == nil {
				j.Date = t
			}
		}

		jobs = append(jobs, j)
	}
	return jobs, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func UpsertCompany(name, ashbySlug string) error {
	slug := strings.ToLower(strings.TrimSpace(ashbySlug))
	_, err := db.Exec(`
		INSERT INTO companies (name, ashby_slug)
		VALUES (?, ?)
		ON CONFLICT(ashby_slug) DO UPDATE SET name = EXCLUDED.name`,
		name, slug)
	return err
}

func UpdateLastScraped(ashbySlug string) error {
	slug := strings.ToLower(strings.TrimSpace(ashbySlug))
	_, err := db.Exec(`
		UPDATE companies SET last_scraped_at = CURRENT_TIMESTAMP
		WHERE LOWER(ashby_slug) = LOWER(?)`,
		slug)
	return err
}

func GetAllCompaniesLastScraped() (map[string]*string, error) {
	rows, err := db.Query(`
		SELECT ashby_slug, last_scraped_at FROM companies
		ORDER BY last_scraped_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*string)
	for rows.Next() {
		var slug string
		var lastScraped *string
		if err := rows.Scan(&slug, &lastScraped); err != nil {
			continue
		}
		key := strings.ToLower(slug)
		if _, exists := result[key]; !exists {
			result[key] = lastScraped
		}
	}
	return result, nil
}
