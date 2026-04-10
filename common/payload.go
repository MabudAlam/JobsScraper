package common

import "time"

type JobPayload struct {
	Id          int       `json:"id"`
	JobName     string    `json:"jobName"`
	Description string    `json:"description"`
	Embedding   []float32 `json:"embedding,omitempty"`
	Keywords    []string  `json:"keywords,omitempty"`
	Date        time.Time `json:"date"`
	ApplyLink   string    `json:"applyLink"`
	CompanyName string    `json:"companyName"`
	Meta        JobMeta   `json:"meta"`
}

type JobMeta struct {
	Location       string `json:"location,omitempty"`
	Remote         bool   `json:"remote,omitempty"`
	Department     string `json:"department,omitempty"`
	Team           string `json:"team,omitempty"`
	EmploymentType string `json:"employmentType,omitempty"`
	Compensation   string `json:"compensation,omitempty"`
	JobURL         string `json:"jobUrl,omitempty"`
	ContentHash    string `json:"contentHash,omitempty"`
	Source         string `json:"source"`
	RawData        any    `json:"rawData,omitempty"`
}

type ChangeType string

const (
	ChangeNew     ChangeType = "JOB_NEW"
	ChangeUpdated ChangeType = "JOB_UPDATED"
	ChangeRemoved ChangeType = "JOB_REMOVED"
)

type Change struct {
	Type ChangeType  `json:"type"`
	Job  *JobPayload `json:"job"`
}
