package repository

import (
	"be-job-portal/internal/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) domain.ApplicationRepository {
	return &applicationRepository{db}
}

func (r *applicationRepository) Create(ctx context.Context, app *domain.Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *applicationRepository) GetByJobID(ctx context.Context, jobID uuid.UUID) ([]domain.Application, error) {
	var apps []domain.Application
	err := r.db.WithContext(ctx).Preload("Seeker").Preload("Seeker.SeekerProfile").Where("job_id = ?", jobID).Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) GetBySeekerID(ctx context.Context, seekerID uuid.UUID) ([]domain.Application, error) {
	var apps []domain.Application
	err := r.db.WithContext(ctx).
		Preload("Job", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, title, category, job_type, recruiter_id")
		}).
		Preload("Job.Company").
		Where("seeker_id = ?", seekerID).
		Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	var app domain.Application
	err := r.db.WithContext(ctx).First(&app, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Application{}).Where("id = ?", id).Update("status", status).Error
}

func (r *applicationRepository) GetDashboardStats(ctx context.Context, recruiterID uuid.UUID) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{
		StatusDistribution: make(map[string]int),
	}

	var totalJobs int64
	if err := r.db.Model(&domain.Job{}).Where("recruiter_id = ?", recruiterID).Count(&totalJobs).Error; err != nil {
		return nil, err
	}
	stats.TotalJobs = int(totalJobs)

	var totalApplicants int64
	err := r.db.Model(&domain.Application{}).
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Where("jobs.recruiter_id = ?", recruiterID).
		Count(&totalApplicants).Error
	if err != nil {
		return nil, err
	}
	stats.TotalApplicants = int(totalApplicants)

	rows, err := r.db.Model(&domain.Application{}).
		Select("applications.status, count(*) as count").
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Where("jobs.recruiter_id = ?", recruiterID).
		Group("applications.status").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		stats.StatusDistribution[status] = count
	}

	trendRows, err := r.db.Model(&domain.Application{}).
		Select("to_char(applications.created_at, 'YYYY-MM-DD') as date, count(*) as count").
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Where("jobs.recruiter_id = ?", recruiterID).
		Group("date").
		Order("date ASC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer trendRows.Close()

	for trendRows.Next() {
		var date string
		var count int
		if err := trendRows.Scan(&date, &count); err != nil {
			continue
		}
		stats.ApplicationsTrend = append(stats.ApplicationsTrend, domain.MonthlyStat{Date: date, Count: count})
	}

	var recentApps []domain.Application
	err = r.db.Preload("Job", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, title, job_type, recruiter_id")
	}).Preload("Job.Company").
		Preload("Seeker").Preload("Seeker.SeekerProfile").
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Where("jobs.recruiter_id = ?", recruiterID).
		Order("applications.created_at DESC").
		Limit(5).
		Find(&recentApps).Error
	if err != nil {
		return nil, err
	}

	var recentAppsParams []domain.RecentApplicationParam
	for _, app := range recentApps {
		param := domain.RecentApplicationParam{
			ID:           app.ID,
			CreatedAt:    app.CreatedAt,
			UpdatedAt:    app.UpdatedAt,
			DeletedAt:    app.DeletedAt,
			JobID:        app.JobID,
			SeekerID:     app.SeekerID,
			Status:       app.Status,
			ResumeURL:    app.ResumeURL,
			CoverLetter:  app.CoverLetter,
			LinkedInURL:  app.LinkedInURL,
			PortfolioURL: app.PortfolioURL,
		}

		if app.Job != nil {
			param.Job = domain.RecentJobParam{
				ID:      app.Job.ID,
				Title:   app.Job.Title,
				JobType: app.Job.JobType,
				Company: domain.RecentCompanyParam{
					CompanyName: app.Job.Company.CompanyName,
				},
			}
		}

		if app.Seeker != nil {
			param.Seeker = domain.RecentSeekerParam{
				ID:    app.Seeker.ID,
				Email: app.Seeker.Email,
			}
			if app.Seeker.SeekerProfile != nil {
				param.Seeker.SeekerProfile = domain.RecentSeekerProfile{
					ID:       app.Seeker.SeekerProfile.ID,
					UserID:   app.Seeker.SeekerProfile.UserID,
					FullName: app.Seeker.SeekerProfile.FullName,
					Phone:    app.Seeker.SeekerProfile.Phone,
					Address:  app.Seeker.SeekerProfile.Address,
				}
			}
		}
		recentAppsParams = append(recentAppsParams, param)
	}
	stats.RecentApplicants = recentAppsParams

	return stats, nil
}
