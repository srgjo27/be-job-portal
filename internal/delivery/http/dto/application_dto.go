package dto

import (
	"time"

	"github.com/google/uuid"
)

type ApplyJobRequest struct {
	JobID        string `json:"job_id" binding:"required"`
	ResumeURL    string `json:"resume_url" binding:"required"`
	CoverLetter  string `json:"cover_letter"`
	LinkedInURL  string `json:"linkedin_url"`
	PortfolioURL string `json:"portfolio_url"`
}

type UpdateApplicationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type JobLiteResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	CompanyName string    `json:"company_name"`
	Category    string    `json:"category"`
	JobType     string    `json:"job_type"`
}

type ApplicationListResponse struct {
	ID           uuid.UUID       `json:"id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at"`
	JobID        uuid.UUID       `json:"job_id"`
	Job          JobLiteResponse `json:"job"`
	SeekerID     uuid.UUID       `json:"seeker_id"`
	Status       string          `json:"status"`
	ResumeURL    string          `json:"resume_url"`
	CoverLetter  string          `json:"cover_letter"`
	LinkedInURL  string          `json:"linkedin_url"`
	PortfolioURL string          `json:"portfolio_url"`
}
