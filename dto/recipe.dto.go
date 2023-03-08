// Package dto defines data transfer objects (DTOs) used for communicating between the input and output of an API
package dto

import "github.com/clementb49/welsh_academy/models"

// RecipeReqBody represents the request body for creating a recipe.
type RecipeReqBody struct {
	Title         string `json:"title" xml:"title" binding:"required"`
	Description   string `json:"description" xml:"description" binding:"required"`
	Difficulty    uint8  `json:"difficulty" xml:"difficulty" binding:"required,min=0,max=5"`
	IngredientsId []uint `json:"ingredients_id" xml:"ingredients_id" binding:"required,min=1"`
}

// ConvertToModdel converts a RecipeReqBody to a Recipe model.
func (r *RecipeReqBody) ConvertToModdel() *models.Recipe {
	return &models.Recipe{
		Title:       r.Title,
		Description: r.Description,
		Difficulty:  r.Difficulty,
	}
}

// RecipeResBody represents the response body for a recipe.
type RecipeResBody struct {
	CommonResBody
	Title       string               `json:"title" xml:"title"`
	Description string               `json:"description" xml:"description"`
	Difficulty  uint8                `json:"difficulty" xml:"difficulty"`
	Ingredients []*IngredientResBody `json:"ingredients" xml:"ingredient"`
	AuthorId    uint                 `json:"author_id"`
}

// ConvertFromModel converts a Recipe model to a RecipeResBody.
func (r *RecipeResBody) ConvertFromModel(model *models.Recipe) {
	r.convertFromGormModel(&model.Model)
	r.Ingredients = make([]*IngredientResBody, len(model.Ingredients))
	for i, v := range model.Ingredients {
		dto := &IngredientResBody{}
		dto.ConvertFromModel(v)
		r.Ingredients[i] = dto
	}
	r.Title = model.Title
	r.Description = model.Description
	r.Difficulty = model.Difficulty
	r.AuthorId = uint(model.AuthorID)
}
