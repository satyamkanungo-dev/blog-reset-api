package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	apiresponse "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_response"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/service"
)

type IAuthMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
}

type Middleware struct {
	authService service.IAuthService
}

func NewMiddleware(authService service.IAuthService) *Middleware {
	return &Middleware{authService: authService}
}

func (am *Middleware) AuthMiddleware() gin.HandlerFunc {
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
