// The package 'services' contains the business logic for handling route
package services

import (
	"github.com/clementb49/welsh_academy/dto"
	"github.com/clementb49/welsh_academy/repositories"
	"go.uber.org/zap"
)

// RecipeService is an interface for defining the methods to manage recipes
type RecipeService interface {
	CreateRecipe(input *dto.RecipeReqBody, userId uint) (*dto.RecipeResBody, error)
	GetAllRecipes(input *dto.CommonQueryPage) (*dto.CommonPageRespBody, error)
	GetRecipeById(input *dto.CommonIdPathUri) (*dto.RecipeResBody, error)
	DeleteRecipeById(input *dto.CommonIdPathUri) error
	AddToFavRecipe(userId uint, input *dto.CommonIdPathUri) (*dto.RecipeResBody, error)
	DeleteFavRecipe(userId uint, input *dto.CommonIdPathUri) error
	GetAllFavRecipes(input *dto.CommonQueryPage, userId uint) (*dto.CommonPageRespBody, error)
}

// recipeService is an implementation of the RecipeService interface
type recipeService struct {
	repo   repositories.RecipeRepository
	logger *zap.Logger
}

// NewRecipeService creates a new RecipeService instance
func NewRecipeService(repo repositories.RecipeRepository) RecipeService {
	return &recipeService{
		repo:   repo,
		logger: zap.L(),
	}
}

// CreateRecipe creates a new recipe with given input and user id
func (s *recipeService) CreateRecipe(input *dto.RecipeReqBody, userId uint) (*dto.RecipeResBody, error) {
	recipeModel := input.ConvertToModdel()
	recipeModel.AuthorID = uint64(userId)
	recipeModel, err := s.repo.CreateRecipe(recipeModel, input.IngredientsId, userId)
	if err != nil {
		return nil, err
	}
	recipeRes := &dto.RecipeResBody{}
	recipeRes.ConvertFromModel(recipeModel)
	return recipeRes, nil
}

// GetAllRecipes returns all the recipes using input to define pagination
func (s *recipeService) GetAllRecipes(input *dto.CommonQueryPage) (*dto.CommonPageRespBody, error) {
	recipes, totalRecipe, err := s.repo.GetAllRecipes(input.PageSize, input.PageNumber)
	if err != nil {
		return nil, err
	}
	recipesRes := make([]interface{}, len(recipes))
	for i, v := range recipes {
		res := dto.RecipeResBody{}
		res.ConvertFromModel(v)
		recipesRes[i] = res
	}
	totalNbPage := int(totalRecipe / int64(input.PageSize))
	pageRes := &dto.CommonPageRespBody{
		CommonQueryPage: *input,
		TotalNbResult:   int(totalRecipe),
		TotablNbPage:    totalNbPage,
		Items:           recipesRes,
	}
	return pageRes, nil
}

// GetRecipeById is a function that returns a recipe specified by ID from the database
func (s *recipeService) GetRecipeById(input *dto.CommonIdPathUri) (*dto.RecipeResBody, error) {
	recipe, err := s.repo.GetRecipeById(input.ID)
	if err != nil {
		return nil, err
	}
	recipeRes := &dto.RecipeResBody{}
	recipeRes.ConvertFromModel(recipe)
	return recipeRes, nil
}

// DeleteRecipeById is a function that delete a recipe specified by ID from the database
func (s *recipeService) DeleteRecipeById(input *dto.CommonIdPathUri) error {
	err := s.repo.DeleteRecipeById(input.ID)
	if err != nil {
		return err
	}
	return nil
}

// AddToFavRecipe adds a recipe to the user's favorite list
func (s *recipeService) AddToFavRecipe(userId uint, input *dto.CommonIdPathUri) (*dto.RecipeResBody, error) {
	recipeModel, err := s.repo.AddToFavRecipe(userId, input.ID)
	if err != nil {
		return nil, err
	}
	recipeRes := dto.RecipeResBody{}
	recipeRes.ConvertFromModel(recipeModel)
	return &recipeRes, nil
}

// DeleteFavRecipe deletes a recipe from the user's favorite list
func (s *recipeService) DeleteFavRecipe(userId uint, input *dto.CommonIdPathUri) error {
	err := s.repo.DeleteFavRecipe(userId, input.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllFavRecipes retrieves all the recipes from the user's favorite list
func (s *recipeService) GetAllFavRecipes(input *dto.CommonQueryPage, userId uint) (*dto.CommonPageRespBody, error) {
	recipes, totalRecipe, err := s.repo.GetAllFavRecipes(input.PageSize, input.PageNumber, userId)
	if err != nil {
		return nil, err
	}
	recipesRes := make([]interface{}, len(recipes))
	for i, v := range recipes {
		res := dto.RecipeResBody{}
		res.ConvertFromModel(v)
		recipesRes[i] = res
	}
	totalNbPage := int(totalRecipe / int64(input.PageSize))
	pageRes := &dto.CommonPageRespBody{
		CommonQueryPage: *input,
		TotalNbResult:   int(totalRecipe),
		TotablNbPage:    totalNbPage,
		Items:           recipesRes,
	}
	return pageRes, nil
}
