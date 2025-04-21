package domain

import "fmt"

// ValueObject represents an immutable object in the domain.
type ValueObject struct {
	Value string
}

// NewValueObject creates a new value object.
func NewValueObject(value string) *ValueObject {
	return &ValueObject{Value: value}
}

// Example usage of ValueObject.
func ExampleValueObject() {
	vo := NewValueObject("example")
	fmt.Printf("ValueObject: %s\n", vo.Value)
}
