// Package dto defines data transfer objects (DTOs) used for communicating between the input and output of an API
package dto

import (
	"time"

	"gorm.io/gorm"
)

// CommonResBody represents the common response body that includes ID, CreatedAt, and UpdatedAt fields
type CommonResBody struct {
	ID        uint      `json:"id" xml:"id"`                 // ID represents the unique identifier of the resource
	CreatedAt time.Time `json:"created_at" xml:"created_at"` // CreatedAt represents the creation timestamp of the resource
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at"` // UpdatedAt represents the last modification timestamp of the resource
}

// convertFromGormModel is a helper method to convert from GORM model to CommonResBody
func (c *CommonResBody) convertFromGormModel(model *gorm.Model) {
	c.ID = model.ID
	c.CreatedAt = model.CreatedAt
	c.UpdatedAt = model.UpdatedAt
}

// CommonQueryPage represents the common pagination query parameters
type CommonQueryPage struct {
	PageSize   int `form:"page_size" json:"page_size" xml:"page_size" binding:"min=-1"`       // PageSize represents the number of items per page
	PageNumber int `form:"page_number" json:"page_number" xml:"page_number" binding:"min=-1"` // PageNumber represents the page number to fetch
}

// CommonPageRespBody represents the common pagination response body
type CommonPageRespBody struct {
	CommonQueryPage               // CommonQueryPage represents the common pagination query parameters
	TotalNbResult   int           `json:"total_result" xml:"total_result"` // TotalNbResult represents the total number of results
	TotablNbPage    int           `json:"total_page" xml:"total_page"`     // TotablNbPage represents the total number of pages
	Items           []interface{} // Items represents the list of items in the current page
}

// CommonIdPathUri represents the common URI parameter for the ID of the resource
type CommonIdPathUri struct {
	ID uint `uri:"id" binding:"required,min=0"` // ID represents the unique identifier of the resource
}
