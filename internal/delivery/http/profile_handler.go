package http

import (
	"net/http"
	"time"

	"strings"

	"be-job-portal/internal/delivery/http/dto"
	"be-job-portal/internal/domain"
	"be-job-portal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileUsecase domain.ProfileUsecase
	userRepo       domain.UserRepository
}

func NewProfileHandler(us domain.ProfileUsecase, ur domain.UserRepository) *ProfileHandler {
	return &ProfileHandler{
		profileUsecase: us,
		userRepo:       ur,
	}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user", err.Error())
		return
	}

	profile, err := h.profileUsecase.GetProfile(c.Request.Context(), userID, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch profile", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile fetched successfully", profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user", err.Error())
		return
	}

	switch strings.ToUpper(user.Role) {
	case "SEEKER":
		var input dto.UpdateSeekerProfileRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}

		var experiences []domain.Experience
		for _, exp := range input.Experiences {
			startDate, _ := time.Parse("2006-01-02", exp.StartDate)
			var endDate *time.Time
			if exp.EndDate != "" {
				parsedEndDate, _ := time.Parse("2006-01-02", exp.EndDate)
				endDate = &parsedEndDate
			}
			experiences = append(experiences, domain.Experience{
				Title:       exp.Title,
				Company:     exp.Company,
				Location:    exp.Location,
				StartDate:   startDate,
				EndDate:     endDate,
				Description: exp.Description,
			})
		}

		var educations []domain.Education
		for _, edu := range input.Educations {
			startDate, _ := time.Parse("2006-01-02", edu.StartDate)
			var endDate *time.Time
			if edu.EndDate != "" {
				parsedEndDate, _ := time.Parse("2006-01-02", edu.EndDate)
				endDate = &parsedEndDate
			}
			educations = append(educations, domain.Education{
				Institution:  edu.Institution,
				Degree:       edu.Degree,
				FieldOfStudy: edu.FieldOfStudy,
				StartDate:    startDate,
				EndDate:      endDate,
			})
		}

		profile := &domain.SeekerProfile{
			FullName:     input.FullName,
			Phone:        input.Phone,
			Address:      input.Address,
			ResumeURL:    input.ResumeURL,
			PortfolioURL: input.PortfolioURL,
			LinkedInURL:  input.LinkedInURL,
			Skills:       input.Skills,
			Description:  input.Description,
			Experiences:  experiences,
			Educations:   educations,
		}
		err = h.profileUsecase.UpdateSeekerProfile(c.Request.Context(), userID, profile)

	case "RECRUITER":
		var input dto.UpdateCompanyProfileRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}

		profile := &domain.CompanyProfile{
			CompanyName: input.CompanyName,
			Website:     input.Website,
			Location:    input.Location,
			Phone:       input.Phone,
			Description: input.Description,
			LogoURL:     input.LogoURL,
		}
		err = h.profileUsecase.UpdateCompanyProfile(c.Request.Context(), userID, profile)

	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role", "User role is not recognized")
		return
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", nil)
}
