// The package 'services' contains the business logic for handling route
package services

import (
	"math"

	"github.com/clementb49/welsh_academy/dto"
	"github.com/clementb49/welsh_academy/repositories"
	"go.uber.org/zap"
)

// IngredientService is an interface that defines the functions that the ingredient service should implement
type IngredientService interface {
	CreateIngredient(input *dto.IngredientReqBody) (*dto.IngredientResBody, error)
	GetAllIngredients(input *dto.CommonQueryPage) (*dto.CommonPageRespBody, error)
	GetIngredientById(id uint) (*dto.IngredientResBody, error)
	DeleteIngredientById(id uint) error
}

// ingredientService is a struct that implements the IngredientService interface
type ingredientService struct {
	repo   repositories.IngredientRepository
	logger *zap.Logger
}

// NewIngredientService is a function that returns a new instance of IngredientService
func NewIngredientService(repo repositories.IngredientRepository) IngredientService {
	return &ingredientService{
		repo:   repo,
		logger: zap.L(),
	}
}

// CreateIngredient is a function that creates a new ingredient in the database
func (s *ingredientService) CreateIngredient(input *dto.IngredientReqBody) (*dto.IngredientResBody, error) {
	ingredientModel := input.ConvertToModel()
	ingredientModel, err := s.repo.CreateIngredient(ingredientModel)
	if err != nil {
		return nil, err
	}
	ingredientRes := &dto.IngredientResBody{}
	ingredientRes.ConvertFromModel(ingredientModel)
	return ingredientRes, nil
}

// GetAllIngredients is a function that returns a page of ingredients from the database
func (s *ingredientService) GetAllIngredients(input *dto.CommonQueryPage) (*dto.CommonPageRespBody, error) {
	ingredients, totalIngredient, err := s.repo.GetAllIngredients(input.PageSize, input.PageNumber)
	if err != nil {
		return nil, err
	}
	ingredientsRes := make([]interface{}, len(ingredients))
	for i, ing := range ingredients {
		res := dto.IngredientResBody{}
		res.ConvertFromModel(ing)
		ingredientsRes[i] = res
	}
	totalNbPage := int(math.Ceil(float64(totalIngredient) / float64(input.PageSize)))
	pageRes := &dto.CommonPageRespBody{
		CommonQueryPage: *input,
		TotalNbResult:   int(totalIngredient),
		TotablNbPage:    totalNbPage,
		Items:           ingredientsRes,
	}
	return pageRes, nil
}

// GetIngredientById is a function that returns an ingredient specified by ID from the database
func (s *ingredientService) GetIngredientById(id uint) (*dto.IngredientResBody, error) {
	ingredient, err := s.repo.GetIngredientById(id)
	if err != nil {
		return nil, err
	}
	ingredientRes := &dto.IngredientResBody{}
	ingredientRes.ConvertFromModel(ingredient)
	return ingredientRes, nil
}

// DeleteIngredientById is a function that delete an ingredients specified by ID from the database
func (s *ingredientService) DeleteIngredientById(id uint) error {
	err := s.repo.DeleteIngredientById(id)
	if err != nil {
		return err
	}
	return nil
}
