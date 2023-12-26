package models

import "encoding/json"

func UnmarshalTemperatures(data []byte) (Temperatures, error) {
	var r Temperatures
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Temperatures) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Temperatures struct {
	HasNext  bool        `json:"hasNext"`
	Next     interface{} `json:"next"`
	Postings []Posting   `json:"postings"`
}

type Posting struct {
	Id               int      `json:"id"`
	PortalId         int      `json:"portalId"`
	Title            string   `json:"title"`
	Type             string   `json:"type"`
	Location         []string `json:"locations"`
	Category         string   `json:"category"`
	Overview         string   `json:"overview"`
	Responsibilities string   `json:"responsibilities"`
	Qualifications   string   `json:"qualifications"`
	ApplyUrl         string   `json:"applyUrl"`
	LeverId          string   `json:"leverId"`
}

type Categories struct {
	Commitment Commitment  `json:"commitment"`
	Department interface{} `json:"department"`
	Level      interface{} `json:"level"`
	Location   string      `json:"location"`
	Team       Team        `json:"team"`
}

type Content struct {
	Description     string      `json:"description"`
	DescriptionHTML string      `json:"descriptionHtml"`
	Lists           []List      `json:"lists"`
	Closing         string      `json:"closing"`
	ClosingHTML     string      `json:"closingHtml"`
	CustomQuestions interface{} `json:"customQuestions"`
}

type List struct {
	Text    string `json:"text"`
	Content string `json:"content"`
}

type Requisition struct {
	ID               string           `json:"id"`
	CreatedAt        int64            `json:"createdAt"`
	Status           Status           `json:"status"`
	Name             string           `json:"name"`
	RequisitionCode  string           `json:"requisitionCode"`
	HeadcountTotal   int64            `json:"headcountTotal"`
	HeadcountHired   int64            `json:"headcountHired"`
	Backfill         bool             `json:"backfill"`
	Creator          interface{}      `json:"creator"`
	Owner            string           `json:"owner"`
	HiringManager    string           `json:"hiringManager"`
	InternalNotes    string           `json:"internalNotes"`
	CompensationBand interface{}      `json:"compensationBand"`
	EmploymentStatus EmploymentStatus `json:"employmentStatus"`
	Location         string           `json:"location"`
	Team             string           `json:"team"`
	Postings         []string         `json:"postings"`
	OfferIDS         []string         `json:"offerIds"`
	CustomFields     CustomFields     `json:"customFields"`
}

type CustomFields struct {
	Creationreason          Creationreason     `json:"creationreason"`
	Budgetedquarter         Budgetedquarter    `json:"budgetedquarter"`
	Workertype              Workertype         `json:"workertype"`
	Managementlevel         string             `json:"managementlevel"`
	Compensationgrade       Compensationgrade  `json:"compensationgrade"`
	Fbogm                   Fbogm              `json:"fbogm"`
	Managerid               int64              `json:"managerid"`
	Jobcode                 int64              `json:"jobcode"`
	Jobprofile              string             `json:"jobprofile"`
	Jobfamily               string             `json:"jobfamily"`
	Jobclassification       Jobclassification  `json:"jobclassification"`
	Jobexempt               Jobexempt          `json:"jobexempt"`
	Companyinsider          CompanyinsiderEnum `json:"companyinsider"`
	Product                 string             `json:"product"`
	Group                   string             `json:"group"`
	Subgroup                string             `json:"subgroup"`
	Team                    string             `json:"team"`
	Recruitingstartdate     string             `json:"recruitingstartdate"`
	Targethiredate          string             `json:"targethiredate"`
	Visarelo                Jobexempt          `json:"visarelo"`
	Lumpsum                 Jobexempt          `json:"lumpsum"`
	Company                 Company            `json:"company"`
	Originalbudgetedquarter CompanyinsiderEnum `json:"originalbudgetedquarter"`
	Costcenter              string             `json:"costcenter"`
	Jobfamilygroup          Jobfamilygroup     `json:"jobfamilygroup"`
}

type Urls struct {
	ListURL   string `json:"listUrl"`
	ShowURL   string `json:"showUrl"`
	ApplyURL  string `json:"applyUrl"`
	HostedURL string `json:"hostedUrl"`
}

type Commitment string

const (
	FullTime Commitment = "Full Time"
)

type Team string

const (
	AnalyticsDataScience Team = "Analytics & Data Science"
	CustomerExperience   Team = "Customer Experience"
	Engineering          Team = "Engineering"
	FinanceAccounting    Team = "Finance & Accounting"
	ProductManagement    Team = "Product Management"
	ProgramManagement    Team = "Program Management"
	Security             Team = "Security"
	Sre                  Team = "SRE"
	SupportCSS           Team = "Support (CSS)"
	TeamIT               Team = "IT"
	TeamLegal            Team = "Legal"
	TeamMarketing        Team = "Marketing"
	TeamPeople           Team = "People"
	TeamSales            Team = "Sales"
)

type ReqCode string

const (
	Empty    ReqCode = ""
	Ref2054N ReqCode = "REF2054N"
)

type Budgetedquarter string

const (
	BudgetedquarterFY22Q1 Budgetedquarter = "FY22 Q1"
	BudgetedquarterFY23Q4 Budgetedquarter = "FY23 Q4"
	Fy19Q4                Budgetedquarter = "FY19 Q4"
	Fy20Q3                Budgetedquarter = "FY20 Q3"
	Fy24Q1                Budgetedquarter = "FY24 Q1"
	Fy24Q2                Budgetedquarter = "FY24 Q2"
	Fy24Q3                Budgetedquarter = "FY24 Q3"
	Fy24Q4                Budgetedquarter = "FY24 Q4"
)

type Company string

const (
	AtlassianBV             Company = "Atlassian B.V."
	AtlassianFrance         Company = "Atlassian France"
	AtlassianGermanyGmbH    Company = "Atlassian Germany GmbH"
	AtlassianIndiaLLP       Company = "Atlassian India LLP"
	AtlassianKK             Company = "Atlassian K.K."
	AtlassianPhilippinesInc Company = "Atlassian Philippines, Inc."
	AtlassianPolandSPZOO    Company = "Atlassian Poland sp. z o.o."
	AtlassianPtyLtd         Company = "Atlassian Pty Ltd"
	AtlassianUSInc          Company = "Atlassian US, Inc."
	OpsGenieAŞ              Company = "OpsGenie A.Ş."
)

type CompanyinsiderEnum string

const (
	Companyinsider       CompanyinsiderEnum = " "
	CompanyinsiderFY22Q1 CompanyinsiderEnum = "FY22 Q1"
	CompanyinsiderFY23Q4 CompanyinsiderEnum = "FY23 Q4"
	Fy22Q4               CompanyinsiderEnum = "FY22 Q4"
	Fy23Q2               CompanyinsiderEnum = "FY23 Q2"
)

type Compensationgrade string

const (
	M40 Compensationgrade = "M40"
	M50 Compensationgrade = "M50"
	M60 Compensationgrade = "M60"
	M70 Compensationgrade = "M70"
	M80 Compensationgrade = "M80"
	P20 Compensationgrade = "P20"
	P30 Compensationgrade = "P30"
	P40 Compensationgrade = "P40"
	P50 Compensationgrade = "P50"
	P60 Compensationgrade = "P60"
	P70 Compensationgrade = "P70"
	P80 Compensationgrade = "P80"
	P90 Compensationgrade = "P90"
)

type Creationreason string

const (
	AnnualPlan       Creationreason = "Annual Plan"
	Attrition        Creationreason = "Attrition"
	FY22Q4SoA        Creationreason = "FY22 Q4 SoA"
	FY23Q1SoA        Creationreason = "FY23 Q1 SoA"
	FY23Q2SoA        Creationreason = "FY23 Q2 SoA"
	FY23Q3SoA        Creationreason = "FY23 Q3 SoA"
	InternalTransfer Creationreason = "Internal Transfer"
	MidYearApproval  Creationreason = "Mid Year Approval"
)

type Fbogm string

const (
	AdityaPhadke   Fbogm = "Aditya Phadke"
	AnuBharadwaj   Fbogm = "Anu Bharadwaj"
	BryanMayo      Fbogm = "Bryan Mayo"
	CameronDeatsch Fbogm = "Cameron Deatsch"
	ErikaFisher    Fbogm = "Erika Fisher"
	JoeBinz        Fbogm = "Joe Binz"
	KevinEgan      Fbogm = "Kevin Egan"
	RajeevRajan    Fbogm = "Rajeev Rajan"
	TalSaraf       Fbogm = "Tal Saraf"
	TaroonMandhana Fbogm = "Taroon Mandhana"
)

type Jobclassification string

const (
	The11ExecutiveSeniorLevelOfficialsAndManagers Jobclassification = "1.1 - Executive/Senior Level Officials and Managers"
	The12FirstMidLevelOfficialsAndManagers        Jobclassification = "1.2 - First/Mid-Level Officials and Managers"
	The2Professionals                             Jobclassification = "2 - Professionals"
	The4SalesWorkers                              Jobclassification = "4 - Sales Workers"
	The5AdministrativeSupportWorkers              Jobclassification = "5 - Administrative Support Workers"
)

type Jobexempt string

const (
	No  Jobexempt = "No"
	Yes Jobexempt = "Yes"
)

type Jobfamilygroup string

const (
	CustomerExperienceCX              Jobfamilygroup = "Customer Experience (CX)"
	DataAndAnalytics                  Jobfamilygroup = "Data and Analytics"
	EngineeringJFG                    Jobfamilygroup = "Engineering JFG"
	EnterpriseBusinessServices        Jobfamilygroup = "Enterprise Business Services"
	Finance                           Jobfamilygroup = "Finance"
	JobfamilygroupIT                  Jobfamilygroup = "IT"
	JobfamilygroupLegal               Jobfamilygroup = "Legal"
	JobfamilygroupMarketing           Jobfamilygroup = "Marketing"
	JobfamilygroupPeople              Jobfamilygroup = "People"
	JobfamilygroupSales               Jobfamilygroup = "Sales"
	Product                           Jobfamilygroup = "Product"
	StrategyPlanningAndTransformation Jobfamilygroup = "Strategy Planning and Transformation"
	Support                           Jobfamilygroup = "Support"
)

type Workertype string

const (
	Regular Workertype = "Regular"
)

type EmploymentStatus string

const (
	EmploymentStatusFullTime EmploymentStatus = "Full time"
)

type Status string

const (
	Open Status = "open"
)

type State string

const (
	Published State = "published"
)
