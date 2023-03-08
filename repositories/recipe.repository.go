// package repositories defines interfaces for managing recipe data in the database
package repositories

import (
	"fmt"

	"github.com/clementb49/welsh_academy/models" // import models package for Recipe and Ingredient structs
	"go.uber.org/zap"                            // import logging package
	"gorm.io/gorm"                               // import gorm package for database operations
)

// RecipeRepository is an interface that defines functions for managing Recipe data in the database
type RecipeRepository interface {
	CreateRecipe(recipe *models.Recipe, ingredientsid []uint, userId uint) (*models.Recipe, error)
	GetAllRecipes(pageSize int, pageNumber int) ([]*models.Recipe, int64, error)
	GetRecipeById(recipeId uint) (*models.Recipe, error)
	DeleteRecipeById(recipeId uint) error
	AddToFavRecipe(userId, recipeId uint) (*models.Recipe, error)
	DeleteFavRecipe(userId, recipeId uint) error
	GetAllFavRecipes(pageSize int, pageNumber int, userId uint) ([]*models.Recipe, int64, error)
}

// Join query to link recipe and user for favorite recipe
const favoriteJoinQuery = "JOIN wac_favorites_recipes ON wac_favorites_recipes.recipe_id = wac_recipes.id AND wac_favorites_recipes.user_id = ?"

// ErrRecipeNotAcceptable is an error that is returned when a Recipe cannot be created due to missing Ingredients in the database
var ErrRecipeNotAcceptable = fmt.Errorf("missing ingredient in the database, recipe not acceptable")

// NewRecipeRepository is a function that returns an implementation of RecipeRepository using the given database connection
func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &repository{
		db:     db,
		logger: zap.L(),
	}
}

// CreateRecipe is a function that creates a new Recipe record in the database
func (r *repository) CreateRecipe(input *models.Recipe, ingredientsId []uint, userId uint) (*models.Recipe, error) {
	db := r.db.Model(input)
	var ingredients []*models.Ingredient
	result := r.db.Where("id IN ?", ingredientsId).Find(&ingredients)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected != int64(len(ingredientsId)) {
		return nil, ErrRecipeNotAcceptable
	}
	input.Ingredients = ingredients
	result = db.Debug().Save(input)
	if err := result.Error; err != nil {
		return nil, err
	}
	return input, nil
}

// GetAllRecipes returns all recipes with pagination
func (r *repository) GetAllRecipes(pageSize int, pageNumber int) ([]*models.Recipe, int64, error) {
	var recipes []*models.Recipe
	var totalRecipes int64
	db := r.db.Model(&models.Recipe{})
	result := db.Count(&totalRecipes)
	if err := result.Error; err != nil {
		return nil, 0, err
	}
	result = r.db.Offset(pageNumber * pageSize).Limit(pageSize).Find(&recipes)
	if err := result.Error; err != nil {
		return nil, 0, err
	}
	return recipes, totalRecipes, nil
}

// GetRecipeById return a recipe by ID
func (r *repository) GetRecipeById(recipeId uint) (*models.Recipe, error) {
	var recipe *models.Recipe
	result := r.db.First(&recipe, recipeId)
	if err := result.Error; err != nil {
		return nil, err
	}
	return recipe, nil
}

// DeleteRecipeById delete a recipe by ID
func (r *repository) DeleteRecipeById(recipeId uint) error {
	result := r.db.Delete(&models.Recipe{}, recipeId)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// AddToFavRecipe add the recipe to the user favorite
func (r *repository) AddToFavRecipe(userId, recipeId uint) (*models.Recipe, error) {
	var recipe *models.Recipe
	result := r.db.First(&recipe, recipeId)
	if err := result.Error; err != nil {
		return nil, err
	}
	var user *models.User
	result = r.db.First(&user, userId)
	if err := result.Error; err != nil {
		return nil, err
	}
	err := r.db.Model(&recipe).Association("LikedUser").Append(user)
	if err != nil {
		return nil, err
	}
	result = r.db.Save(&recipe)
	if err := result.Error; err != nil {
		return nil, err
	}
	return recipe, nil
}

// DeleteToFavRecipe delete the recipe to the user favorite
func (r *repository) DeleteFavRecipe(userId, recipeId uint) error {
	var recipe *models.Recipe
	result := r.db.First(&recipe, recipeId)
	if err := result.Error; err != nil {
		return err
	}
	var user *models.User
	result = r.db.First(&user, userId)
	if err := result.Error; err != nil {
		return err
	}
	err := r.db.Model(&recipe).Association("LikedUser").Delete(user)
	if err != nil {
		return err
	}
	result = r.db.Save(&recipe)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

	// GetAllFavRecipe returns all favorite recipe with pagination
func (r *repository) GetAllFavRecipes(pageSize int, pageNumber int, userId uint) ([]*models.Recipe, int64, error) {
	var recipes []*models.Recipe
	var totalRecipes int64
	db := r.db.Model(&models.Recipe{}).Joins(favoriteJoinQuery, userId)
	err := db.Count(&totalRecipes).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Offset(pageNumber * pageSize).Limit(pageSize).Find(&recipes).Error
	if err != nil {
		return nil, 0, err
	}
	return recipes, totalRecipes, nil
}
