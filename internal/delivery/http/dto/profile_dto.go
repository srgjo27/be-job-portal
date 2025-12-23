package dto

type ExperienceInput struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	StartDate   string `json:"start_date"` // Format: 2006-01-02
	EndDate     string `json:"end_date"`   // Format: 2006-01-02
	Description string `json:"description"`
}

type EducationInput struct {
	Institution  string `json:"institution"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study"`
	StartDate    string `json:"start_date"` // Format: 2006-01-02
	EndDate      string `json:"end_date"`   // Format: 2006-01-02
}

type UpdateSeekerProfileRequest struct {
	FullName     string            `json:"full_name"`
	Phone        string            `json:"phone"`
	Address      string            `json:"address"`
	ResumeURL    string            `json:"resume_url"`
	PortfolioURL string            `json:"portfolio_url"`
	LinkedInURL  string            `json:"linkedin_url"`
	Description  string            `json:"description"`
	Skills       []string          `json:"skills"`
	Experiences  []ExperienceInput `json:"experiences"`
	Educations   []EducationInput  `json:"educations"`
}

type UpdateCompanyProfileRequest struct {
	CompanyName string `json:"company_name"`
	Website     string `json:"website"`
	Phone       string `json:"phone"`
	Location    string `json:"location"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}
