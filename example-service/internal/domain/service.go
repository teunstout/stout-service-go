package domain

import "fmt"

type DomainService struct{}

func (s *DomainService) PerformBusinessLogic(input string) string {
	return "Processed: " + input
}

func ExampleDomainService() {
	service := &DomainService{}
	result := service.PerformBusinessLogic("input data")
	fmt.Println(result)
}
