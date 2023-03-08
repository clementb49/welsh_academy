// The package 'services' contains the business logic for handling route
package services

import (
	"errors"

	"github.com/clementb49/welsh_academy/dto"
	"github.com/clementb49/welsh_academy/repositories"
	"github.com/clementb49/welsh_academy/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// The bcrypt cost used to hash passord
const BCRYPT_COST = 15

// ErrInvalidAuth is returned when the email or password provided are invalid
var ErrInvalidAuth = errors.New("invalid email or password")

// UserService provides an interface for interacting with user data
type UserService interface {
	CreateUser(input *dto.CreateUserReqBody) (*dto.UserResBody, error)
	LoginUser(input *dto.LoginReqBody) (*dto.LoginResBody, error)
	GetUserById(userId uint) (*dto.UserResBody, error)
}

// userService is an implementation of UserService
type userService struct {
	repo   repositories.UserRepository
	logger *zap.Logger
}

// NewUserService creates a new instance of userService
func NewUserService(repo repositories.UserRepository) *userService {
	return &userService{
		repo:   repo,
		logger: zap.L(),
	}
}

// CreateUser creates a new user with the provided input data
func (s *userService) CreateUser(input *dto.CreateUserReqBody) (*dto.UserResBody, error) {
	// CreateUser creates a new user with the provided input data
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), BCRYPT_COST)
	if err != nil {
		return nil, err
	}
	input.Password = string(hashedPassword)
	userModel := input.ConvertToModdel()
	userModel, err = s.repo.CreateUser(userModel)
	if err != nil {
		return nil, err
	}
	userRes := &dto.UserResBody{}
	userRes.ConvertFromModel(userModel)
	return userRes, nil
}

// LoginUser logs in a user with the provided input data
func (s *userService) LoginUser(input *dto.LoginReqBody) (*dto.LoginResBody, error) {
	// Find the user with the provided email
	userModel, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		return nil, ErrInvalidAuth
	}
		// Compare the password provided with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidAuth
	}
		// Generate a JWT token and return it in a response object
	token, err := utils.GetSignedJwt(userModel.ID)
	if err != nil {
		return nil, ErrInvalidAuth
	}
	return &dto.LoginResBody{AccessToken: token}, nil
}
// GetUserById retrieves a user with the provided user ID
func (s *userService) GetUserById(userId uint) (*dto.UserResBody, error) {
	userModel, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	userRes := &dto.UserResBody{}
	userRes.ConvertFromModel(userModel)
	return userRes, nil
}
