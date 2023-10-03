package models

import "encoding/json"

type FrontRow []FrontRowElement

func UnmarshalFrontRow(data []byte) (FrontRow, error) {
	var r FrontRow
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FrontRow) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FrontRowElement struct {
	ID             int64         `json:"id"`
	Title          string        `json:"title"`
	DepartmentID   int64         `json:"departmentId"`
	DepartmentName string        `json:"departmentName"`
	JobLocations   []JobLocation `json:"jobLocations"`
	JobType        int64         `json:"jobType"`
}

type JobLocation struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	State       string `json:"state"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
}
