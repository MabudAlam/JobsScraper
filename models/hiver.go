package models

import "encoding/json"

// HiverPosting represents a job posting on Hiver.
type HiverPosting struct {
	Title       string `json:"title"`
	Location    string `json:"location"`
	Description string `json:"description"`
	ApplyURL    string `json:"applyUrl"`
	// Add other fields as needed
}

// UnmarshalHiverPosting parses JSON data into a HiverPosting structure.
func UnmarshalHiverPosting(data []byte) (HiverPosting, error) {
	var posting HiverPosting
	err := json.Unmarshal(data, &posting)
	return posting, err
}

// Marshal converts a HiverPosting structure into JSON data.
func (posting *HiverPosting) Marshal() ([]byte, error) {
	return json.Marshal(posting)
}
