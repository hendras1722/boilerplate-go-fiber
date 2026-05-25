package repository

import (
	"gorm.io/gorm"
)

type ExampleRepository interface {
	// Add repository methods here
}

type exampleRepository struct {
	db *gorm.DB
}

func NewExampleRepository(db *gorm.DB) ExampleRepository {
	return &exampleRepository{
		db: db,
	}
}
