package models

import "gorm.io/gorm"

// Struct to store the user, it embed the gorm model strut which define common fields
type User struct {
	gorm.Model
	FirstName      string    `gorm:"type:varchar(100);not null"` // the user first name 
	LastName       string    `gorm:"type:varchar(100);not null"` // the user last name 
	Email          string    `gorm:"type:varchar(255);unique;not null"` // the user email 
	Password       string    `gorm:"type:char(60);not null"` // the hashed version of the user password
	FavRecipes     []*Recipe `gorm:"many2many:favorites_recipes;"` // the favorite recipe for the user 
	CreatedRecipes []*Recipe `gorm:"foreignKey:AuthorID"` // the recipe created by the user 
}
