package service

import (
	apirequest "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_request"
	models "github.com/satyamkanungo-dev/blog-rest-api/internal/models/core"
)

type IUserService interface {
	Create(rr *apirequest.RegisterRequest) (*models.User, error)
	Update(ur *apirequest.UpdateUserRequest, userId string) (*models.User, error)
	Get(lr *apirequest.LoginRequest) (*models.User, error)
}
