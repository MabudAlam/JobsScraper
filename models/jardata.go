package models

type JarAPIResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []JarJob `json:"results"`
}

type JarJob struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	JobType        string   `json:"job_type"`
	Description    string   `json:"description"`
	MaxSalary      float64  `json:"max_salary"`
	MinSalary      float64  `json:"min_salary"`
	MaxExperience  float64  `json:"max_experience"`
	Skill          []string `json:"skill"`
	Location       string   `json:"location"`
	DepartmentName string   `json:"department_name"`
	CompanyName    string   `json:"company_name"`
	WorkplaceType  string   `json:"workplace_type"`
	Remote         bool     `json:"remote"`
	CreatedAt      string   `json:"created_at"`
	SalaryType     string   `json:"salary_type"`
	ApplyLink      string   `json:"apply_link"`
}
