// This package make the interface between the database and the service
package repositories

import (
	"github.com/clementb49/welsh_academy/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IngredientRepository is an interface for interacting with the Ingredient model
type IngredientRepository interface {
	CreateIngredient(input *models.Ingredient) (*models.Ingredient, error)               // Create a new ingredient
	GetAllIngredients(pageSize int, pageNumber int) ([]*models.Ingredient, int64, error) // Get all ingredients with pagination
	GetIngredientById(ingredientId uint) (*models.Ingredient, error)                       // Get an ingredient by ID
	DeleteIngredientById(ingredientId uint) error                                        // Delete an ingredient by ID
}

// NewIngredientRepository returns a new instance of the IngredientRepository interface
func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return &repository{
		db:     db,
		logger: zap.L(),
	}
}

// CreateIngredient creates a new ingredient in the database
func (r *repository) CreateIngredient(input *models.Ingredient) (*models.Ingredient, error) {
	db := r.db.Model(input)
	result := db.Create(input)
	if err := result.Error; err != nil {
		return nil, err
	}
	return input, nil
}

// GetAllIngredients returns all ingredients with pagination
func (r *repository) GetAllIngredients(pageSize int, pageNumber int) ([]*models.Ingredient, int64, error) {
	var ingredients []*models.Ingredient
	var totalIngredient int64
	db := r.db.Model(&models.Ingredient{})
	result := db.Count(&totalIngredient)
	if err := result.Error; err != nil {
		return nil, 0, err
	}
	result = r.db.Offset(pageNumber * pageSize).Limit(pageSize).Find(&ingredients)
	if err := result.Error; err != nil {
		return nil, 0, err
	}
	return ingredients, totalIngredient, nil
}

// GetIngredientById returns an ingredient by ID
func (r *repository) GetIngredientById(ingredientId uint) (*models.Ingredient, error) {
	var ingredient *models.Ingredient
	result := r.db.First(&ingredient, ingredientId)
	if err := result.Error; err != nil {
		return nil, err
	}
	return ingredient, nil
}

// DeleteIngredientById deletes an ingredient by ID
func (r *repository) DeleteIngredientById(ingredientId uint) error {
	result := r.db.Delete(&models.Ingredient{}, ingredientId)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
