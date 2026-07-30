package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	apirequest "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_request"
	apiresponse "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_response"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/service"
)

type IAuthMiddleware interface {
	AccessMiddleware() gin.HandlerFunc
	RefreshMiddleware() gin.HandlerFunc
}

type AuthMiddleware struct {
	authService service.IAuthService
}

func NewAuthMiddleware(authService service.IAuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (am *AuthMiddleware) AccessMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, apiresponse.APIResponse{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: Error.ErrAuthorizationToken.Error(),
			})
			return
		}

		userId, err := am.authService.ValidateToken(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, apiresponse.APIResponse{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}

func (am *AuthMiddleware) RefreshMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input apirequest.RefreshTokenRequest
		if err := ctx.BindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		userId, err := am.authService.ValidateToken(input.Token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, apiresponse.APIResponse{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}
