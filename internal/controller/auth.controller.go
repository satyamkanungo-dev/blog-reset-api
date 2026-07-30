package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/auth"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/middleware"
	apiresponse "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_response"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/service"
)

type authController struct {
	authService service.IAuthService
}

func NewAuthController(service service.IAuthService) *authController {
	return &authController{authService: service}
}

func (ac *authController) RegisterRoutes(r gin.IRouter, middleware middleware.IAuthMiddleware) {
	auth := r.Group("/users/refresh").Use(middleware.RefreshMiddleware())
	{
		auth.POST("", ac.Get)
	}
}

func (ac *authController) Get(ctx *gin.Context) {
	userId, _ := getUserIdFromMiddleware(ctx)

	accessToken, err := ac.authService.GetToken(auth.AccessTokenExp, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return

	}

	refreshToken, err := ac.authService.GetToken(auth.RefreshTokenExp, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return

	}

	ctx.JSON(http.StatusOK, apiresponse.APIResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data: apiresponse.AuthResponse{
			RefreshResponse: apiresponse.RefreshResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
			User: apiresponse.UserResponse{
				Id: userId,
			},
		},
	})
}
