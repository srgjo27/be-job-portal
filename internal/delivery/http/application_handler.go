package http

import (
	"net/http"
	"time"

	"be-job-portal/internal/delivery/http/dto"
	"be-job-portal/internal/domain"
	"be-job-portal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplicationHandler struct {
	appUsecase domain.ApplicationUsecase
}

func NewApplicationHandler(us domain.ApplicationUsecase) *ApplicationHandler {
	return &ApplicationHandler{
		appUsecase: us,
	}
}

func (h *ApplicationHandler) ApplyJob(c *gin.Context) {
	var input dto.ApplyJobRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	jobID, err := uuid.Parse(input.JobID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Job ID", err.Error())
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	err = h.appUsecase.ApplyJob(c.Request.Context(), jobID, userID, input.ResumeURL, input.CoverLetter, input.LinkedInURL, input.PortfolioURL)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to apply for job", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Applied successfully", nil)
}

func (h *ApplicationHandler) ListApplications(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	role := "SEEKER"

	apps, err := h.appUsecase.ListApplications(c.Request.Context(), userID, role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch applications", err.Error())
		return
	}

	var response []dto.ApplicationListResponse
	for _, app := range apps {
		var deletedAt *time.Time
		if app.DeletedAt.Valid {
			deletedAt = &app.DeletedAt.Time
		}

		jobLite := dto.JobLiteResponse{
			ID:       app.JobID,
			Title:    "",
			Category: "",
			JobType:  "",
		}

		if app.Job != nil {
			jobLite.ID = app.Job.ID
			jobLite.Title = app.Job.Title
			jobLite.CompanyName = app.Job.Company.CompanyName
			jobLite.Category = app.Job.Category
			jobLite.JobType = app.Job.JobType
		}

		response = append(response, dto.ApplicationListResponse{
			ID:           app.ID,
			CreatedAt:    app.CreatedAt,
			UpdatedAt:    app.UpdatedAt,
			DeletedAt:    deletedAt,
			JobID:        app.JobID,
			Job:          jobLite,
			SeekerID:     app.SeekerID,
			Status:       app.Status,
			ResumeURL:    app.ResumeURL,
			CoverLetter:  app.CoverLetter,
			LinkedInURL:  app.LinkedInURL,
			PortfolioURL: app.PortfolioURL,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Applications fetched successfully", response)
}

func (h *ApplicationHandler) ListJobApplicants(c *gin.Context) {
	jobIDStr := c.Param("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Job ID", err.Error())
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	apps, err := h.appUsecase.ListJobApplicants(c.Request.Context(), jobID, userID)
	if err != nil {
		if err == domain.ErrUnauthorized {
			utils.ErrorResponse(c, http.StatusForbidden, "Unauthorized", "You are not authorized to view applicants for this job")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch applicants", err.Error())
		return
	}

	var response []gin.H
	for _, app := range apps {
		seekerName := "Unknown"
		seekerEmail := ""
		if app.Seeker != nil {
			seekerEmail = app.Seeker.Email
			if app.Seeker.SeekerProfile != nil {
				seekerName = app.Seeker.SeekerProfile.FullName
			}
		}

		response = append(response, gin.H{
			"id":         app.ID,
			"applied_at": app.CreatedAt,
			"status":     app.Status,
			"seeker": gin.H{
				"id":        app.SeekerID,
				"email":     seekerEmail,
				"full_name": seekerName,
				"phone":     app.Seeker.SeekerProfile.Phone,
			},
			"resume_url":   app.ResumeURL,
			"cover_letter": app.CoverLetter,
			"linkedin_url": app.LinkedInURL,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Applicants fetched successfully", response)
}

func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Application ID", err.Error())
		return
	}

	var input dto.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	err = h.appUsecase.UpdateStatus(c.Request.Context(), appID, userID, input.Status)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			utils.ErrorResponse(c, http.StatusForbidden, "Unauthorized", "You are not authorized to update this application")
		case domain.ErrBadRequest:
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status", "Status must be PENDING, PROCESS, ACCEPTED, or REJECTED")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update status", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Application status updated successfully", nil)
}
