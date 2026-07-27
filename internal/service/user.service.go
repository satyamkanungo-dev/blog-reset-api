package service

import (
	"context"

	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	apirequest "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_request"
	models "github.com/satyamkanungo-dev/blog-rest-api/internal/models/core"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	passwordLength = 8
)

type UserService struct {
	userRepository repository.IUserRepsoitory
}

func NewUserService(ur repository.IUserRepsoitory) *UserService {
	return &UserService{userRepository: ur}
}

func (us *UserService) Create(rr *apirequest.RegisterRequest, role string) (*models.User, error) {
	// validate the values
	if rr.Email == "" || rr.Password == "" || rr.Name == "" {
		return nil, Error.ErrAllFieldRequired
	}

	if len(rr.Password) < passwordLength {
		return nil, Error.ErrPasswordTooShort
	}

	// check the email already register
	if existingUser, err := us.userRepository.Get(context.Background(), rr.Email); err == nil && existingUser != nil {
		return nil, Error.ErrEmailAlreadyExists
	}

	// generate hashpassword
	hassPassword, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := us.userRepository.Create(context.Background(), rr.Name, rr.Email, string(hassPassword), role)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Get(lr *apirequest.LoginRequest) (*models.User, error) {
	// validate the values
	if lr.Email == "" || lr.Password == "" {
		return nil, Error.ErrAllFieldRequired
	}

	if len(lr.Password) < passwordLength {
		return nil, Error.ErrPasswordTooShort
	}

	// get the user from userRepository
	user, err := us.userRepository.Get(context.Background(), lr.Email)
	if err != nil {
		return nil, err
	}

	// check the Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lr.Password))
	if err != nil {
		return nil, Error.ErrIncorrectPassword
	}

	return &models.User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (us *UserService) Update(ur *apirequest.UpdateUserRequest, userId string) (*models.User, error) {
	// validate the values
	if ur.Name == "" && ur.Password == "" {
		return nil, Error.ErrAllFieldRequired
	}

	if ur.Name == "" || len(ur.Password) >= passwordLength {
		return nil, Error.ErrPasswordTooShort
	}

	user, err := us.userRepository.Update(context.Background(), userId, ur.Name, ur.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
