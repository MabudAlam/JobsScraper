package models

// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    zohoJobs, err := UnmarshalZohoJobs(bytes)
//    bytes, err = zohoJobs.Marshal()

import "encoding/json"

func UnmarshalZohoJobs(data []byte) (ZohoJobs, error) {
	var r ZohoJobs
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ZohoJobs) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ZohoJobs struct {
	Code string  `json:"code"`
	Data []Datum `json:"data"`
	Info Info    `json:"info"`
}

type Datum struct {
	JobDescription string `json:"Job_Description"`
	JobType        string `json:"Job_Type"`
	JobOpeningName string `json:"Job_Opening_Name"`
	PostingTitle   string `json:"Posting_Title"`
	Country1       string `json:"Country1"`
	URL            string `json:"$url"`
	ID             string `json:"id"`
}

type Info struct {
	PageName string  `json:"page_name"`
	Groups   []Facet `json:"groups"`
	Fields   []Facet `json:"fields"`
	Facets   []Facet `json:"facets"`
}

type Facet struct {
	APIName    string `json:"api_name"`
	FieldLabel string `json:"field_label"`
}
