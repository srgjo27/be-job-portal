package http

import (
	"net/http"
	"strconv"

	"be-job-portal/internal/delivery/http/dto"
	"be-job-portal/internal/domain"
	"be-job-portal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobHandler struct {
	jobUsecase domain.JobUsecase
}

func NewJobHandler(us domain.JobUsecase) *JobHandler {
	return &JobHandler{
		jobUsecase: us,
	}
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	var input dto.CreateJobRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	role, exists := c.Get("role")
	if !exists || role.(string) != "RECRUITER" {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "Only recruiters can create jobs")
		return
	}

	err = h.jobUsecase.CreateJob(c.Request.Context(), input.Title, input.Description, input.Category, input.JobType, input.Salary, input.Benefits, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create job", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Job created successfully", nil)
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	page := 1
	limit := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 100 {
				limit = 100
			}
		}
	}

	params := domain.PaginationParams{
		Page:  page,
		Limit: limit,
	}

	result, err := h.jobUsecase.ListJobs(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch jobs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Jobs fetched successfully", result)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid job ID", err.Error())
		return
	}

	job, err := h.jobUsecase.GetJob(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch job", err.Error())
		return
	}
	if job == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Job not found", "Job with given ID does not exist")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Job fetched successfully", job)
}

func (h *JobHandler) ListJobsByRecruiter(c *gin.Context) {
	// Get recruiter_id from query parameter (optional)
	recruiterIDStr := c.Query("recruiter_id")

	var recruiterID uuid.UUID
	var err error

	if recruiterIDStr != "" {
		recruiterID, err = uuid.Parse(recruiterIDStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid recruiter ID", err.Error())
			return
		}
	} else {
		recruiterID, err = utils.GetUserID(c)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
			return
		}
	}

	jobs, err := h.jobUsecase.ListJobsByRecruiter(c.Request.Context(), recruiterID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch jobs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Jobs fetched successfully", jobs)
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid job ID", err.Error())
		return
	}

	var input dto.UpdateJobRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	err = h.jobUsecase.UpdateJob(c.Request.Context(), id, userID, input.Title, input.Description, input.Category, input.JobType, input.Salary, input.Benefits)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Job not found", "Job with given ID does not exist")
		case domain.ErrUnauthorized:
			utils.ErrorResponse(c, http.StatusForbidden, "Unauthorized", "You are not authorized to update this job")
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update job", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Job updated successfully", nil)
}
