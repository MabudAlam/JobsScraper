package models

import "encoding/json"

func UnmarshalNiyoSolution(data []byte) (NiyoSolution, error) {
	var r NiyoSolution
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NiyoSolution) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NiyoSolution struct {
	Status  string  `json:"status"`
	Message Message `json:"message"`
}

type NiyoSolutionMessage struct {
	Jobscount int64 `json:"jobscount"`
	Jobs      []Job `json:"jobs"`
}

type NiyoSolutionJob struct {
	ID                     string                `json:"id"`
	DesignationDisplayName string                `json:"designation_display_name"`
	ExperienceFromNum      int64                 `json:"experience_from_num"`
	ExperienceToNum        int64                 `json:"experience_to_num"`
	CreatedOn              string                `json:"created_on"`
	OfficelocationShowArr  OfficelocationShowArr `json:"officelocation_show_arr"`
	OfficelocationArr      []OfficelocationArr   `json:"officelocation_arr"`
	JobPostingOn           int64                 `json:"job_posting_on"`
	Department             string                `json:"department"`
	DepartmentID           string                `json:"department_id"`
	EmpType                EmpType               `json:"emp_type"`
	EmpTypeID              EmpTypeID             `json:"emp_type_id"`
	FunctionalArea         string                `json:"functional_area"`
	FunctionalAreaID       string                `json:"functional_area_id"`
	Title                  string                `json:"title"`
	ToolTipLocations       OfficelocationShowArr `json:"tool_tip_locations"`
	Timezone               Timezone              `json:"timezone"`
}

type NiyoSolutionEmpType string

const (
	NiyoSolutionFte EmpType = "FTE"
)

type NiyoSolutionEmpTypeID string

const (
	A63317545D4E33 EmpTypeID = "a63317545d4e33"
)

type NiyoSolutionOfficelocationArr string

const (
	A635Bf581E6F26 OfficelocationArr = "a635bf581e6f26"
)

type NiyoSolutionOfficelocationShowArr string

const (
	NiyoSolutionBengaluruKarnatakaIndia OfficelocationShowArr = "Bengaluru, Karnataka, India "
)

type NiyoSolutionTimezone string

const (
	AsiaKolkata Timezone = "Asia/Kolkata"
)
