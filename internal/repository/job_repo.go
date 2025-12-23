package repository

import (
	"be-job-portal/internal/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) domain.JobRepository {
	return &jobRepository{db}
}

func (r *jobRepository) Create(ctx context.Context, job *domain.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) Update(ctx context.Context, job *domain.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *jobRepository) GetAll(ctx context.Context, params domain.PaginationParams) ([]domain.Job, int64, error) {
	var jobs []domain.Job
	var totalCount int64

	if err := r.db.WithContext(ctx).Model(&domain.Job{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := r.db.WithContext(ctx).
		Preload("Company").
		Order("created_at DESC").
		Limit(params.Limit).
		Offset(offset).
		Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, totalCount, nil
}

func (r *jobRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	var job domain.Job
	if err := r.db.WithContext(ctx).Preload("Company").First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetByRecruiterID(ctx context.Context, recruiterID uuid.UUID) ([]domain.Job, error) {
	var jobs []domain.Job
	if err := r.db.WithContext(ctx).Preload("Company").Where("recruiter_id = ?", recruiterID).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}
