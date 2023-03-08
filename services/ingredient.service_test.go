package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/clementb49/welsh_academy/dto"
	"github.com/clementb49/welsh_academy/models"
	"github.com/clementb49/welsh_academy/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockIngredientRepository struct{}

func (m *mockIngredientRepository) CreateIngredient(input *models.Ingredient) (*models.Ingredient, error) {
	if input.Name == "not_exist_ingredient" {
		input.ID = 1
		input.CreatedAt = time.Now()
		input.UpdatedAt = time.Now()
		return input, nil
	}
	return nil, gorm.ErrDuplicatedKey
}

func (m *mockIngredientRepository) GetAllIngredients(pageSize int, pageNumber int) ([]*models.Ingredient, int64, error) {
	if pageNumber == 1 {
		ingredients := make([]*models.Ingredient, 5)
		for i := 0; i < 5; i++ {
			ingredients[i] = &models.Ingredient{
				Model: gorm.Model{
					ID: uint(i),
				},
				Name: fmt.Sprintf("ingredient_%d", i),
				Type: "test",
			}
		}
		return ingredients, 5, nil
	}
	return make([]*models.Ingredient, 0), 5, nil
}

func (m *mockIngredientRepository) GetIngredientById(ingredientId uint) (*models.Ingredient, error) {
	if ingredientId == 1 {
		return &models.Ingredient{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name: "test",
			Type: "test_type",
		}, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockIngredientRepository) DeleteIngredientById(ingredientId uint) error {
	if ingredientId == 1 {
		return nil
	}
	return gorm.ErrRecordNotFound
}
func TestCreateIngredient(t *testing.T) {
	repo := &mockIngredientRepository{}
	ingredientService := services.NewIngredientService(repo)
	// test happy path
	input := dto.IngredientReqBody{
		Name: "not_exist_ingredient",
		Type: "test_type",
	}
	ingredientRes, err := ingredientService.CreateIngredient(&input)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), ingredientRes.ID)
	assert.Equal(t, "not_exist_ingredient", ingredientRes.Name)
	assert.Equal(t, "test_type", ingredientRes.Type)

	// test error: ingredient exists
	input.Name = "exist_ingredient"
	ingredientRes, err = ingredientService.CreateIngredient(&input)
	assert.ErrorAs(t, err, &gorm.ErrDuplicatedKey)
	assert.Nil(t, ingredientRes)
}

func TestGetIngredientById(t *testing.T) {
	repo := &mockIngredientRepository{}
	ingredientService := services.NewIngredientService(repo)

	// test happy path
	ingredienRes, err := ingredientService.GetIngredientById(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), ingredienRes.ID)
	assert.Equal(t, "test", ingredienRes.Name)
	assert.Equal(t, "test_type", ingredienRes.Type)
	// test: record not found error
	ingredienRes, err = ingredientService.GetIngredientById(2)
	assert.ErrorAs(t, err, &gorm.ErrRecordNotFound)
	assert.Nil(t, ingredienRes)
}

func TestDeleteIngrdientById(t *testing.T) {
	repo := &mockIngredientRepository{}
	ingredientService := services.NewIngredientService(repo)
	// test happy path
	err := ingredientService.DeleteIngredientById(1)
	assert.NoError(t, err)
	// test record not found
	err = ingredientService.DeleteIngredientById(2)
	assert.ErrorAs(t, err, &gorm.ErrRecordNotFound)
}

func TestGetAllIngredients(t *testing.T) {
	repo := &mockIngredientRepository{}
	ingredientService := services.NewIngredientService(repo)
	// test happy path
	input := &dto.CommonQueryPage{
		PageSize:   10,
		PageNumber: 1,
	}
	ingredientsRes, err := ingredientService.GetAllIngredients(input)
	assert.NoError(t, err)
	assert.Equal(t, 1, ingredientsRes.PageNumber)
	assert.Equal(t, 10, ingredientsRes.PageSize)
	assert.Equal(t, 5, ingredientsRes.TotalNbResult)
	assert.Equal(t, 1, ingredientsRes.TotablNbPage)
	assert.Equal(t, 5, len(ingredientsRes.Items))
	for i := 0; i < len(ingredientsRes.Items); i++ {
		ingredient, ok := ingredientsRes.Items[i].(dto.IngredientResBody)
		assert.Equal(t, true, ok)
		assert.Equal(t, uint(i), ingredient.ID)
	}
	// test page empty
	input.PageNumber = 2
	ingredientsRes, err = ingredientService.GetAllIngredients(input)
	assert.NoError(t, err)
	assert.Equal(t, 1, ingredientsRes.TotablNbPage)
	assert.Equal(t, 5, ingredientsRes.TotalNbResult)
	assert.Equal(t, 0, len(ingredientsRes.Items))
}
