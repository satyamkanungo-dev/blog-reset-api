package service

import (
	"time"

	apirequest "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_request"
	apiresponse "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_response"
	models "github.com/satyamkanungo-dev/blog-rest-api/internal/models/core"
)

type IUserService interface {
	Create(rr *apirequest.RegisterRequest, role string) (*models.User, error)
	Update(ur *apirequest.UpdateUserRequest, userId string) (*models.User, error)
	Get(lr *apirequest.LoginRequest) (*models.User, error)
}

type IBlogService interface {
	Create(br *apirequest.BlogRequest, userId string) (*models.Blog, error)
	Get(id, userId string) (*models.Blog, error)
	Update(br *apirequest.BlogRequest, id, userId string) (*models.Blog, error)
	Delete(id, userId string) error
	GetAll(userId, cursor string) (*apiresponse.BlogsResponse, error)
	DeleteMultiple(db *apirequest.DeleteBlogRequest, userId string) (*apiresponse.BulkDeleteResponse, error)
}

type IAuthService interface {
	GetToken(exp time.Time, userId string) (string, error)
	ValidateToken(tokenStr string) (string, error)
}
