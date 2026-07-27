package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/auth"
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

func (uc *UserController) RegisterRoutes(r gin.IRouter) {
	users := r.Group("/users")
	{
		users.POST("/register", uc.Create)
		users.GET("/login", uc.Get)
		users.PUT("", uc.Update)
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
		log.Println("from accessToken")
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
		log.Println("from refreshToken")
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
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
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
		log.Println("from accessToken")

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
		log.Println("from refreshToken")
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
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
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

	// TODO: needed work
	// demo
	var userId string

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
