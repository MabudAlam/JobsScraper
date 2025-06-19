// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    atlanJobData, err := UnmarshalAtlanJobData(bytes)
//    bytes, err = atlanJobData.Marshal()

package models

import "encoding/json"

func UnmarshalAtlanJobData(data []byte) (AtlanJobData, error) {
	var r AtlanJobData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AtlanJobData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type AtlanJobData struct {
	Data Data `json:"data"`
}

type Data struct {
	JobBoard JobBoard `json:"jobBoard"`
}

type JobBoard struct {
	Teams       []AtlanTeam       `json:"teams"`
	JobPostings []JobPosting `json:"jobPostings"`
	Typename    string       `json:"__typename"`
}

type AtlanTeam struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	ParentTeamID interface{} `json:"parentTeamId"`
	Typename     string      `json:"__typename"`
}

type JobPosting struct {
	ID                      string              `json:"id"`
	Title                   string              `json:"title"`
	TeamID                  string              `json:"teamId"`
	LocationID              string              `json:"locationId"`
	LocationName            string              `json:"locationName"`
	WorkplaceType           *string             `json:"workplaceType"` // nullable
	EmploymentType          *string             `json:"employmentType"` // nullable
	SecondaryLocations      []SecondaryLocation `json:"secondaryLocations"`
	CompensationTierSummary interface{}         `json:"compensationTierSummary"`
	Typename                string              `json:"__typename"`
}

type SecondaryLocation struct {
	LocationID   string `json:"locationId"`
	LocationName string `json:"locationName"`
	Typename     string `json:"__typename"`
}
