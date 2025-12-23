package dto

type CreateJobRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Category    string   `json:"category"`
	JobType     string   `json:"job_type"`
	Salary      string   `json:"salary"`
	Benefits    []string `json:"benefits"`
}

type UpdateJobRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Category    string   `json:"category"`
	JobType     string   `json:"job_type"`
	Salary      string   `json:"salary"`
	Benefits    []string `json:"benefits"`
}
