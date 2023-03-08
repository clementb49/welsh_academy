// This package make the interface between the database and the service
package repositories

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Common structure for all repository interface type. 
type repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

