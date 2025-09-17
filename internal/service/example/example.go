package example

import (
	"context"
	"errors"
	"fmt"
	"go-template/internal/repository/example"
)

// Service represents handler for repo used
type service struct {
	exampleRepository example.Repository
}

// Service is the interface for exampleService
type Service interface {
	ExampleGet(ctx context.Context) ([]*example.Example, error)
	ExampleAdd(ctx context.Context, name []string) error
}

// NewExampleService initiate Example Service
func NewExampleService(exampleRepository example.Repository) Service {
	return &service{
		exampleRepository: exampleRepository,
	}
}

// ExampleGet will get all records from db
func (s *service) ExampleGet(ctx context.Context) ([]*example.Example, error) {
	res, err := s.exampleRepository.GetRecords(ctx, []int{})
	if err != nil {
		return nil, errors.New("internal server error")
	}
	for _, val := range res {
		val.Name = "a"
	}
	return res, nil
}

// ExampleAdd will add records based on []string provided
func (s *service) ExampleAdd(ctx context.Context, name []string) error {
	// validate for empty array
	if len(name) == 0 {
		return fmt.Errorf("name array is empty")
	}
	// process name into struct
	payload := []*example.Example{}
	// index, value | _ means that we doesn't need index
	for _, value := range name {
		payload = append(payload, &example.Example{
			Name: value,
		})
	}
	return s.exampleRepository.AddRecords(ctx, payload)
}
