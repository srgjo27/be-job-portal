package usecase

import (
	"be-job-portal/internal/domain"
	"context"

	"github.com/google/uuid"
)

type applicationUsecase struct {
	appRepo domain.ApplicationRepository
	jobRepo domain.JobRepository
}

func NewApplicationUsecase(appRepo domain.ApplicationRepository, jobRepo domain.JobRepository) domain.ApplicationUsecase {
	return &applicationUsecase{appRepo, jobRepo}
}

func (u *applicationUsecase) ApplyJob(ctx context.Context, jobID, seekerID uuid.UUID, resumeURL, coverLetter, linkedInURL, portfolioURL string) error {
	_, err := u.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return err
	}

	app := &domain.Application{
		JobID:        jobID,
		SeekerID:     seekerID,
		Status:       "PENDING",
		ResumeURL:    resumeURL,
		CoverLetter:  coverLetter,
		LinkedInURL:  linkedInURL,
		PortfolioURL: portfolioURL,
	}
	return u.appRepo.Create(ctx, app)
}

func (u *applicationUsecase) ListApplications(ctx context.Context, userID uuid.UUID, role string) ([]domain.Application, error) {
	if role == "SEEKER" {
		return u.appRepo.GetBySeekerID(ctx, userID)
	}

	return []domain.Application{}, nil
}

func (u *applicationUsecase) ListJobApplicants(ctx context.Context, jobID, recruiterID uuid.UUID) ([]domain.Application, error) {
	job, err := u.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job.RecruiterID != recruiterID {
		return nil, domain.ErrUnauthorized
	}

	return u.appRepo.GetByJobID(ctx, jobID)
}

func (u *applicationUsecase) UpdateStatus(ctx context.Context, appID, recruiterID uuid.UUID, status string) error {
	validStatuses := map[string]bool{
		domain.StatusPending:  true,
		domain.StatusProcess:  true,
		domain.StatusAccepted: true,
		domain.StatusRejected: true,
	}

	if !validStatuses[status] {
		return domain.ErrBadRequest
	}

	app, err := u.appRepo.GetByID(ctx, appID)
	if err != nil {
		return err
	}

	job, err := u.jobRepo.GetByID(ctx, app.JobID)
	if err != nil {
		return err
	}

	if job.RecruiterID != recruiterID {
		return domain.ErrUnauthorized
	}

	return u.appRepo.UpdateStatus(ctx, appID, status)
}

func (u *applicationUsecase) GetDashboardStats(ctx context.Context, recruiterID uuid.UUID) (*domain.DashboardStats, error) {
	return u.appRepo.GetDashboardStats(ctx, recruiterID)
}
