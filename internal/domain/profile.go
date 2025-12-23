package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Experience struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	SeekerProfileID uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"seeker_profile_id"`
	SeekerProfile   *SeekerProfile `gorm:"foreignKey:SeekerProfileID;references:ID" json:"-"`
	Title           string         `json:"title"`
	Company         string         `json:"company"`
	Location        string         `json:"location"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
	Description     string         `gorm:"type:text" json:"description"`
}

type Education struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	SeekerProfileID uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"seeker_profile_id"`
	SeekerProfile   *SeekerProfile `gorm:"foreignKey:SeekerProfileID;references:ID" json:"-"`
	Institution     string         `json:"institution"`
	Degree          string         `json:"degree"`
	FieldOfStudy    string         `json:"field_of_study"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
}

type SeekerProfile struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID;references:ID" json:"-"`
	FullName     string         `json:"full_name"`
	Phone        string         `json:"phone"`
	Address      string         `json:"address"`
	ResumeURL    string         `json:"resume_url"`
	PortfolioURL string         `json:"portfolio_url"`
	LinkedInURL  string         `json:"linkedin_url"`
	Description  string         `json:"description"`
	Skills       []string       `gorm:"serializer:json" json:"skills"`
	Experiences  []Experience   `gorm:"foreignKey:SeekerProfileID" json:"experiences"`
	Educations   []Education    `gorm:"foreignKey:SeekerProfileID" json:"educations"`
}

type CompanyProfile struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID;references:ID" json:"-"`
	CompanyName string         `json:"company_name"`
	Website     string         `json:"website"`
	Phone       string         `json:"phone"`
	Location    string         `json:"location"`
	Description string         `gorm:"type:text" json:"description"`
	LogoURL     string         `json:"logo_url"`
}

type ProfileRepository interface {
	GetSeekerProfile(ctx context.Context, userID uuid.UUID) (*SeekerProfile, error)
	UpdateSeekerProfile(ctx context.Context, profile *SeekerProfile) error
	GetCompanyProfile(ctx context.Context, userID uuid.UUID) (*CompanyProfile, error)
	UpdateCompanyProfile(ctx context.Context, profile *CompanyProfile) error
}

type ProfileUsecase interface {
	GetProfile(ctx context.Context, userID uuid.UUID, role string) (interface{}, error)
	UpdateSeekerProfile(ctx context.Context, userID uuid.UUID, profile *SeekerProfile) error
	UpdateCompanyProfile(ctx context.Context, userID uuid.UUID, profile *CompanyProfile) error
}
