package utility

import (
	"goscraper/models"
)

// Helper function to check if a job posting already exists in the 'postings' slice
func CheckDuplicates(postings []models.Job, posting models.Job) bool {
	for _, p := range postings {
		if p.Title == posting.Title && p.Location == posting.Location {
			return true
		}
	}
	return false
}
