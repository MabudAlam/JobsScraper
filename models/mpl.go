// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    mPLJobData, err := UnmarshalMPLJobData(bytes)
//    bytes, err = mPLJobData.Marshal()

package models

import "time"

import "encoding/json"

func UnmarshalMPLJobData(data []byte) (MPLJobData, error) {
	var r MPLJobData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MPLJobData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MPLJobData struct {
	Status  string  `json:"status"`
	Message Message `json:"message"`
}

type Message struct {
	Jobscount int64 `json:"jobscount"`
	Jobs      []MplJob `json:"jobs"`
}

type MplJob struct {
	ID                     string                `json:"id"`
	DesignationDisplayName string                `json:"designation_display_name"`
	ExperienceFromNum      string                `json:"experience_from_num"`
	ExperienceToNum        string                `json:"experience_to_num"`
	CreatedOn              time.Time             `json:"created_on"`
	OfficelocationShowArr  OfficelocationShowArr `json:"officelocation_show_arr"`
	OfficelocationArr      string                `json:"officelocation_arr"`
	JobPostingOn           int64                 `json:"job_posting_on"`
	Department             string                `json:"department"`
	DepartmentID           string                `json:"department_id"`
	EmpType                EmpType               `json:"emp_type"`
	EmpTypeID              EmpTypeID             `json:"emp_type_id"`
	FunctionalArea         FunctionalArea        `json:"functional_area"`
	FunctionalAreaID       FunctionalAreaID      `json:"functional_area_id"`
	Title                  string                `json:"title"`
	ToolTipLocations       []ToolTipLocation     `json:"tool_tip_locations"`
	Timezone               Timezone              `json:"timezone"`
}

type EmpType string

const (
	Contract EmpType = "Contract"
	Fte      EmpType = "FTE"
	Intern   EmpType = "Intern"
)

type EmpTypeID string

const (
	The5Ecbc7019233C EmpTypeID = "5ecbc7019233c"
	The5Ecbc7079278A EmpTypeID = "5ecbc7079278a"
	The5Ecbc70Eb2E96 EmpTypeID = "5ecbc70eb2e96"
)

type FunctionalArea string

const (
	BusinessEnablementBE FunctionalArea = "Business Enablement (BE)"
	ProductPROD          FunctionalArea = "Product (PROD)"
	RevenueGrowthRG      FunctionalArea = "Revenue & Growth (R&G)"
	TechTech             FunctionalArea = "TECH (TECH)"
)

type FunctionalAreaID string

const (
	A62203F695B355   FunctionalAreaID = "a62203f695b355"
	A622229Afb2Ec4   FunctionalAreaID = "a622229afb2ec4"
	A6238Ccd093Ca7   FunctionalAreaID = "a6238ccd093ca7"
	The5Ff3Ab07997Bd FunctionalAreaID = "5ff3ab07997bd"
)

type OfficelocationShowArr string

const (
	MultipleLocations                            OfficelocationShowArr = "Multiple locations"
	OfficelocationShowArrBengaluruKarnatakaIndia OfficelocationShowArr = "Bengaluru, Karnataka, India "
)

type Timezone string

const (
	AsiaJakarta Timezone = "Asia/Jakarta"
)

type ToolTipLocation string

const (
	PuneMaharashtraIndia                   ToolTipLocation = "Pune, Maharashtra, India "
	ToolTipLocationBengaluruKarnatakaIndia ToolTipLocation = "Bengaluru, Karnataka, India "
)
