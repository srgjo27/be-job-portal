package http

import (
	"be-job-portal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authHandler *AuthHandler, jobHandler *JobHandler, appHandler *ApplicationHandler, profileHandler *ProfileHandler, dashboardHandler *DashboardHandler) {
	// Auth Routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/google/callback", authHandler.GoogleCallback)
	}

	// Job Routes
	jobs := r.Group("/api/jobs")
	jobs.Use(utils.AuthMiddleware())
	{
		jobs.POST("", jobHandler.CreateJob)
		jobs.GET("", jobHandler.ListJobs)
		jobs.GET("/recruiter", jobHandler.ListJobsByRecruiter)
		jobs.GET("/:id", jobHandler.GetJob)
		jobs.PUT("/:id", jobHandler.UpdateJob)
		jobs.GET("/:id/applicants", appHandler.ListJobApplicants)
	}

	// Application Routes
	apps := r.Group("/api/applications")
	apps.Use(utils.AuthMiddleware())
	{
		apps.POST("", appHandler.ApplyJob)
		apps.GET("", appHandler.ListApplications)
		apps.PUT("/:id/status", appHandler.UpdateStatus)
	}

	// Profile Routes
	profile := r.Group("/api/profile")
	profile.Use(utils.AuthMiddleware())
	{
		profile.GET("", profileHandler.GetProfile)
		profile.PUT("", profileHandler.UpdateProfile)
	}

	// Dashboard Routes
	dashboard := r.Group("/api/dashboard")
	dashboard.Use(utils.AuthMiddleware())
	{
		dashboard.GET("/stats", dashboardHandler.GetRecruiterStats)
	}
}
