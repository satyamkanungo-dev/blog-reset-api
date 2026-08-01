package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/auth"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/middleware"
	apirequest "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_request"
	apiresponse "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_response"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/service"
)

var (
	user = "user"
)

type UserController struct {
	UserService service.IUserService
	AuthService service.IAuthService
}

func NewUserController(userservice service.IUserService, authservice service.IAuthService) *UserController {
	return &UserController{UserService: userservice, AuthService: authservice}
}

func (uc *UserController) RegisterRoutes(r gin.IRouter, middleware middleware.IAuthMiddleware) {
	users := r.Group("/users")
	{
		users.POST("/register", uc.Create)
		users.POST("/login", uc.Get)
		users.PUT("", middleware.AccessMiddleware(), uc.Update)
	}
}

func (uc *UserController) Create(ctx *gin.Context) {
	var input apirequest.RegisterRequest
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	user, err := uc.UserService.Create(&input, user)
	if err != nil {
		if errors.Is(err, Error.ErrEmailAlreadyExists) {
			ctx.JSON(http.StatusConflict, apiresponse.APIResponse{
				Code:    http.StatusConflict,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, Error.ErrAllFieldRequired) || errors.Is(err, Error.ErrPasswordTooShort) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// access token
	accessToken, err := uc.AuthService.GetToken(auth.AccessTokenExp, user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// refresh token
	refreshToken, err := uc.AuthService.GetToken(auth.RefreshTokenExp, user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, apiresponse.APIResponse{
		Code:   http.StatusCreated,
		Status: "success",
		Data: apiresponse.AuthResponse{
			RefreshResponse: apiresponse.RefreshResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
			User: apiresponse.UserResponse{
				Id:       user.Id,
				UserName: user.Name,
				Email:    user.Email,
				Role:     user.Role,
			},
		},
	})
}

func (uc *UserController) Get(ctx *gin.Context) {
	var input apirequest.LoginRequest
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	user, err := uc.UserService.Get(&input)
	if err != nil {
		if errors.Is(err, Error.ErrAllFieldRequired) || errors.Is(err, Error.ErrPasswordTooShort) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, Error.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, apiresponse.APIResponse{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, Error.ErrIncorrectPassword) {
			ctx.JSON(http.StatusUnauthorized, apiresponse.APIResponse{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return

	}

	// access token
	accessToken, err := uc.AuthService.GetToken(auth.AccessTokenExp, user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// refresh token
	refreshToken, err := uc.AuthService.GetToken(auth.RefreshTokenExp, user.Id)
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
				Id:       user.Id,
				UserName: user.Name,
				Email:    user.Email,
				Role:     user.Role,
			},
		},
	})
}

func (uc *UserController) Update(ctx *gin.Context) {
	var input apirequest.UpdateUserRequest
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	userId, _ := getUserIdFromMiddleware(ctx)

	user, err := uc.UserService.Update(&input, userId)
	if err != nil {
		if errors.Is(err, Error.ErrAllFieldRequired) || errors.Is(err, Error.ErrPasswordTooShort) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, Error.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, apiresponse.APIResponse{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

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
		Data:   user,
	})
}
