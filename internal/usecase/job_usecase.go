package usecase

import (
	"be-job-portal/internal/domain"
	"context"

	"github.com/google/uuid"
)

type jobUsecase struct {
	jobRepo domain.JobRepository
}

func NewJobUsecase(jobRepo domain.JobRepository) domain.JobUsecase {
	return &jobUsecase{jobRepo}
}

func (u *jobUsecase) CreateJob(ctx context.Context, title, description, category, jobType, salary string, benefits []string, recruiterID uuid.UUID) error {
	job := &domain.Job{
		Title:       title,
		Description: description,
		Category:    category,
		JobType:     jobType,
		Salary:      salary,
		Benefits:    benefits,
		RecruiterID: recruiterID,
	}
	return u.jobRepo.Create(ctx, job)
}

func (u *jobUsecase) UpdateJob(ctx context.Context, id, recruiterID uuid.UUID, title, description, category, jobType, salary string, benefits []string) error {
	job, err := u.jobRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if job == nil {
		return domain.ErrNotFound
	}

	if job.RecruiterID != recruiterID {
		return domain.ErrUnauthorized
	}

	job.Title = title
	job.Description = description
	job.Category = category
	job.JobType = jobType
	job.Salary = salary
	job.Benefits = benefits

	return u.jobRepo.Update(ctx, job)
}

func (u *jobUsecase) ListJobs(ctx context.Context, params domain.PaginationParams) (*domain.PaginatedJobsResponse, error) {
	jobs, totalCount, err := u.jobRepo.GetAll(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := int(totalCount) / params.Limit
	if int(totalCount)%params.Limit != 0 {
		totalPages++
	}

	paginationMeta := domain.PaginationMeta{
		CurrentPage:  params.Page,
		TotalPages:   totalPages,
		TotalItems:   totalCount,
		ItemsPerPage: params.Limit,
		HasNext:      params.Page < totalPages,
		HasPrev:      params.Page > 1,
	}

	return &domain.PaginatedJobsResponse{
		Jobs:       jobs,
		Pagination: paginationMeta,
	}, nil
}

func (u *jobUsecase) GetJob(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	return u.jobRepo.GetByID(ctx, id)
}

func (u *jobUsecase) ListJobsByRecruiter(ctx context.Context, recruiterID uuid.UUID) ([]domain.Job, error) {
	return u.jobRepo.GetByRecruiterID(ctx, recruiterID)
}
