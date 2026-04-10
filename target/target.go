package target

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type CompanyInfo struct {
	Company        string
	AshbySlug      string
	Enabled        bool
	FrequencyHours int
}

var DefaultCompanies []CompanyInfo

func init() {
	loadCompaniesFromFile()
}

func loadCompaniesFromFile() {
	data, err := os.ReadFile("companies.json")
	if err != nil {
		loadCompaniesFromEnv()
		return
	}

	var companies []CompanyInfo
	if err := json.Unmarshal(data, &companies); err != nil {
		loadCompaniesFromEnv()
		return
	}

	DefaultCompanies = companies
}

func loadCompaniesFromEnv() {
	envCompanies := os.Getenv("ASHBY_COMPANIES")
	if envCompanies != "" {
		var parsed []CompanyInfo
		if err := json.Unmarshal([]byte(envCompanies), &parsed); err == nil {
			DefaultCompanies = parsed
			return
		}
	}

	commaSep := os.Getenv("ASHBY_COMPANIES_COMMA")
	if commaSep != "" {
		parts := strings.Split(commaSep, ",")
		for _, p := range parts {
			kv := strings.Split(strings.TrimSpace(p), ":")
			if len(kv) == 2 {
				DefaultCompanies = append(DefaultCompanies, CompanyInfo{
					Company: kv[0], AshbySlug: kv[1], Enabled: true, FrequencyHours: 12,
				})
			}
		}
	}
}

func GetEnabledCompanies() []CompanyInfo {
	companies := []CompanyInfo{}

	envCompanies := os.Getenv("ASHBY_COMPANIES")
	if envCompanies != "" {
		var parsed []CompanyInfo
		if err := json.Unmarshal([]byte(envCompanies), &parsed); err == nil {
			companies = append(companies, parsed...)
		}
	}

	commaSep := os.Getenv("ASHBY_COMPANIES_COMMA")
	if commaSep != "" {
		parts := strings.Split(commaSep, ",")
		for _, p := range parts {
			kv := strings.Split(strings.TrimSpace(p), ":")
			if len(kv) == 2 {
				companies = append(companies, CompanyInfo{
					Company: kv[0], AshbySlug: kv[1], Enabled: true, FrequencyHours: 12,
				})
			}
		}
	}

	if len(companies) == 0 {
		for _, c := range DefaultCompanies {
			if c.Enabled {
				companies = append(companies, c)
			}
		}
	}

	return companies
}

func GetAllCompanies() []CompanyInfo {
	return DefaultCompanies
}

func GetDueCompanies(lastScraped map[string]*string, allCompanies []CompanyInfo) []CompanyInfo {
	var due []CompanyInfo
	for _, c := range allCompanies {
		if !c.Enabled {
			continue
		}
		last := lastScraped[c.AshbySlug]
		if last == nil {
			due = append(due, c)
			continue
		}
		if shouldScrape(*last, c.FrequencyHours) {
			due = append(due, c)
		}
	}
	return due
}

func shouldScrape(lastScraped string, frequencyHours int) bool {
	if lastScraped == "" {
		return true
	}
	t, err := time.Parse(time.RFC3339, lastScraped)
	if err != nil {
		return true
	}
	return time.Since(t) > time.Duration(frequencyHours)*time.Hour
}
