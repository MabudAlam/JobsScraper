package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	embeddings "jobscraper/db"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 80
	maxLimit     = 200
)

func GetAllJobs(c *gin.Context) {
	offset := 0
	limit := defaultLimit
	search := c.Query("search")
	company := c.Query("company")
	location := c.Query("location")
	sort := c.DefaultQuery("sort", "newest")

	allowedSorts := map[string]bool{"newest": true, "date": true, "title": true}
	if !allowedSorts[sort] {
		sort = "newest"
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	jobs, total, err := embeddings.GetActiveJobsPaginated(offset, limit, search, company, location, sort)
	if err != nil {
		slog.Error("failed to get jobs", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"jobs":   jobs,
		"offset": offset,
		"limit":  limit,
		"total":  total,
	})
}

func GetCompanies(c *gin.Context) {
	companies, err := embeddings.GetAllCompanies()
	if err != nil {
		slog.Error("failed to get companies", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve companies"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"companies": companies})
}

func GetLocations(c *gin.Context) {
	locations, err := embeddings.GetAllLocations()
	if err != nil {
		slog.Error("failed to get locations", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve locations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"locations": locations})
}

func GetJobByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := embeddings.GetJobById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}
