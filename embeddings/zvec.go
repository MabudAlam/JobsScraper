package embeddings

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ashbyimpl/common"

	"github.com/oliveagle/zvec-go"
)

const (
	CollectionName = "jobs"
	Dimension      = 1536
)

type UpsertResult string

const (
	Inserted  UpsertResult = "inserted"
	Updated   UpsertResult = "updated"
	Unchanged UpsertResult = "unchanged"
)

var (
	coll          *zvec.Collection
	zvecPath      string
	dbInitOnce    sync.Once
	dbInitialized bool
)

func InitDB() error {
	var initErr error
	dbInitOnce.Do(func() {
		zvec.Init(zvec.DefaultConfig())

		zvecPath = os.Getenv("ZVEC_PATH")
		if zvecPath == "" {
			zvecPath = filepath.Join(os.Getenv("PWD"), "zvec_data")
		}

		if err := os.MkdirAll(zvecPath, 0755); err != nil {
			initErr = fmt.Errorf("failed to create zvec data dir: %w", err)
			return
		}

		logPath := filepath.Join(zvecPath, "logs")
		if err := os.MkdirAll(logPath, 0755); err != nil {
			initErr = fmt.Errorf("failed to create log dir: %w", err)
			return
		}

		schema := zvec.NewCollectionSchema(CollectionName)
		schema.AddField(zvec.NewFieldSchema("job_id", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("company", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("title", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("location", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("team", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("department", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("employment_type", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("remote", zvec.DataTypeBool))
		schema.AddField(zvec.NewFieldSchema("description", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("apply_url", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("job_url", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("compensation", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("content_hash", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("is_active", zvec.DataTypeBool))
		schema.AddField(zvec.NewFieldSchema("status", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("published_at", zvec.DataTypeString))
		schema.AddField(zvec.NewFieldSchema("scraped_at", zvec.DataTypeString))

		schema.AddVectorField(
			zvec.NewVectorSchema("embedding", zvec.DataTypeVectorFP32, Dimension).
				WithMetricType(zvec.MetricTypeCOSINE),
		)

		collPath := filepath.Join(zvecPath, CollectionName)
		var err error
		coll, err = zvec.CreateAndOpen(collPath, schema, nil)
		if err != nil {
			initErr = fmt.Errorf("failed to create/open collection: %w", err)
			return
		}

		loadDocsFromDisk()

		dbInitialized = true
		log.Println("zvec database initialized at:", zvecPath)
	})
	return initErr
}

func loadDocsFromDisk() {
	if coll == nil {
		return
	}

	docsDir := filepath.Join(zvecPath, CollectionName, "docs")
	entries, err := os.ReadDir(docsDir)
	if err != nil {
		log.Printf("[zvec] Warning: could not read docs directory: %v", err)
		return
	}

	var docs []*zvec.Document
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := entry.Name()[:len(entry.Name())-5]
		doc, err := coll.Get(id)
		if err == nil && doc != nil {
			docs = append(docs, doc)
		}
	}

	if len(docs) > 0 {
		_, err := coll.InsertBatch(docs)
		if err != nil {
			log.Printf("[zvec] Warning: failed to load docs into memory: %v", err)
		} else {
			log.Printf("[zvec] Loaded %d documents from disk", len(docs))
		}
	}
}

func Close() error {
	if coll != nil {
		return coll.Close()
	}
	return nil
}

func CreateCollectionIfNotExists() error {
	return InitDB()
}

func docToJob(doc *zvec.Document) common.JobPayload {
	job := common.JobPayload{}

	if v, ok := doc.GetField("job_id"); ok {
		job.Meta.ContentHash = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("company"); ok {
		job.CompanyName = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("title"); ok {
		job.JobName = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("location"); ok {
		job.Meta.Location = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("team"); ok {
		job.Meta.Team = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("department"); ok {
		job.Meta.Department = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("employment_type"); ok {
		job.Meta.EmploymentType = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("remote"); ok {
		if b, ok := v.(bool); ok {
			job.Meta.Remote = b
		}
	}
	if v, ok := doc.GetField("description"); ok {
		job.Description = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("apply_url"); ok {
		job.ApplyLink = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("job_url"); ok {
		job.Meta.JobURL = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("compensation"); ok {
		job.Meta.Compensation = fmt.Sprintf("%v", v)
	}
	if v, ok := doc.GetField("published_at"); ok {
		if s, ok := v.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				job.Date = t
			}
		}
	}
	job.Meta.Source = "ashby"

	return job
}

func jobToDoc(job *common.JobPayload, id string) *zvec.Document {
	doc := zvec.NewDocument(id)
	doc.SetField("job_id", job.Meta.ContentHash)
	doc.SetField("company", job.CompanyName)
	doc.SetField("title", job.JobName)
	doc.SetField("location", job.Meta.Location)
	doc.SetField("team", job.Meta.Team)
	doc.SetField("department", job.Meta.Department)
	doc.SetField("employment_type", job.Meta.EmploymentType)
	doc.SetField("remote", job.Meta.Remote)
	doc.SetField("description", job.Description)
	doc.SetField("apply_url", job.ApplyLink)
	doc.SetField("job_url", job.Meta.JobURL)
	doc.SetField("compensation", job.Meta.Compensation)
	doc.SetField("content_hash", job.Meta.ContentHash)
	doc.SetField("is_active", true)
	doc.SetField("status", "new")

	if !job.Date.IsZero() {
		doc.SetField("published_at", job.Date.Format(time.RFC3339))
	}
	doc.SetField("scraped_at", time.Now().Format(time.RFC3339))

	embedding := job.Embedding
	if embedding == nil {
		embedding = make([]float32, Dimension)
	}
	doc.SetVector("embedding", embedding)

	return doc
}

func InsertJobs(jobs []*common.JobPayload) error {
	if len(jobs) == 0 || coll == nil {
		return nil
	}

	docs := make([]*zvec.Document, len(jobs))
	for i, job := range jobs {
		docs[i] = jobToDoc(job, fmt.Sprintf("job_%d", i))
	}

	_, err := coll.InsertBatch(docs)
	return err
}

func UpsertJob(job *common.JobPayload) (UpsertResult, error) {
	if coll == nil {
		return Inserted, fmt.Errorf("collection not initialized")
	}

	doc := jobToDoc(job, job.Meta.ContentHash)
	err := coll.Upsert(doc)
	if err != nil {
		return Inserted, err
	}
	return Updated, nil
}

func MarkRemovedJobs(company string, activeJobIds []string) ([]string, error) {
	if coll == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	allDocs, err := coll.ListIDs()
	if err != nil {
		return nil, err
	}

	activeSet := make(map[string]bool)
	for _, id := range activeJobIds {
		activeSet[id] = true
	}

	var removed []string
	for _, id := range allDocs {
		doc, err := coll.Get(id)
		if err != nil {
			continue
		}
		if comp, ok := doc.GetField("company"); ok {
			if compStr, ok := comp.(string); ok && compStr == company {
				if active, ok := doc.GetField("is_active"); ok {
					if activeBool, ok := active.(bool); ok && activeBool {
						if !activeSet[id] {
							if err := coll.Delete(id); err == nil {
								removed = append(removed, id)
							}
						}
					}
				}
			}
		}
	}
	return removed, nil
}

func GetAllActiveJobs() ([]common.JobPayload, error) {
	if coll == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	allDocs, err := coll.ListIDs()
	if err != nil {
		return nil, err
	}

	var jobs []common.JobPayload
	for _, id := range allDocs {
		doc, err := coll.Get(id)
		if err != nil {
			continue
		}
		if active, ok := doc.GetField("is_active"); ok {
			if activeBool, ok := active.(bool); ok && activeBool {
				jobs = append(jobs, docToJob(doc))
			}
		}
	}
	return jobs, nil
}

func GetActiveJobsPaginated(offset, limit int, search, company, location, sort string) ([]common.JobPayload, int, error) {
	if coll == nil {
		return nil, 0, fmt.Errorf("collection not initialized")
	}

	allDocs, err := coll.ListIDs()
	if err != nil {
		return nil, 0, err
	}

	var jobs []common.JobPayload
	for _, id := range allDocs {
		doc, err := coll.Get(id)
		if err != nil {
			continue
		}
		if active, ok := doc.GetField("is_active"); ok {
			if activeBool, ok := active.(bool); ok && activeBool {
				job := docToJob(doc)

				if company != "" {
					if job.CompanyName != company {
						continue
					}
				}
				if location != "" {
					if job.Meta.Location != location {
						continue
					}
				}
				if search != "" {
					found := false
					searchLower := lowercase(search)
					if containsLower(job.JobName, searchLower) {
						found = true
					} else if containsLower(job.Description, searchLower) {
						found = true
					} else if containsLower(job.Meta.Location, searchLower) {
						found = true
					}
					if !found {
						continue
					}
				}

				jobs = append(jobs, job)
			}
		}
	}

	total := len(jobs)
	if offset < len(jobs) {
		end := offset + limit
		if end > len(jobs) {
			end = len(jobs)
		}
		jobs = jobs[offset:end]
	} else {
		jobs = nil
	}

	return jobs, total, nil
}

func lowercase(s string) string {
	return s
}

func containsLower(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func SearchJobs(queryVector []float32, limit int) ([]map[string]interface{}, error) {
	if coll == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	query := zvec.NewVectorQueryByVector("embedding", queryVector).WithTopK(limit)
	results, err := coll.Search(query)
	if err != nil {
		return nil, err
	}

	var searchResults []map[string]interface{}
	for _, r := range results {
		result := make(map[string]interface{})
		result["id"] = r.ID
		result["score"] = r.Score
		if r.Document != nil {
			for k, v := range r.Document.Fields {
				result[k] = v
			}
		}
		searchResults = append(searchResults, result)
	}
	return searchResults, nil
}

func SearchJobsWithFilter(queryVector []float32, limit int, filter string) ([]map[string]interface{}, error) {
	return SearchJobs(queryVector, limit)
}

func GetAllLocations() ([]string, error) {
	if coll == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	seen := make(map[string]bool)
	var locations []string

	allDocs, err := coll.ListIDs()
	if err != nil {
		return nil, err
	}

	for _, id := range allDocs {
		doc, err := coll.Get(id)
		if err != nil {
			continue
		}
		if active, ok := doc.GetField("is_active"); ok {
			if activeBool, ok := active.(bool); ok && activeBool {
				if loc, ok := doc.GetField("location"); ok {
					if locStr := fmt.Sprintf("%v", loc); locStr != "" && !seen[locStr] {
						seen[locStr] = true
						locations = append(locations, locStr)
					}
				}
			}
		}
	}
	return locations, nil
}

func GetAllCompanies() ([]string, error) {
	if coll == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	seen := make(map[string]bool)
	var companies []string

	allDocs, err := coll.ListIDs()
	if err != nil {
		return nil, err
	}

	for _, id := range allDocs {
		doc, err := coll.Get(id)
		if err != nil {
			continue
		}
		if active, ok := doc.GetField("is_active"); ok {
			if activeBool, ok := active.(bool); ok && activeBool {
				if comp, ok := doc.GetField("company"); ok {
					if compStr := fmt.Sprintf("%v", comp); compStr != "" && !seen[compStr] {
						seen[compStr] = true
						companies = append(companies, compStr)
					}
				}
			}
		}
	}
	return companies, nil
}

func GetClient() interface{} {
	return coll
}

func SaveSnapshot(job *common.JobPayload) error {
	return nil
}

func UpsertCompany(name, ashbySlug string) error {
	return nil
}

func UpdateLastScraped(ashbySlug string) error {
	return nil
}

func GetAllCompaniesLastScraped() (map[string]*string, error) {
	result := make(map[string]*string)
	return result, nil
}

func GetJobById(id int) (*common.JobPayload, error) {
	if coll == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	doc, err := coll.Get(fmt.Sprintf("job_%d", id))
	if err != nil {
		return nil, fmt.Errorf("job not found")
	}
	job := docToJob(doc)
	return &job, nil
}

func SaveJobs(jobs []*common.JobPayload) error {
	return InsertJobs(jobs)
}

func CreateTableIfNotExists() error {
	return InitDB()
}

func MarkJobAsDeleted(jobId string) error {
	if coll == nil {
		return fmt.Errorf("collection not initialized")
	}
	return coll.Delete(jobId)
}

func GetDeletedJobs(company string) ([]string, error) {
	return nil, nil
}

func GetJobIdByContentHash(contentHash string) (string, error) {
	if coll == nil {
		return "", fmt.Errorf("collection not initialized")
	}

	doc, err := coll.Get(contentHash)
	if err != nil {
		return "", fmt.Errorf("not found")
	}
	if v, ok := doc.GetField("job_id"); ok {
		return fmt.Sprintf("%v", v), nil
	}
	return contentHash, nil
}

func MarkJobsInactive(company string, jobIds []string) error {
	return nil
}

func GetAllJobsDeleted() ([]common.JobPayload, error) {
	return nil, nil
}

func SaveDiff(deleted []string, updated []*common.JobPayload, new []*common.JobPayload) error {
	return nil
}

func GetDeletedDiff() ([]string, error) {
	return nil, nil
}

func GetAllJobIds() ([]string, error) {
	if coll == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	allDocs, err := coll.ListIDs()
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, id := range allDocs {
		doc, err := coll.Get(id)
		if err != nil {
			continue
		}
		if active, ok := doc.GetField("is_active"); ok {
			if activeBool, ok := active.(bool); ok && activeBool {
				if jobId, ok := doc.GetField("job_id"); ok {
					ids = append(ids, fmt.Sprintf("%v", jobId))
				}
			}
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
	return nil, nil
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
	return make(map[string]string), nil
}

func SetLastScrapedForCompany(ashbySlug, company string, timestamp time.Time) error {
	return nil
}

func GetLastScrapedForCompany(company string) (*time.Time, error) {
	return nil, nil
}

func UpdateJobsInDB(deleted []string, updated []*common.JobPayload, new []*common.JobPayload) error {
	return nil
}

func GetUpdatedJobs() ([]*common.JobPayload, error) {
	return nil, nil
}

func contentHashToId(contentHash string) (string, error) {
	return contentHash, nil
}
