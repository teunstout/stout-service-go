package usecase

import (
	"context"
	"example-service/internal/domain"
	"example-service/internal/repository"
)

type ExampleUsecase struct {
	repo *repository.ExampleRepository
}

func NewExampleUsecase(repo *repository.ExampleRepository) *ExampleUsecase {
	return &ExampleUsecase{repo: repo}
}

func (u *ExampleUsecase) ProcessExample(input domain.ExampleInput) (*domain.Example, error) {
	example := &domain.Example{
		ID:   generateID(), // Assume generateID is a helper function
		Name: input.Name,
	}

	if err := u.repo.SaveExample(context.Background(), *example); err != nil {
		return nil, err
	}

	return example, nil
}
