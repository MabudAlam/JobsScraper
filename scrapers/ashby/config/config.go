package config

import (
	"os"

	"github.com/joho/godotenv"
)

var cachedConfig *Config

func LoadConfig() *Config {
	if cachedConfig != nil {
		return cachedConfig
	}

	_ = godotenv.Load()

	cfg := &Config{
		Fetch: FetchConfig{
			DelayBetweenCompaniesMin: 100,
			DelayBetweenCompaniesMax: 500,
			MaxRetries:               3,
		},
		Notify: NotifyConfig{
			CLI:        true,
			Markdown:   false,
			ReportsDir: "./reports",
		},
		Intelligence: IntelligenceConfig{
			MinScore:          10,
			FreshnessBoostHrs: 72,
			FreshnessBoost:    5,
			LocationBoost:     3,
			RemoteBoost:       2,
			DepartmentBoost:   3,
			Keywords: map[string]int{
				"senior":    5,
				"principal": 5,
				"staff":     5,
				"engineer":  3,
				"developer": 3,
				"architect": 4,
				"manager":   3,
				"lead":      4,
			},
			Locations: []string{
				"san francisco",
				"new york",
				"remote",
				"us",
			},
			Departments: []string{
				"engineering",
				"product",
				"design",
			},
		},
	}

	if v := os.Getenv("FETCH_DELAY_MIN"); v != "" {
		cfg.Fetch.DelayBetweenCompaniesMin = 100
	}
	if v := os.Getenv("FETCH_DELAY_MAX"); v != "" {
		cfg.Fetch.DelayBetweenCompaniesMax = 500
	}
	if v := os.Getenv("MIN_SCORE"); v != "" {
		cfg.Intelligence.MinScore = 10
	}

	cachedConfig = cfg
	return cfg
}

type Config struct {
	Fetch        FetchConfig
	Notify       NotifyConfig
	Intelligence IntelligenceConfig
}

type FetchConfig struct {
	DelayBetweenCompaniesMin int
	DelayBetweenCompaniesMax int
	MaxRetries               int
}

type NotifyConfig struct {
	CLI        bool
	Markdown   bool
	ReportsDir string
}

type IntelligenceConfig struct {
	MinScore          int
	FreshnessBoostHrs int
	FreshnessBoost    int
	LocationBoost     int
	RemoteBoost       int
	DepartmentBoost   int
	Keywords          map[string]int
	Locations         []string
	Departments       []string
}

type IntelligenceRules = IntelligenceConfig
