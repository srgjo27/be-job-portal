package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
	Email          string          `gorm:"uniqueIndex;not null" json:"email"`
	Password       string          `json:"-"`
	Role           string          `gorm:"default:'SEEKER'" json:"role"`    // RECRUITER, SEEKER
	Provider       string          `gorm:"default:'local'" json:"provider"` // local, google
	SeekerProfile  *SeekerProfile  `gorm:"foreignKey:UserID" json:"seeker_profile,omitempty"`
	CompanyProfile *CompanyProfile `gorm:"foreignKey:UserID" json:"company_profile,omitempty"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type AuthUsecase interface {
	Register(ctx context.Context, email, password, role string) error
	Login(ctx context.Context, email, password string) (string, error)
	GoogleLogin(ctx context.Context, code string) (string, error)
}
