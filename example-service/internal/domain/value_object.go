package domain

import "fmt"

type ValueObject struct {
	Value string
}

func NewValueObject(value string) *ValueObject {
	return &ValueObject{Value: value}
}

func ExampleValueObject() {
	vo := NewValueObject("example")
	fmt.Printf("ValueObject: %s\n", vo.Value)
}
