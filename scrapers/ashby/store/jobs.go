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
	jobs, _, err := GetActiveJobsPaginated(0, 0, "", "", "", "")
	return jobs, err
}

func GetActiveJobsPaginated(offset, limit int, search, company, location, sort string) ([]common.JobPayload, int, error) {
	var total int
	var countQuery string
	var args []interface{}

	countQuery = `SELECT COUNT(*) FROM jobs WHERE is_active = 1`
	queryArgs := []interface{}{}

	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery = `SELECT COUNT(*) FROM jobs WHERE is_active = 1 AND (
			title LIKE ? OR description LIKE ? OR location LIKE ?
		)`
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	if company != "" {
		countQuery += ` AND company = ?`
		args = append(args, company)
	}

	if location != "" {
		countQuery += ` AND location = ?`
		args = append(args, location)
	}

	row := db.QueryRow(countQuery, args...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, job_id, company, title, location, team, department,
			employment_type, remote, description,
			apply_url, job_url, published_at, compensation_summary
		FROM jobs WHERE is_active = 1`
	queryArgs = []interface{}{}

	if search != "" {
		searchPattern := "%" + search + "%"
		query += ` AND (
			title LIKE ? OR description LIKE ? OR location LIKE ?
		)`
		queryArgs = append(queryArgs, searchPattern, searchPattern, searchPattern)
	}

	if company != "" {
		query += ` AND company = ?`
		queryArgs = append(queryArgs, company)
	}

	if location != "" {
		query += ` AND location = ?`
		queryArgs = append(queryArgs, location)
	}

	if sort == "oldest" {
		query += ` ORDER BY published_at ASC`
	} else {
		query += ` ORDER BY published_at DESC`
	}

	if limit > 0 {
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, limit, offset)
	}

	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []common.JobPayload
	for rows.Next() {
		var j common.JobPayload
		var jobID, loc, team, department, employmentType, description, applyURL, jobURL, publishedAt, compSummary *string
		var id int
		var remote int

		if err := rows.Scan(&id, &jobID, &j.CompanyName, &j.JobName, &loc, &team, &department,
			&employmentType, &remote, &description,
			&applyURL, &jobURL, &publishedAt, &compSummary); err != nil {
			continue
		}

		j.Id = id

		j.ApplyLink = derefString(applyURL)
		j.Meta.JobURL = derefString(jobURL)
		j.Meta.Location = derefString(loc)
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
	return jobs, total, nil
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

func GetAllLocations() ([]string, error) {
	rows, err := db.Query(`
		SELECT DISTINCT location FROM jobs 
		WHERE is_active = 1 AND location IS NOT NULL AND location != ''
		ORDER BY location`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []string
	for rows.Next() {
		var loc string
		if err := rows.Scan(&loc); err != nil {
			continue
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

func GetAllCompanies() ([]string, error) {
	rows, err := db.Query(`
		SELECT DISTINCT company FROM jobs 
		WHERE is_active = 1 AND company IS NOT NULL AND company != ''
		ORDER BY company`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []string
	for rows.Next() {
		var company string
		if err := rows.Scan(&company); err != nil {
			continue
		}
		companies = append(companies, company)
	}
	return companies, nil
}

func GetJobById(id int) (*common.JobPayload, error) {
	query := `
		SELECT id, job_id, company, title, location, team, department,
			employment_type, remote, description,
			apply_url, job_url, published_at, compensation_summary
		FROM jobs WHERE id = ? AND is_active = 1`

	var j common.JobPayload
	var jobID, location, team, department, employmentType, description, applyURL, jobURL, publishedAt, compSummary *string
	var remote int

	err := db.QueryRow(query, id).Scan(
		&j.Id, &jobID, &j.CompanyName, &j.JobName, &location, &team, &department,
		&employmentType, &remote, &description,
		&applyURL, &jobURL, &publishedAt, &compSummary)
	if err != nil {
		return nil, err
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

	return &j, nil
}
