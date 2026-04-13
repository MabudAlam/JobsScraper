package embeddings

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"jobscraper/common"

	_ "github.com/mattn/go-sqlite3"
)

type UpsertResult string

const (
	Inserted  UpsertResult = "inserted"
	Updated   UpsertResult = "updated"
	Unchanged UpsertResult = "unchanged"
)

var (
	db            *sql.DB
	dbPath        string
	dbInitOnce    sync.Once
	dbInitialized bool
)

func InitDB() error {
	var initErr error
	dbInitOnce.Do(func() {
		dbPath = os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = filepath.Join(os.Getenv("PWD"), "jobs.db")
		}

		var err error
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		if err = db.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		if err = createTables(); err != nil {
			initErr = fmt.Errorf("failed to create tables: %w", err)
			return
		}

		dbInitialized = true
		log.Println("sqlite database initialized at:", dbPath)
	})
	return initErr
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id TEXT UNIQUE NOT NULL,
		content_hash TEXT NOT NULL,
		company TEXT NOT NULL,
		title TEXT NOT NULL,
		location TEXT,
		team TEXT,
		department TEXT,
		employment_type TEXT,
		remote INTEGER,
		description TEXT,
		apply_url TEXT,
		job_url TEXT,
		compensation TEXT,
		is_active INTEGER DEFAULT 1,
		status TEXT DEFAULT 'new',
		published_at TEXT,
		scraped_at TEXT NOT NULL,
		embedding TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_jobs_content_hash ON jobs(content_hash);
	CREATE INDEX IF NOT EXISTS idx_jobs_company ON jobs(company);
	CREATE INDEX IF NOT EXISTS idx_jobs_is_active ON jobs(is_active);

	CREATE TABLE IF NOT EXISTS companies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ashby_slug TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		last_scraped TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_companies_slug ON companies(ashby_slug);

	CREATE TABLE IF NOT EXISTS deleted_jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content_hash TEXT NOT NULL,
		company TEXT NOT NULL,
		deleted_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_deleted_hash ON deleted_jobs(content_hash);
	`
	_, err := db.Exec(schema)
	return err
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func CreateCollectionIfNotExists() error {
	return InitDB()
}

func InsertJobs(jobs []*common.JobPayload) error {
	if len(jobs) == 0 || db == nil {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO jobs (job_id, content_hash, company, title, location, team, department, employment_type, remote, description, apply_url, job_url, compensation, is_active, status, published_at, scraped_at, embedding)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(job_id) DO UPDATE SET
			title = excluded.title,
			location = excluded.location,
			team = excluded.team,
			department = excluded.department,
			employment_type = excluded.employment_type,
			remote = excluded.remote,
			description = excluded.description,
			apply_url = excluded.apply_url,
			job_url = excluded.job_url,
			compensation = excluded.compensation,
			is_active = excluded.is_active,
			status = excluded.status,
			published_at = excluded.published_at,
			scraped_at = excluded.scraped_at,
			embedding = excluded.embedding
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, job := range jobs {
		embeddingJSON := ""
		if job.Embedding != nil {
			data, err := json.Marshal(job.Embedding)
			if err != nil {
				return err
			}
			embeddingJSON = string(data)
		}

		remote := 0
		if job.Meta.Remote {
			remote = 1
		}

		contentHash := job.Meta.ContentHash
		if contentHash == "" {
			contentHash = fmt.Sprintf("job_%d", i)
		}

		jobID := contentHash

		publishedAt := ""
		if !job.Date.IsZero() {
			publishedAt = job.Date.Format(time.RFC3339)
		}

		_, err = stmt.Exec(
			jobID,
			contentHash,
			job.CompanyName,
			job.JobName,
			job.Meta.Location,
			job.Meta.Team,
			job.Meta.Department,
			job.Meta.EmploymentType,
			remote,
			job.Description,
			job.ApplyLink,
			job.Meta.JobURL,
			job.Meta.Compensation,
			1,
			"new",
			publishedAt,
			time.Now().Format(time.RFC3339),
			embeddingJSON,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func UpsertJob(job *common.JobPayload) (UpsertResult, error) {
	if db == nil {
		return Inserted, fmt.Errorf("database not initialized")
	}

	embeddingJSON := ""
	if job.Embedding != nil {
		data, err := json.Marshal(job.Embedding)
		if err != nil {
			return Inserted, err
		}
		embeddingJSON = string(data)
	}

	remote := 0
	if job.Meta.Remote {
		remote = 1
	}

	contentHash := job.Meta.ContentHash
	if contentHash == "" {
		contentHash = fmt.Sprintf("job_%d", job.Id)
	}

	jobID := contentHash

	publishedAt := ""
	if !job.Date.IsZero() {
		publishedAt = job.Date.Format(time.RFC3339)
	}

	result, err := db.Exec(`
		INSERT INTO jobs (job_id, content_hash, company, title, location, team, department, employment_type, remote, description, apply_url, job_url, compensation, is_active, status, published_at, scraped_at, embedding)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(job_id) DO UPDATE SET
			title = excluded.title,
			location = excluded.location,
			team = excluded.team,
			department = excluded.department,
			employment_type = excluded.employment_type,
			remote = excluded.remote,
			description = excluded.description,
			apply_url = excluded.apply_url,
			job_url = excluded.job_url,
			compensation = excluded.compensation,
			is_active = excluded.is_active,
			status = excluded.status,
			published_at = excluded.published_at,
			scraped_at = excluded.scraped_at,
			embedding = excluded.embedding
	`,
		jobID,
		contentHash,
		job.CompanyName,
		job.JobName,
		job.Meta.Location,
		job.Meta.Team,
		job.Meta.Department,
		job.Meta.EmploymentType,
		remote,
		job.Description,
		job.ApplyLink,
		job.Meta.JobURL,
		job.Meta.Compensation,
		1,
		"new",
		publishedAt,
		time.Now().Format(time.RFC3339),
		embeddingJSON,
	)
	if err != nil {
		return Inserted, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		return Inserted, nil
	}
	return Updated, nil
}

func MarkRemovedJobs(company string, activeJobIds []string) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT job_id FROM jobs WHERE company = ? AND is_active = 1", company)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activeSet := make(map[string]bool)
	for _, id := range activeJobIds {
		activeSet[id] = true
	}

	var removed []string
	for rows.Next() {
		var jobID string
		if err := rows.Scan(&jobID); err != nil {
			continue
		}
		if !activeSet[jobID] {
			_, err := db.Exec("UPDATE jobs SET is_active = 0, status = 'removed' WHERE job_id = ?", jobID)
			if err == nil {
				removed = append(removed, jobID)
			}
		}
	}
	return removed, nil
}

func GetAllActiveJobs() ([]common.JobPayload, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT job_id, content_hash, company, title, location, team, department, employment_type, remote, description, apply_url, job_url, compensation, published_at, embedding FROM jobs WHERE is_active = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []common.JobPayload
	for rows.Next() {
		job := common.JobPayload{Meta: common.JobMeta{Source: "ashby"}}
		var remote int
		var publishedAt, embeddingJSON sql.NullString
		var location, team, department, employmentType, compensation, jobURL, description, applyURL sql.NullString

		err := rows.Scan(
			&job.Meta.ContentHash,
			&job.Meta.ContentHash,
			&job.CompanyName,
			&job.JobName,
			&location,
			&team,
			&department,
			&employmentType,
			&remote,
			&description,
			&applyURL,
			&jobURL,
			&compensation,
			&publishedAt,
			&embeddingJSON,
		)
		if err != nil {
			continue
		}

		job.Meta.Remote = remote == 1
		if location.Valid {
			job.Meta.Location = location.String
		}
		if team.Valid {
			job.Meta.Team = team.String
		}
		if department.Valid {
			job.Meta.Department = department.String
		}
		if employmentType.Valid {
			job.Meta.EmploymentType = employmentType.String
		}
		if compensation.Valid {
			job.Meta.Compensation = compensation.String
		}
		if jobURL.Valid {
			job.Meta.JobURL = jobURL.String
		}
		if description.Valid {
			job.Description = description.String
		}
		if applyURL.Valid {
			job.ApplyLink = applyURL.String
		}
		if publishedAt.Valid && publishedAt.String != "" {
			if t, err := time.Parse(time.RFC3339, publishedAt.String); err == nil {
				job.Date = t
			}
		}
		if embeddingJSON.Valid && embeddingJSON.String != "" {
			var emb []float32
			if err := json.Unmarshal([]byte(embeddingJSON.String), &emb); err == nil {
				job.Embedding = emb
			}
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func GetActiveJobsPaginated(offset, limit int, search, company, location, sort string) ([]common.JobPayload, int, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("database not initialized")
	}

	whereClause := "is_active = 1"
	args := []interface{}{}

	if company != "" {
		whereClause += " AND company = ?"
		args = append(args, company)
	}
	if location != "" {
		whereClause += " AND location = ?"
		args = append(args, location)
	}
	if search != "" {
		whereClause += " AND (title LIKE ? OR description LIKE ? OR location LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM jobs WHERE " + whereClause
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := "scraped_at DESC"
	if sort == "date" {
		orderBy = "published_at DESC"
	} else if sort == "title" {
		orderBy = "title ASC"
	}

	query := fmt.Sprintf("SELECT id, job_id, content_hash, company, title, location, team, department, employment_type, remote, description, apply_url, job_url, compensation, published_at, embedding FROM jobs WHERE %s ORDER BY %s LIMIT ? OFFSET ?", whereClause, orderBy)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []common.JobPayload
	for rows.Next() {
		job := common.JobPayload{Meta: common.JobMeta{Source: "ashby"}}
		var remote int
		var publishedAt, embeddingJSON sql.NullString
		var loc, team, dept, empType, comp, jobURL, desc, appURL sql.NullString

		err := rows.Scan(
			&job.Id,
			&job.Meta.ContentHash,
			&job.Meta.ContentHash,
			&job.CompanyName,
			&job.JobName,
			&loc,
			&team,
			&dept,
			&empType,
			&remote,
			&desc,
			&appURL,
			&jobURL,
			&comp,
			&publishedAt,
			&embeddingJSON,
		)
		if err != nil {
			continue
		}

		job.Meta.Remote = remote == 1
		if loc.Valid {
			job.Meta.Location = loc.String
		}
		if team.Valid {
			job.Meta.Team = team.String
		}
		if dept.Valid {
			job.Meta.Department = dept.String
		}
		if empType.Valid {
			job.Meta.EmploymentType = empType.String
		}
		if comp.Valid {
			job.Meta.Compensation = comp.String
		}
		if jobURL.Valid {
			job.Meta.JobURL = jobURL.String
		}
		if desc.Valid {
			job.Description = desc.String
		}
		if appURL.Valid {
			job.ApplyLink = appURL.String
		}
		if publishedAt.Valid && publishedAt.String != "" {
			if t, err := time.Parse(time.RFC3339, publishedAt.String); err == nil {
				job.Date = t
			}
		}
		if embeddingJSON.Valid && embeddingJSON.String != "" {
			var emb []float32
			if err := json.Unmarshal([]byte(embeddingJSON.String), &emb); err == nil {
				job.Embedding = emb
			}
		}

		jobs = append(jobs, job)
	}

	return jobs, total, nil
}

func SearchJobs(queryVector []float32, limit int) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("vector search not supported in sqlite backend")
}

func SearchJobsWithFilter(queryVector []float32, limit int, filter string) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("vector search not supported in sqlite backend")
}

func GetAllLocations() ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT DISTINCT location FROM jobs WHERE is_active = 1 AND location IS NOT NULL AND location != ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []string
	for rows.Next() {
		var loc string
		if err := rows.Scan(&loc); err == nil && loc != "" {
			locations = append(locations, loc)
		}
	}
	return locations, nil
}

func GetAllCompanies() ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT DISTINCT company FROM jobs WHERE is_active = 1 AND company IS NOT NULL AND company != '' ORDER BY company ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []string
	for rows.Next() {
		var comp string
		if err := rows.Scan(&comp); err == nil && comp != "" {
			companies = append(companies, comp)
		}
	}
	return companies, nil
}

func GetClient() interface{} {
	return db
}

func SaveSnapshot(job *common.JobPayload) error {
	return nil
}

func UpsertCompany(name, ashbySlug string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := db.Exec(`
		INSERT INTO companies (ashby_slug, name)
		VALUES (?, ?)
		ON CONFLICT(ashby_slug) DO UPDATE SET name = excluded.name
	`, ashbySlug, name)
	return err
}

func UpdateLastScraped(ashbySlug string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := db.Exec("UPDATE companies SET last_scraped = ? WHERE ashby_slug = ?", time.Now().Format(time.RFC3339), ashbySlug)
	return err
}

func GetAllCompaniesLastScraped() (map[string]*string, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT ashby_slug, last_scraped FROM companies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*string)
	for rows.Next() {
		var slug string
		var lastScraped sql.NullString
		if err := rows.Scan(&slug, &lastScraped); err != nil {
			continue
		}
		if lastScraped.Valid {
			result[slug] = &lastScraped.String
		} else {
			result[slug] = nil
		}
	}
	return result, nil
}

func GetJobById(id int) (*common.JobPayload, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	job := &common.JobPayload{Meta: common.JobMeta{Source: "ashby"}}
	var remote int
	var publishedAt, embeddingJSON sql.NullString
	var location, team, department, employmentType, compensation, jobURL, description, applyURL sql.NullString
	var jobID string

	err := db.QueryRow("SELECT id, job_id, content_hash, company, title, location, team, department, employment_type, remote, description, apply_url, job_url, compensation, published_at, embedding FROM jobs WHERE id = ?", id).Scan(
		&job.Id,
		&jobID,
		&job.Meta.ContentHash,
		&job.CompanyName,
		&job.JobName,
		&location,
		&team,
		&department,
		&employmentType,
		&remote,
		&description,
		&applyURL,
		&jobURL,
		&compensation,
		&publishedAt,
		&embeddingJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("job not found")
	}

	job.Meta.Remote = remote == 1
	if location.Valid {
		job.Meta.Location = location.String
	}
	if team.Valid {
		job.Meta.Team = team.String
	}
	if department.Valid {
		job.Meta.Department = department.String
	}
	if employmentType.Valid {
		job.Meta.EmploymentType = employmentType.String
	}
	if compensation.Valid {
		job.Meta.Compensation = compensation.String
	}
	if jobURL.Valid {
		job.Meta.JobURL = jobURL.String
	}
	if description.Valid {
		job.Description = description.String
	}
	if applyURL.Valid {
		job.ApplyLink = applyURL.String
	}
	if publishedAt.Valid && publishedAt.String != "" {
		if t, err := time.Parse(time.RFC3339, publishedAt.String); err == nil {
			job.Date = t
		}
	}
	if embeddingJSON.Valid && embeddingJSON.String != "" {
		var emb []float32
		if err := json.Unmarshal([]byte(embeddingJSON.String), &emb); err == nil {
			job.Embedding = emb
		}
	}

	return job, nil
}

func SaveJobs(jobs []*common.JobPayload) error {
	return InsertJobs(jobs)
}

func CreateTableIfNotExists() error {
	return InitDB()
}

func MarkJobAsDeleted(jobId string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := db.Exec("UPDATE jobs SET is_active = 0, status = 'deleted' WHERE job_id = ?", jobId)
	return err
}

func GetDeletedJobs(company string) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT content_hash FROM jobs WHERE company = ? AND is_active = 0", company)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func GetJobIdByContentHash(contentHash string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var jobID string
	err := db.QueryRow("SELECT job_id FROM jobs WHERE content_hash = ?", contentHash).Scan(&jobID)
	if err != nil {
		return "", fmt.Errorf("not found")
	}
	return jobID, nil
}

func MarkJobsInactive(company string, jobIds []string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	for _, id := range jobIds {
		_, err := db.Exec("UPDATE jobs SET is_active = 0, status = 'inactive' WHERE job_id = ? AND company = ?", id, company)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllJobsDeleted() ([]common.JobPayload, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT job_id, content_hash, company, title, location, team, department, employment_type, remote, description, apply_url, job_url, compensation, published_at, embedding FROM jobs WHERE is_active = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []common.JobPayload
	for rows.Next() {
		job := common.JobPayload{Meta: common.JobMeta{Source: "ashby"}}
		var remote int
		var publishedAt, embeddingJSON sql.NullString
		var location, team, department, employmentType, compensation, jobURL, description, applyURL sql.NullString

		err := rows.Scan(
			&job.Meta.ContentHash,
			&job.Meta.ContentHash,
			&job.CompanyName,
			&job.JobName,
			&location,
			&team,
			&department,
			&employmentType,
			&remote,
			&description,
			&applyURL,
			&jobURL,
			&compensation,
			&publishedAt,
			&embeddingJSON,
		)
		if err != nil {
			continue
		}

		job.Meta.Remote = remote == 1
		if location.Valid {
			job.Meta.Location = location.String
		}
		if team.Valid {
			job.Meta.Team = team.String
		}
		if department.Valid {
			job.Meta.Department = department.String
		}
		if employmentType.Valid {
			job.Meta.EmploymentType = employmentType.String
		}
		if compensation.Valid {
			job.Meta.Compensation = compensation.String
		}
		if jobURL.Valid {
			job.Meta.JobURL = jobURL.String
		}
		if description.Valid {
			job.Description = description.String
		}
		if applyURL.Valid {
			job.ApplyLink = applyURL.String
		}
		if publishedAt.Valid && publishedAt.String != "" {
			if t, err := time.Parse(time.RFC3339, publishedAt.String); err == nil {
				job.Date = t
			}
		}
		if embeddingJSON.Valid && embeddingJSON.String != "" {
			var emb []float32
			if err := json.Unmarshal([]byte(embeddingJSON.String), &emb); err == nil {
				job.Embedding = emb
			}
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func SaveDiff(deleted []string, updated []*common.JobPayload, new []*common.JobPayload) error {
	return nil
}

func GetDeletedDiff() ([]string, error) {
	return nil, nil
}

func GetAllJobIds() ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT job_id FROM jobs WHERE is_active = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func SetDeletedJobIds(jobIds []string) error {
	return nil
}

func ClearDeletedJobIds() error {
	return nil
}

func GetDeletedJobIds() ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT content_hash FROM deleted_jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func SetDeletedDiff(deleted []string) error {
	return nil
}

func GetDeletedDiffContent() ([]string, error) {
	return nil, nil
}

func MarkJobsAsDeleted(company string, activeJobIds []string) error {
	return nil
}

func GetAllAshbySlugToCompany() (map[string]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT ashby_slug, name FROM companies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var slug, name string
		if err := rows.Scan(&slug, &name); err == nil {
			result[slug] = name
		}
	}
	return result, nil
}

func SetLastScrapedForCompany(ashbySlug, company string, timestamp time.Time) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := db.Exec("INSERT INTO companies (ashby_slug, name, last_scraped) VALUES (?, ?, ?) ON CONFLICT(ashby_slug) DO UPDATE SET last_scraped = excluded.last_scraped", ashbySlug, company, timestamp.Format(time.RFC3339))
	return err
}

func GetLastScrapedForCompany(company string) (*time.Time, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var lastScraped string
	err := db.QueryRow("SELECT last_scraped FROM companies WHERE name = ?", company).Scan(&lastScraped)
	if err != nil {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, lastScraped); err == nil {
		return &t, nil
	}
	return nil, nil
}

func UpdateJobsInDB(deleted []string, updated []*common.JobPayload, new []*common.JobPayload) error {
	return nil
}

func GetUpdatedJobs() ([]*common.JobPayload, error) {
	return nil, nil
}
