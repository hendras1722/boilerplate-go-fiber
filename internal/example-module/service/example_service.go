package service

import (
	"github.com/username/msa-boilerplate-go/internal/example-module/dto"
	exampleRepo "github.com/username/msa-boilerplate-go/internal/example-module/repository"
)

type ExampleService interface {
	GetExampleMessage() (*dto.ExampleResponse, error)
}

type exampleService struct {
	repo exampleRepo.ExampleRepository
}

func NewExampleService(repo exampleRepo.ExampleRepository) ExampleService {
	return &exampleService{
		repo: repo,
	}
}

func (s *exampleService) GetExampleMessage() (*dto.ExampleResponse, error) {
	return &dto.ExampleResponse{
		Message: "Hello World from Boilerplate",
	}, nil
}
