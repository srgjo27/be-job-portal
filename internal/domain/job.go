package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Title       string         `json:"title" binding:"required"`
	Description string         `gorm:"type:text" json:"description" binding:"required"`
	Category    string         `json:"category"`
	JobType     string         `json:"job_type"`
	Salary      string         `json:"salary"`
	Benefits    []string       `gorm:"serializer:json" json:"benefits"`
	RecruiterID uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"recruiter_id"`
	Recruiter   *User          `gorm:"foreignKey:RecruiterID;references:ID" json:"-"`
	Company     JobCompany     `gorm:"foreignKey:RecruiterID;references:UserID" json:"company"`
}

type JobCompany struct {
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid" json:"-"`
	CompanyName string    `json:"company_name"`
	Location    string    `json:"location"`
	LogoURL     string    `json:"logo_url"`
}

func (JobCompany) TableName() string {
	return "company_profiles"
}

type PaginationParams struct {
	Page  int
	Limit int
}

type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

type PaginatedJobsResponse struct {
	Jobs       []Job          `json:"jobs"`
	Pagination PaginationMeta `json:"pagination"`
}

type JobRepository interface {
	Create(ctx context.Context, job *Job) error
	Update(ctx context.Context, job *Job) error
	GetAll(ctx context.Context, params PaginationParams) ([]Job, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)
	GetByRecruiterID(ctx context.Context, recruiterID uuid.UUID) ([]Job, error)
}

type JobUsecase interface {
	CreateJob(ctx context.Context, title, description, category, jobType, salary string, benefits []string, recruiterID uuid.UUID) error
	UpdateJob(ctx context.Context, id, recruiterID uuid.UUID, title, description, category, jobType, salary string, benefits []string) error
	ListJobs(ctx context.Context, params PaginationParams) (*PaginatedJobsResponse, error)
	GetJob(ctx context.Context, id uuid.UUID) (*Job, error)
	ListJobsByRecruiter(ctx context.Context, recruiterID uuid.UUID) ([]Job, error)
}
