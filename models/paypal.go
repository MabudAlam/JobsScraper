package models

import "encoding/json"

func UnmarshalPayPalJobs(data []byte) (PayPalJobs, error) {
	var r PayPalJobs
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *PayPalJobs) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PayPalJobs struct {
	Domain                         string               `json:"domain"`
	User                           string               `json:"user"`
	IsWillingToRelocate            bool                 `json:"isWillingToRelocate"`
	IsUserAuthenticated            bool                 `json:"isUserAuthenticated"`
	IsUserETXCandidate             bool                 `json:"isUserETXCandidate"`
	IsDomainETX                    bool                 `json:"isDomainETX"`
	IsDomainVeteran                bool                 `json:"isDomainVeteran"`
	SignUpConfig                   Debug                `json:"signUpConfig"`
	GetHelpButton                  Debug                `json:"getHelpButton"`
	IsCareerPlannerEnabled         bool                 `json:"isCareerPlannerEnabled"`
	EnableRememberMeOption         bool                 `json:"enableRememberMeOption"`
	IsMyApplicationsEnabled        bool                 `json:"isMyApplicationsEnabled"`
	ShowVeteranEmployerSignUp      bool                 `json:"showVeteranEmployerSignUp"`
	EnableUserPreferenceSelection  bool                 `json:"enableUserPreferenceSelection"`
	IsUserPreferenceApplied        bool                 `json:"isUserPreferenceApplied"`
	Candidate                      Candidate            `json:"candidate"`
	Branding                       Branding             `json:"branding"`
	Positions                      []Position           `json:"positions"`
	Debug                          Debug                `json:"debug"`
	Count                          int64                `json:"count"`
	CountFilterByMatchScore        interface{}          `json:"countFilterByMatchScore"`
	LocationUser                   []string             `json:"location_user"`
	LocationUsed                   string               `json:"location_used"`
	LocationInsights               interface{}          `json:"locationInsights"`
	ShowWizard                     int64                `json:"showWizard"`
	DisableLocationSearchDropdown  bool                 `json:"disableLocationSearchDropdown"`
	PcsAutocompleteLocationEnabled int64                `json:"pcsAutocompleteLocationEnabled"`
	RecommendedStarThreshold       int64                `json:"recommended_star_threshold"`
	MocTitle                       interface{}          `json:"mocTitle"`
	FuzzyResultsReturned           bool                 `json:"fuzzyResultsReturned"`
	Query                          Query                `json:"query"`
	UserTitles                     []string             `json:"userTitles"`
	IsThinProfile                  bool                 `json:"isThinProfile"`
	VeteranProgramDetails          []interface{}        `json:"veteranProgramDetails"`
	EnableTargetedResume           int64                `json:"enableTargetedResume"`
	JobCardConfig                  interface{}          `json:"jobCardConfig"`
	Facets                         Facets               `json:"facets"`
	LocationRadiusConfig           LocationRadiusConfig `json:"locationRadiusConfig"`
	IsSubQuery                     bool                 `json:"isSubQuery"`
}

type Branding struct {
	PostApplyTitle                          string                  `json:"postApplyTitle"`
	PageImage                               string                  `json:"page_image"`
	EnableTalentNetwork                     int64                   `json:"enableTalentNetwork"`
	Favicons                                Favicons                `json:"favicons"`
	TalentNetworkHeroBanner                 TalentNetworkHeroBanner `json:"talentNetworkHeroBanner"`
	TalentFormMandatoryFields               []string                `json:"talent_form_mandatory_fields"`
	ShowJobID                               int64                   `json:"showJobId"`
	I18NOverridesMaster                     I18NOverridesMaster     `json:"i18n_overrides_master"`
	CustomContent                           BrandingCustomContent   `json:"customContent"`
	TalentNetworkBranding                   TalentNetworkBranding   `json:"talentNetworkBranding"`
	HomePageHeroBanner                      HomePageHeroBanner      `json:"homePageHeroBanner"`
	CompanyName                             string                  `json:"companyName"`
	Perks                                   []Perk                  `json:"perks"`
	CustomStyle                             BrandingCustomStyle     `json:"custom_style"`
	HideJobCart                             bool                    `json:"hideJobCart"`
	Privacy                                 Privacy                 `json:"privacy"`
	CustomHTML                              BrandingCustomHTML      `json:"custom_html"`
	CustomHeadScripts                       CustomHeadScripts       `json:"custom_head_scripts"`
	NavBar                                  NavBar                  `json:"navBar"`
	PageTitle                               string                  `json:"page_title"`
	JobPageTitle                            string                  `json:"job_page_title"`
	PageDescription                         string                  `json:"page_description"`
	DefaultState                            DefaultState            `json:"defaultState"`
	RecaptchaEnabled                        int64                   `json:"recaptcha_enabled"`
	ShowLoggedOutNotificationsPrivacyPolicy bool                    `json:"showLoggedOutNotificationsPrivacyPolicy"`
}

type BrandingCustomContent struct {
	PositionSections []PositionSection `json:"positionSections"`
}

type PositionSection struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}

type BrandingCustomHTML struct {
	Footer string `json:"footer"`
	Header string `json:"header"`
}

type CustomHeadScripts struct {
	HomePageLoad              string `json:"HOME_PAGE_LOAD"`
	SinglePositionPageLoad    string `json:"SINGLE_POSITION_PAGE_LOAD"`
	ApplyFormPageLoad         string `json:"APPLY_FORM_PAGE_LOAD"`
	SuccessFormPageLoad       string `json:"SUCCESS_FORM_PAGE_LOAD"`
	JoinTalentNetworkPageLoad string `json:"JOIN_TALENT_NETWORK_PAGE_LOAD"`
}

type BrandingCustomStyle struct {
	CSS string `json:"css"`
}

type DefaultState struct {
	Pymww bool `json:"pymww"`
}

type Favicons struct {
	Favicon string `json:"favicon"`
}

type HomePageHeroBanner struct {
	UseImage         int64   `json:"useImage"`
	Opacity          float64 `json:"opacity"`
	HideInMobileView bool    `json:"hideInMobileView"`
	Image            string  `json:"image"`
	Title            string  `json:"title"`
}

type I18NOverridesMaster struct {
	CustomHTML            TalentNetworkBrandingClass       `json:"custom_html"`
	CustomContent         I18NOverridesMasterCustomContent `json:"customContent"`
	TalentNetworkBranding TalentNetworkBrandingClass       `json:"talentNetworkBranding"`
}

type I18NOverridesMasterCustomContent struct {
	En CustomContentEn `json:"en"`
}

type CustomContentEn struct {
	K1 string `json:"k1"`
	K2 string `json:"k2"`
	K3 string `json:"k3"`
	K4 string `json:"k4"`
	ID string `json:"ID"`
}

type TalentNetworkBrandingClass struct {
	En CustomHTMLEn `json:"en"`
}

type CustomHTMLEn struct {
	JoinTalentNetwork string `json:"Join Talent Network"`
	ID                string `json:"ID"`
}

type NavBar struct {
	Opacity    int64  `json:"opacity"`
	Background string `json:"background"`
	Image      string `json:"image"`
}

type Perk struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type Privacy struct {
	Text   string `json:"text"`
	Button string `json:"button"`
	Title  string `json:"title"`
}

type TalentNetworkBranding struct {
	CustomHTML  BrandingCustomHTML               `json:"custom_html"`
	CustomStyle TalentNetworkBrandingCustomStyle `json:"custom_style"`
}

type TalentNetworkBrandingCustomStyle struct {
	Font string `json:"font"`
	CSS  string `json:"css"`
}

type TalentNetworkHeroBanner struct {
	UseImage int64  `json:"useImage"`
	Opacity  int64  `json:"opacity"`
	Image    string `json:"image"`
}

type Candidate struct {
	EncID               string        `json:"enc_id"`
	Fullname            string        `json:"fullname"`
	Firstname           string        `json:"firstname"`
	Lastname            string        `json:"lastname"`
	Skills              []string      `json:"skills"`
	Email               string        `json:"email"`
	Phone               string        `json:"phone"`
	Location            string        `json:"location"`
	Filename            string        `json:"filename"`
	StarredPositions    []interface{} `json:"starred_positions"`
	ResumeURL           string        `json:"resumeUrl"`
	OnboardingCompleted bool          `json:"onboardingCompleted"`
	IsUserInPcsIjp      bool          `json:"isUserInPcsIjp"`
	LinkedinURL         string        `json:"linkedinUrl"`
}

type Debug struct {
}

type PayPalJobData struct {
	Company        Debug            `json:"Company"`
	Skills         map[string]int64 `json:"Skills"`
	JobCategory    map[string]int64 `json:"Job Category"`
	Country        map[string]int64 `json:"Country"`
	Brand          Brand            `json:"Brand"`
	EmploymentType EmploymentType   `json:"Employment Type"`
	TypeOfWork     TypeOfWork       `json:"Type of Work"`
	Locations      map[string]int64 `json:"locations"`
}

type Brand struct {
	PayPal       int64 `json:"PayPal"`
	Venmo        int64 `json:"venmo"`
	Zettle       int64 `json:"zettle"`
	HappyReturns int64 `json:"happy returns"`
	Xoom         int64 `json:"xoom"`
	Braintree    int64 `json:"braintree"`
	Honey        int64 `json:"honey"`
}

type EmploymentType struct {
	FullTime int64 `json:"full time"`
}

type TypeOfWork struct {
	InPerson int64 `json:"in person"`
	Remote   int64 `json:"remote"`
}

type LocationRadiusConfig struct {
	ShowLocationRadius bool   `json:"showLocationRadius"`
	LocationRadiusType string `json:"locationRadiusType"`
}

type Position struct {
	ID                   int64              `json:"id"`
	Name                 string             `json:"name"`
	Location             Location           `json:"location"`
	Locations            []Location         `json:"locations"`
	Hot                  int64              `json:"hot"`
	Department           Department         `json:"department"`
	BusinessUnit         BusinessUnit       `json:"business_unit"`
	TUpdate              int64              `json:"t_update"`
	TCreate              int64              `json:"t_create"`
	AtsJobID             string             `json:"ats_job_id"`
	DisplayJobID         string             `json:"display_job_id"`
	Type                 Type               `json:"type"`
	IDLocale             string             `json:"id_locale"`
	JobDescription       string             `json:"job_description"`
	Locale               Locale             `json:"locale"`
	Stars                float64            `json:"stars"`
	MedallionProgram     interface{}        `json:"medallionProgram"`
	LocationFlexibility  interface{}        `json:"location_flexibility"`
	WorkLocationOption   WorkLocationOption `json:"work_location_option"`
	CanonicalPositionURL string             `json:"canonicalPositionUrl"`
	MatchDetails         []MatchDetail      `json:"match_details,omitempty"`
	MatchSkills          []string           `json:"match_skills,omitempty"`
	Insights             *Insights          `json:"insights,omitempty"`
}

type Insights struct {
	Skills          Debug `json:"skills"`
	Companies       Debug `json:"companies"`
	Schools         Debug `json:"schools"`
	Degrees         Debug `json:"degrees"`
	ExperienceYears Debug `json:"experience_years"`
	Seniority       Debug `json:"seniority"`
	Titles          Debug `json:"titles"`
	NumVeterans     int64 `json:"num_veterans"`
}

type MatchDetail struct {
	Key  string `json:"key"`
	Text string `json:"text"`
	Desc string `json:"desc"`
}

type Query struct {
	Query          interface{}   `json:"query"`
	Location       string        `json:"location"`
	Department     []interface{} `json:"department"`
	Skill          []interface{} `json:"skill"`
	Seniority      []interface{} `json:"seniority"`
	PID            string        `json:"pid"`
	Company        []interface{} `json:"Company"`
	Skills         []interface{} `json:"Skills"`
	JobCategory    []interface{} `json:"Job Category"`
	Country        []string      `json:"Country"`
	Brand          []interface{} `json:"Brand"`
	EmploymentType []interface{} `json:"Employment Type"`
	TypeOfWork     []interface{} `json:"Type of Work"`
}

type BusinessUnit string

const (
	Unknown BusinessUnit = "Unknown"
)

type Department string

const (
	SoftwareDevelopment        Department = "Software Development"
	TechnicalProductManagement Department = "Technical Product Management"
)

type Locale string

const (
	En Locale = "en"
)

type Location string

const (
	Bangalore               Location = "Bangalore, Karnataka, India"
	ChennaiTamilNaduIndia   Location = "Chennai, Tamil Nadu, India"
	MaharashtraIndiaVirtual Location = "Maharashtra, India (Virtual)"
)

type Type string

const (
	Ats Type = "ATS"
)

type WorkLocationOption string

const (
	OnSite WorkLocationOption = "On-site"
	Remote WorkLocationOption = "Remote"
)
