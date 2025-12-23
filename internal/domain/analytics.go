package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DashboardStats struct {
	TotalJobs          int                      `json:"total_jobs"`
	TotalApplicants    int                      `json:"total_applicants"`
	StatusDistribution map[string]int           `json:"status_distribution"`
	ApplicationsTrend  []MonthlyStat            `json:"applications_trend"`
	RecentApplicants   []RecentApplicationParam `json:"recent_applicants"`
}

type MonthlyStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type RecentApplicationParam struct {
	ID           uuid.UUID         `json:"id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `json:"deleted_at"`
	JobID        uuid.UUID         `json:"job_id"`
	Job          RecentJobParam    `json:"job"`
	SeekerID     uuid.UUID         `json:"seeker_id"`
	Seeker       RecentSeekerParam `json:"seeker"`
	Status       string            `json:"status"`
	ResumeURL    string            `json:"resume_url"`
	CoverLetter  string            `json:"cover_letter"`
	LinkedInURL  string            `json:"linkedin_url"`
	PortfolioURL string            `json:"portfolio_url"`
}

type RecentJobParam struct {
	ID      uuid.UUID          `json:"id"`
	Title   string             `json:"title"`
	JobType string             `json:"job_type"`
	Company RecentCompanyParam `json:"company"`
}

type RecentCompanyParam struct {
	CompanyName string `json:"company_name"`
}

type RecentSeekerParam struct {
	ID            uuid.UUID           `json:"id"`
	Email         string              `json:"email"`
	SeekerProfile RecentSeekerProfile `json:"seeker_profile"`
}

type RecentSeekerProfile struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	Address  string    `json:"address"`
}
