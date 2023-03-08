// package which contains database model definition
package models

import "gorm.io/gorm"

// Struct to store the ingredient, it embed the gorm model strut which define common fields
type Ingredient struct {
	gorm.Model
	Name    string    `gorm:"type:varchar(100);uninque;not null"` // ingedient name
	Type    string    `gorm:"type:varchar(100);not null"`         // ingredient type
	Recipes []*Recipe `gorm:"many2many:ingredients_recipes;"`     // Reference of each recipe which use this ingredient
}
