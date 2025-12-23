package repository

import (
	"context"
	"errors"

	"be-job-portal/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) domain.ProfileRepository {
	return &profileRepository{db}
}

func (r *profileRepository) GetSeekerProfile(ctx context.Context, userID uuid.UUID) (*domain.SeekerProfile, error) {
	var profile domain.SeekerProfile
	err := r.db.WithContext(ctx).Preload("Experiences", func(db *gorm.DB) *gorm.DB {
		return db.Order("start_date DESC")
	}).Preload("Educations", func(db *gorm.DB) *gorm.DB {
		return db.Order("start_date DESC")
	}).Where("user_id = ?", userID).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) UpdateSeekerProfile(ctx context.Context, profile *domain.SeekerProfile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.SeekerProfile
		err := tx.Where("user_id = ?", profile.UserID).First(&existing).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(profile).Error
			}
			return err
		}

		profile.ID = existing.ID

		if err := tx.Where("seeker_profile_id = ?", existing.ID).Delete(&domain.Experience{}).Error; err != nil {
			return err
		}
		if err := tx.Where("seeker_profile_id = ?", existing.ID).Delete(&domain.Education{}).Error; err != nil {
			return err
		}

		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(profile).Error
	})
}

func (r *profileRepository) GetCompanyProfile(ctx context.Context, userID uuid.UUID) (*domain.CompanyProfile, error) {
	var profile domain.CompanyProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) UpdateCompanyProfile(ctx context.Context, profile *domain.CompanyProfile) error {
	var existing domain.CompanyProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", profile.UserID).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.WithContext(ctx).Create(profile).Error
		}
		return err
	}

	profile.ID = existing.ID
	return r.db.WithContext(ctx).Save(profile).Error
}
