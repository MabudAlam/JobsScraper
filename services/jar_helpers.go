package services

import (
	"fmt"
	"goscraper/models"
)

func buildJarDescription(item models.JarJob) string {
	desc := "Title: " + item.Title + "\n" +
		"Description: " + item.Description + "\n" +
		"Department: " + item.DepartmentName + "\n" +
		"Location: " + item.Location + "\n" +
		"Company: " + item.CompanyName + "\n" +
		"Job Type: " + item.JobType + "\n" +
		"Workplace Type: " + item.WorkplaceType + "\n" +
		"Remote: " + boolToYesNo(item.Remote) + "\n" +
		"Salary Type: " + item.SalaryType + "\n" +
		"Min Salary: " + formatFloat(item.MinSalary) + "\n" +
		"Max Salary: " + formatFloat(item.MaxSalary) + "\n" +
		"Max Experience: " + formatFloat(item.MaxExperience) + "\n" +
		"Skills: " + joinStrings(item.Skill, ", ") + "\n" +
		"Apply Link: " + item.ApplyLink + "\n" +
		"Created At: " + item.CreatedAt
	return desc
}

func boolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func formatFloat(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%.2f", f)
}

func joinStrings(arr []string, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result += sep + arr[i]
	}
	return result
}
