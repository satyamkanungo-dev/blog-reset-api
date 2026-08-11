package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	apiresponse "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_response"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/service"
)

type IMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
	LogMiddleware() gin.HandlerFunc
}

type Middleware struct {
	authService service.IAuthService
}

func NewMiddleware(authService service.IAuthService) *Middleware {
	return &Middleware{authService: authService}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
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

		userId, err := m.authService.ValidateToken(tokenStr)
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

func (m *Middleware) LogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		log.Printf("[REQUEST] %s %s from %s %s", ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.UserAgent(), ctx.ClientIP())

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()
		log.Printf("[RESPONSE] %s %s -> %d (%s)", ctx.Request.Method, ctx.Request.URL.Path, status, latency)
	}
}
