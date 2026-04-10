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

var DefaultCompanies = []CompanyInfo{
	{Company: "Ashby", AshbySlug: "Ashby", Enabled: true, FrequencyHours: 12},
	{Company: "OpenAI", AshbySlug: "openai", Enabled: true, FrequencyHours: 12},
	{Company: "Notion", AshbySlug: "Notion", Enabled: true, FrequencyHours: 12},
	{Company: "Ramp", AshbySlug: "ramp", Enabled: true, FrequencyHours: 12},
	{Company: "Deel", AshbySlug: "deel", Enabled: true, FrequencyHours: 12},
	{Company: "Linear", AshbySlug: "linear", Enabled: true, FrequencyHours: 12},
	{Company: "Cursor", AshbySlug: "cursor", Enabled: true, FrequencyHours: 12},
	{Company: "Snowflake", AshbySlug: "snowflake", Enabled: true, FrequencyHours: 12},
	{Company: "Vanta", AshbySlug: "vanta", Enabled: true, FrequencyHours: 12},
	{Company: "PostHog", AshbySlug: "posthog", Enabled: true, FrequencyHours: 12},
	{Company: "Replit", AshbySlug: "replit", Enabled: true, FrequencyHours: 12},
	{Company: "Supabase", AshbySlug: "supabase", Enabled: true, FrequencyHours: 12},
	{Company: "Zapier", AshbySlug: "zapier", Enabled: true, FrequencyHours: 12},
	{Company: "Harvey", AshbySlug: "harvey", Enabled: true, FrequencyHours: 12},
	{Company: "Stytch", AshbySlug: "stytch", Enabled: true, FrequencyHours: 12},
	{Company: "1Password", AshbySlug: "1password", Enabled: true, FrequencyHours: 12},
	{Company: "Deliveroo", AshbySlug: "deliveroo", Enabled: true, FrequencyHours: 12},
	{Company: "Trainline", AshbySlug: "trainline", Enabled: true, FrequencyHours: 12},
	{Company: "Cohere", AshbySlug: "cohere", Enabled: true, FrequencyHours: 12},
	{Company: "Lemonade", AshbySlug: "lemonade", Enabled: true, FrequencyHours: 12},
	{Company: "Gorgias", AshbySlug: "gorgias", Enabled: true, FrequencyHours: 12},
	{Company: "UiPath", AshbySlug: "uipath", Enabled: true, FrequencyHours: 12},
	{Company: "HackerOne", AshbySlug: "hackerone", Enabled: true, FrequencyHours: 12},
	{Company: "FullStory", AshbySlug: "fullstory", Enabled: true, FrequencyHours: 12},
	{Company: "Quora", AshbySlug: "quora", Enabled: true, FrequencyHours: 12},
	{Company: "Multiverse", AshbySlug: "multiverse", Enabled: true, FrequencyHours: 12},
	{Company: "Sequoia", AshbySlug: "sequoia", Enabled: true, FrequencyHours: 12},
	{Company: "Help Scout", AshbySlug: "helpscout", Enabled: true, FrequencyHours: 12},
	{Company: "Hopper", AshbySlug: "hopper", Enabled: true, FrequencyHours: 12},
	{Company: "Oyster", AshbySlug: "oyster", Enabled: true, FrequencyHours: 12},
	{Company: "Andela", AshbySlug: "andela", Enabled: true, FrequencyHours: 12},
	{Company: "Leapsome", AshbySlug: "leapsome", Enabled: true, FrequencyHours: 12},
	{Company: "Modern Treasury", AshbySlug: "moderntreasury", Enabled: true, FrequencyHours: 12},
	{Company: "Monte Carlo", AshbySlug: "montecarlodata", Enabled: true, FrequencyHours: 12},
	{Company: "Aurora Solar", AshbySlug: "aurorasolar", Enabled: true, FrequencyHours: 12},
	{Company: "EightSleep", AshbySlug: "eightsleep", Enabled: true, FrequencyHours: 12},
	{Company: "Coder", AshbySlug: "coder", Enabled: true, FrequencyHours: 12},
	{Company: "Form Energy", AshbySlug: "formenergy", Enabled: true, FrequencyHours: 12},
	{Company: "NETGEAR", AshbySlug: "netgear", Enabled: true, FrequencyHours: 12},
	{Company: "Alan", AshbySlug: "alan", Enabled: true, FrequencyHours: 12},
	{Company: "Dave", AshbySlug: "dave", Enabled: true, FrequencyHours: 12},
	{Company: "Marshmallow", AshbySlug: "marshmallow", Enabled: true, FrequencyHours: 12},
	{Company: "GameChanger", AshbySlug: "gamechanger", Enabled: true, FrequencyHours: 12},
	{Company: "CoinTracker", AshbySlug: "cointracker", Enabled: true, FrequencyHours: 12},
	{Company: "Nivoda", AshbySlug: "nivoda", Enabled: true, FrequencyHours: 12},
	{Company: "TheyDo", AshbySlug: "theydo", Enabled: true, FrequencyHours: 12},
	{Company: "Zello", AshbySlug: "zello", Enabled: true, FrequencyHours: 12},
	{Company: "Juniper Square", AshbySlug: "junipersquare", Enabled: true, FrequencyHours: 12},
	{Company: "Rad AI", AshbySlug: "radai", Enabled: true, FrequencyHours: 12},
	{Company: "Rula", AshbySlug: "rula", Enabled: true, FrequencyHours: 12},
	{Company: "Sanity", AshbySlug: "sanity", Enabled: true, FrequencyHours: 12},
	{Company: "Cryptio", AshbySlug: "cryptio", Enabled: true, FrequencyHours: 12},
	{Company: "Decimal", AshbySlug: "decimal", Enabled: true, FrequencyHours: 12},
	{Company: "Shopify", AshbySlug: "shopify", Enabled: false, FrequencyHours: 12},
	{Company: "Plaid", AshbySlug: "plaid", Enabled: false, FrequencyHours: 12},
	{Company: "Clay", AshbySlug: "clay", Enabled: false, FrequencyHours: 12},
	{Company: "Ironclad", AshbySlug: "ironclad", Enabled: false, FrequencyHours: 12},
	{Company: "Lime", AshbySlug: "lime", Enabled: false, FrequencyHours: 12},
	{Company: "Opendoor", AshbySlug: "opendoor", Enabled: false, FrequencyHours: 12},
	{Company: "Duolingo", AshbySlug: "duolingo", Enabled: false, FrequencyHours: 12},
	{Company: "Calm", AshbySlug: "calm", Enabled: false, FrequencyHours: 12},
	{Company: "incident.io", AshbySlug: "incidentio", Enabled: false, FrequencyHours: 12},
	{Company: "Flock Safety", AshbySlug: "flocksafety", Enabled: false, FrequencyHours: 12},
	{Company: "Brightline", AshbySlug: "brightline", Enabled: false, FrequencyHours: 12},
	{Company: "Scribe", AshbySlug: "scribehow", Enabled: false, FrequencyHours: 12},
	{Company: "mParticle", AshbySlug: "mparticle", Enabled: false, FrequencyHours: 12},
	{Company: "Teal", AshbySlug: "teal", Enabled: false, FrequencyHours: 12},
	{Company: "Luma AI", AshbySlug: "lumaai", Enabled: false, FrequencyHours: 12},
	{Company: "WeTransfer", AshbySlug: "wetransfer", Enabled: false, FrequencyHours: 12},
	{Company: "Solana Foundation", AshbySlug: "solanafoundation", Enabled: false, FrequencyHours: 12},
	{Company: "Anthropic", AshbySlug: "anthropic", Enabled: false, FrequencyHours: 12},
	{Company: "Scale AI", AshbySlug: "scaleai", Enabled: false, FrequencyHours: 12},
	{Company: "Brex", AshbySlug: "brex", Enabled: false, FrequencyHours: 12},
	{Company: "Reddit", AshbySlug: "reddit", Enabled: false, FrequencyHours: 12},
	{Company: "Mercury", AshbySlug: "mercury", Enabled: false, FrequencyHours: 12},
	{Company: "Retool", AshbySlug: "retool", Enabled: false, FrequencyHours: 12},
	{Company: "Airtable", AshbySlug: "airtable", Enabled: false, FrequencyHours: 12},
	{Company: "Marqeta", AshbySlug: "marqeta", Enabled: false, FrequencyHours: 12},
	{Company: "Superhuman", AshbySlug: "superhuman", Enabled: false, FrequencyHours: 12},
	{Company: "Boomi", AshbySlug: "boomi", Enabled: false, FrequencyHours: 12},
	{Company: "Nango", AshbySlug: "nango", Enabled: false, FrequencyHours: 12},
	{Company: "Vercel", AshbySlug: "vercel", Enabled: false, FrequencyHours: 12},
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
