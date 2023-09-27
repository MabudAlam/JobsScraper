

package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
"goscraper/models")

func GetAllJobsFromSqlite(c *gin.Context) {
	db, err := sql.Open("sqlite3", "jobs.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	query := `
		SELECT title, location, created_at, company, apply_url, image_url
		FROM jobs;
	`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var allPostings []models.Job

	for rows.Next() {
		var posting models.Job
		err := rows.Scan(&posting.Title, &posting.Location, &posting.CreatedAt, &posting.Company, &posting.ApplyURL, &posting.ImageUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		allPostings = append(allPostings, posting)
	}

	// Return the retrieved job postings as a JSON response.
	c.JSON(http.StatusOK, allPostings)
}
