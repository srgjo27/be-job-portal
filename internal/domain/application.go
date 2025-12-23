package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Application struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	JobID        uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"job_id" binding:"required"`
	Job          *Job           `gorm:"foreignKey:JobID;references:ID" json:"job,omitempty"`
	SeekerID     uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"seeker_id"`
	Seeker       *User          `gorm:"foreignKey:SeekerID;references:ID" json:"seeker,omitempty"`
	Status       string         `gorm:"default:'PENDING'" json:"status"` // PENDING, PROCESS, ACCEPTED, REJECTED
	ResumeURL    string         `json:"resume_url" binding:"required"`
	CoverLetter  string         `gorm:"type:text" json:"cover_letter"`
	LinkedInURL  string         `json:"linkedin_url"`
	PortfolioURL string         `json:"portfolio_url"`
}

const (
	StatusPending  = "PENDING"
	StatusProcess  = "PROCESS"
	StatusAccepted = "ACCEPTED"
	StatusRejected = "REJECTED"
)

type ApplicationRepository interface {
	Create(ctx context.Context, app *Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*Application, error)
	GetByJobID(ctx context.Context, jobID uuid.UUID) ([]Application, error)
	GetBySeekerID(ctx context.Context, seekerID uuid.UUID) ([]Application, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	GetDashboardStats(ctx context.Context, recruiterID uuid.UUID) (*DashboardStats, error)
}

type ApplicationUsecase interface {
	ApplyJob(ctx context.Context, jobID, seekerID uuid.UUID, resumeURL, coverLetter, linkedInURL, portfolioURL string) error
	ListApplications(ctx context.Context, userID uuid.UUID, role string) ([]Application, error)
	ListJobApplicants(ctx context.Context, jobID, recruiterID uuid.UUID) ([]Application, error)
	UpdateStatus(ctx context.Context, appID, recruiterID uuid.UUID, status string) error
	GetDashboardStats(ctx context.Context, recruiterID uuid.UUID) (*DashboardStats, error)
}
