package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/clementb49/welsh_academy/dto"
	"github.com/clementb49/welsh_academy/models"
	"github.com/clementb49/welsh_academy/services"
	"github.com/clementb49/welsh_academy/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepository struct{}

func (m *mockUserRepository) CreateUser(userModel *models.User) (*models.User, error) {
	if userModel.Email == "existing_user@example.com" {
		return nil, gorm.ErrDuplicatedKey
	}

	userModel.ID = 1
	userModel.CreatedAt = time.Now()
	userModel.UpdatedAt = time.Now()
	return userModel, nil
}

func (m *mockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	if email == "existing_user@example.com" {
		hashed_password, _ := bcrypt.GenerateFromPassword([]byte("password"), services.BCRYPT_COST)

		return &models.User{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:    "existing_user@example.com",
			Password: string(hashed_password),
		}, nil
	}

	return nil, errors.New("user not found")
}

func (m *mockUserRepository) GetUserById(userId uint) (*models.User, error) {
	if userId == 1 {
		return &models.User{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:    "test@example.com",
			Password: "hashed_password",
		}, nil
	}

	return nil, gorm.ErrRecordNotFound
}

func TestCreateUser(t *testing.T) {
	repo := &mockUserRepository{}
	userService := services.NewUserService(repo)

	// test happy path
	input := &dto.CreateUserReqBody{
		LoginReqBody: dto.LoginReqBody{
			Email:    "test@example.com",
			Password: "hashed_password",
		},
		FirstName: "toto",
		LastName:  "tata",
	}
	userRes, err := userService.CreateUser(input)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), userRes.ID)
	assert.Equal(t, "test@example.com", userRes.Email)
	assert.Equal(t, "toto", userRes.FirstName)
	assert.Equal(t, "tata", userRes.LastName)

	// Test error: email already exists
	input.Email = "existing_user@example.com"
	userRes, err = userService.CreateUser(input)
	assert.Error(t, err)
	assert.Nil(t, userRes)
}

func TestLoginUser(t *testing.T) {
	repo := &mockUserRepository{}
	userService := services.NewUserService(repo)

	// test happy path
	input := &dto.LoginReqBody{
		Email:    "existing_user@example.com",
		Password: "password",
	}
	loginRes, err := userService.LoginUser(input)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginRes.AccessToken)
	claims, err := utils.VerifyToken(loginRes.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, "1", claims.ID)

	// Test error: invalid password
	input.Password = "wrong_password"
	loginRes, err = userService.LoginUser(input)
	assert.Error(t, err)
	assert.Nil(t, loginRes)

	// Test error: user not found
	input.Email = "non_existing_user@example.com"
	loginRes, err = userService.LoginUser(input)
	assert.Error(t, err)
	assert.Nil(t, loginRes)
}

func TestGetUserById(t *testing.T) {
	repo := &mockUserRepository{}
	userService := services.NewUserService(repo)

	// Test happy path
	userRes, err := userService.GetUserById(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), userRes.ID)
	assert.Equal(t, "test@example.com", userRes.Email)

	// Test error: user not found
	userRes, err = userService.GetUserById(2)
	assert.Error(t, err)
	assert.Nil(t, userRes)
}
