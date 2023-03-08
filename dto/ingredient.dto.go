// Package dto defines data transfer objects (DTOs) used for communicating between the input and output of an API
package dto

import "github.com/clementb49/welsh_academy/models"

// IngredientReqBody defines the request body for creating or updating an Ingredient
type IngredientReqBody struct {
	Name string `json:"name" xml:"name" binding:"required"`
	Type string `json:"type" xml:"type" binding:"required"`
}

// ConvertToModel converts an IngredientReqBody to a models.Ingredient
func (i *IngredientReqBody) ConvertToModel() *models.Ingredient {
	return &models.Ingredient{
		Name: i.Name,
		Type: i.Type,
	}
}

// IngredientResBody defines the response body for getting an Ingredient
type IngredientResBody struct {
	CommonResBody
	IngredientReqBody
}

// ConvertFromModel converts a models.Ingredient to an IngredientResBody
func (i *IngredientResBody) ConvertFromModel(model *models.Ingredient) {
	i.convertFromGormModel(&model.Model)
	i.Name = model.Name
	i.Type = model.Type
}
