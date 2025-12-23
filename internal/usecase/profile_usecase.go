package usecase

import (
	"context"
	"errors"

	"strings"

	"be-job-portal/internal/domain"

	"github.com/google/uuid"
)

type profileUsecase struct {
	profileRepo domain.ProfileRepository
}

func NewProfileUsecase(profileRepo domain.ProfileRepository) domain.ProfileUsecase {
	return &profileUsecase{profileRepo}
}

func (u *profileUsecase) GetProfile(ctx context.Context, userID uuid.UUID, role string) (interface{}, error) {
	switch strings.ToUpper(role) {
	case "SEEKER":
		profile, err := u.profileRepo.GetSeekerProfile(ctx, userID)
		if err != nil {
			return nil, err
		}
		if profile == nil {
			return &domain.SeekerProfile{UserID: userID}, nil
		}
		return profile, nil
	case "RECRUITER":
		profile, err := u.profileRepo.GetCompanyProfile(ctx, userID)
		if err != nil {
			return nil, err
		}
		if profile == nil {
			return &domain.CompanyProfile{UserID: userID}, nil
		}
		return profile, nil
	}
	return nil, errors.New("invalid role")
}

func (u *profileUsecase) UpdateSeekerProfile(ctx context.Context, userID uuid.UUID, profile *domain.SeekerProfile) error {
	profile.UserID = userID
	return u.profileRepo.UpdateSeekerProfile(ctx, profile)
}

func (u *profileUsecase) UpdateCompanyProfile(ctx context.Context, userID uuid.UUID, profile *domain.CompanyProfile) error {
	profile.UserID = userID
	return u.profileRepo.UpdateCompanyProfile(ctx, profile)
}
