package http

import (
	"be-job-portal/internal/domain"
	"be-job-portal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	appUsecase domain.ApplicationUsecase
}

func NewDashboardHandler(appUsecase domain.ApplicationUsecase) *DashboardHandler {
	return &DashboardHandler{
		appUsecase: appUsecase,
	}
}

func (h *DashboardHandler) GetRecruiterStats(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	role, exists := c.Get("role")
	if !exists || role.(string) != "RECRUITER" {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "Only recruiters can access dashboard")
		return
	}

	stats, err := h.appUsecase.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dashboard stats", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard stats fetched successfully", stats)
}
