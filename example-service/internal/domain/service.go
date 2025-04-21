package domain

import "fmt"

// DomainService represents a service that encapsulates domain logic.
type DomainService struct{}

// PerformBusinessLogic performs a business operation.
func (s *DomainService) PerformBusinessLogic(input string) string {
	return "Processed: " + input
}

// Example usage of DomainService.
func ExampleDomainService() {
	service := &DomainService{}
	result := service.PerformBusinessLogic("input data")
	fmt.Println(result)
}
