// package which contains database model definition
package models

import "gorm.io/gorm"

// Struct to store the recipe, it embed the gorm model strut which define common fields
type Recipe struct {
	gorm.Model
	Title       string        `gorm:"type:varchar(200);unique;not null"` // the recipe title
	Description string        `gorm:"not null"`                          // text for the recipe description
	Difficulty  uint8         `gorm:"not null;check:difficulty <= 5"`    // the defficuty of the recipe
	Ingredients []*Ingredient `gorm:"many2many:ingredients_recipes;"`    // the ingredient required to make the recipe
	LikedUser   []*User       `gorm:"many2many:favorites_recipes;"`      // the users who liked the recipe
	AuthorID    uint64        // the refence of the user who created the recipe
}
