package config

import (
	"sync"

	"github.com/joho/godotenv"
)

var (
	cachedConfig *Config
	configMu     sync.RWMutex
	configOnce   sync.Once
	configErr    error
)

func LoadConfig() *Config {
	configMu.RLock()
	if cachedConfig != nil {
		configMu.RUnlock()
		return cachedConfig
	}
	configMu.RUnlock()

	configMu.Lock()
	defer configMu.Unlock()
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
