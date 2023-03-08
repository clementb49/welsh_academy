// package repositories defines interfaces for managing user data in the database
package repositories

import (
	"github.com/clementb49/welsh_academy/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Define UserRepository interface with three methods
type UserRepository interface {
	CreateUser(input *models.User) (*models.User, error)
	GetUserById(userId uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

// Define a function to create a new UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &repository{
		db:     db,
		logger: zap.L(),
	}
}

// This function creates a new user in the database by taking a pointer to a user model
// as input and returns a pointer to the same user model along with any error encountered
func (r *repository) CreateUser(input *models.User) (*models.User, error) {
	db := r.db.Model(input)              // Get the database model for the input user
	result := db.Create(input)           // Create a new user record in the database
	if err := result.Error; err != nil { // Check for any errors in the result of the create operation
		return nil, err
	}
	return input, nil // Return the input user model with no errors
}

// This function retrieves a user from the database by taking a user ID as input
// and returns a pointer to the retrieved user model along with any error encountered
func (r *repository) GetUserById(userID uint) (*models.User, error) {
	var user models.User                 // Create a variable to hold the retrieved user model
	result := r.db.First(&user, userID)  // Retrieve the user with the given user ID from the database
	if err := result.Error; err != nil { // Check for any errors in the result of the retrieve operation
		return nil, err
	}
	return &user, nil // Return a pointer to the retrieved user model with no errors
}

// This function retrieves a user from the database by taking an email address as input
// and returns a pointer to the retrieved user model along with any error encountered
func (r *repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User                            // Create a variable to hold the retrieved user model
	result := r.db.First(&user, "email = ?", email) // Retrieve the user with the given email address from the database
	if err := result.Error; err != nil {            // Check for any errors in the result of the retrieve operation
		return nil, err
	}
	return &user, nil // Return a pointer to the retrieved user model with no errors
}
