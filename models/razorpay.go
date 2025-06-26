package models

type RazorpayData struct {
	Jobs []RazorpayJobs `json:"jobs"`
	Meta RazorpayMeta   `json:"meta"`
}

type RazorpayJobs struct {
	AbsoluteURL    string               `json:"absolute_url"`
	DataCompliance []DataCompliance     `json:"data_compliance"`
	Employment     string               `json:"employment"`
	InternalJobID  int64                `json:"internal_job_id"`
	Location       RazorpayLocation     `json:"location"`
	Metadata       []Metadata           `json:"metadata"`
	ID             int64                `json:"id"`
	UpdatedAt      string               `json:"updated_at"`
	RequisitionID  string               `json:"requisition_id"`
	Title          string               `json:"title"`
	CompanyName    string               `json:"company_name"`
	FirstPublished string               `json:"first_published"`
	Content        string               `json:"content"`
	Departments    []RazorpayDepartment `json:"departments"`
	Offices        []RazorpayOffice     `json:"offices"`
}

type DataCompliance struct {
	Type                          string `json:"type"`
	RequiresConsent               bool   `json:"requires_consent"`
	RequiresProcessingConsent     bool   `json:"requires_processing_consent"`
	RequiresRetentionConsent      bool   `json:"requires_retention_consent"`
	RetentionPeriod               *int   `json:"retention_period"`
	DemographicDataConsentApplies bool   `json:"demographic_data_consent_applies"`
}

type RazorpayLocation struct {
	Name string `json:"name"`
}

type Metadata struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"` // Can be string or []string
	ValueType string      `json:"value_type"`
}

type RazorpayDepartment struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	ChildIDs []int64 `json:"child_ids"`
	ParentID int64   `json:"parent_id"`
}

type RazorpayOffice struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Location *string `json:"location"`
	ChildIDs []int64 `json:"child_ids"`
	ParentID *int64  `json:"parent_id"`
}

type RazorpayMeta struct {
	Total int `json:"total"`
}
