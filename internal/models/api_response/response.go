package apiresponse

import models "github.com/satyamkanungo-dev/blog-rest-api/internal/models/core"

// core
type APIResponse struct {
	Code    int    `json:"code"`   // http status
	Status  string `json:"status"` // success or error
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// user
type UserResponse struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// blog
type BlogsResponse struct {
	Data       []models.Blog `json:"data"`
	NextCursor any           `json:"next_cursor"`
	HasNext    bool          `json:"has_next"`
}

type BulkDeleteSummary struct {
	Total   int `json:"total"`
	Deleted int `json:"deleted"`
	Failed  int `json:"failed"`
}

type BulkDeleteFailed struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

type BulkDeleteResponse struct {
	Succeeded []string           `json:"succeeded"`
	Failed    []BulkDeleteFailed `json:"failed"`
	Summary   BulkDeleteSummary  `json:"summary"`
}
