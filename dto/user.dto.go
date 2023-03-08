package dto

import "github.com/clementb49/welsh_academy/models"

// LoginReqBody represents the request body for logging in.
type LoginReqBody struct {
	Email    string `json:"email" xml:"email" binding:"required,email"`
	Password string `json:"password" xml:"password" binding:"required,alphanumunicode,min=8"`
}

// CreateUserReqBody represents the request body for creating a user.
type CreateUserReqBody struct {
	FirstName string `json:"first_name" xml:"firstName" binding:"required"`
	LastName  string `json:"last_name" xml:"lastName" binding:"required"`
	LoginReqBody
}

// ConvertToModdel converts CreateUserReqBody to models.User.
func (u *CreateUserReqBody) ConvertToModdel() *models.User {
	return &models.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

// UserResBody represents the response body for a user.
type UserResBody struct {
	CommonResBody
	FirstName string `json:"first_name" xml:"firstName"`
	LastName  string `json:"last_name" xml:"lastName"`
	Email     string `json:"email" xml:"email"`
}

// ConvertFromModel converts models.User to UserResBody.
func (u *UserResBody) ConvertFromModel(model *models.User) {
	u.convertFromGormModel(&model.Model)
	u.FirstName = model.FirstName
	u.LastName = model.LastName
	u.Email = model.Email
}

// LoginResBody represents the response body for logging in.
type LoginResBody struct {
	AccessToken string `json:"access_token" xml:"access_token"`
}
