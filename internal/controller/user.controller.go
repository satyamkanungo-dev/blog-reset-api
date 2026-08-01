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
		users.POST("/refresh", middleware.AuthMiddleware(), uc.GetToken)
		users.PUT("", middleware.AuthMiddleware(), uc.Update)
	}
}

// Create  		 godoc
// @Summary      Register a new user
// @Description  Creates a new user account and returns access and refresh tokens along with user details
// @Tags         users
// @Accept       json
// @Produce      json
// @Security	 BearerAuth
// @Param        input  body      apirequest.RegisterRequest  true  "User registration details"
// @Success      201    {object}  apiresponse.APIResponse{data=apiresponse.AuthResponse}  "User created successfully"
// @Failure      400    {object}  apiresponse.APIResponse  "Invalid input, missing required fields, or password too short"
// @Failure      409    {object}  apiresponse.APIResponse  "Email already exists"
// @Failure      500    {object}  apiresponse.APIResponse  "Internal server error"
// @Router       /users/register [Post]
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

// Get 				godoc
// @Summary 		Login a user
// @Description		Get a user along with access and refresh token
// @Tags 			users
// @Accept			json
// @Produce 		json
// @Security	 BearerAuth
// @Param			input	body 		apirequest.LoginRequest		true 	"User login details"
// @Success			200		{object}	apiresponse.APIResponse{data=apiresponse.AuthResponse}	"User login successfully"
// @Failure			400		{object}	apiresponse.APIResponse			"Invalid input, missing required fields, or password too short"
// @Failure			401 	{object}	apiresponse.APIResponse			"Unauthorized"
// @Failure			404		{object}	apiresponse.APIResponse			"Invalid user"
// @Failure			500		{object}	apiresponse.APIResponse 		"Internal server error"
// @Router			/users/login 	[Get]
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

// Update 			godoc
// @Summary 		Update users details
// @Description		Get a updated user
// @Tags 			users
// @Accept			json
// @Produce 		json
// @Security	 	BearerAuth
// @Param			input 	body 		apirequest.UpdateUserRequest true "Update user details"
// @Success 		200		{object}	apiresponse.APIResponse{data=models.User}	 "User details updated successfully"
// @Failure			400		{object}	apiresponse.APIResponse	"Invalid input, missing required fields, or password too short"
// @Failure 		401		{object}	apiresponse.APIResponse "Unauthorized"
// @Failure			404		{object}	apiresponse.APIResponse	"Invalid user"
// @Failure 		500 	{object} 	apiresponse.APIResponse	"Internal server error"
// @Router			/users [Put]
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

// GetToken			godoc
// @Summary			Get a new authentication tokens
// @Description		Exchanges a valid refresh token for a new access token and a newly generated refresh token. The refresh token must be provided in the Authorization header as a Bearer token.
// @Tags			users
// @Produce 		json
// @Security		BearerAuth
// @Param			Authorization	header		string		true		"Refresh token, format: Bearer <refresh-token>"
// @Success			200		{object}	apiresponse.APIResponse{data=apiresponse.AuthResponse}
// @Failure			400		{object}	apiresponse.APIResponse			"Invalid input, missing required fields, or malformed token"
// @Failure			401 	{object}	apiresponse.APIResponse			"Unauthorized, invalid or expired refresh token"
// @Failure			500		{object}	apiresponse.APIResponse 		"Internal server error"
// @Router 			/users/refresh [Post]
func (uc *UserController) GetToken(ctx *gin.Context) {
	userId, _ := getUserIdFromMiddleware(ctx)

	accessToken, err := uc.AuthService.GetToken(auth.AccessTokenExp, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return

	}

	refreshToken, err := uc.AuthService.GetToken(auth.RefreshTokenExp, userId)
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
		},
	})
}
