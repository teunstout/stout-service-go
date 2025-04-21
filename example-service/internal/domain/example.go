package domain

// Example represents a domain entity with a unique identity.
type Example struct {
	ID   string
	Name string
}

// ExampleInput represents input data for processing an example.
type ExampleInput struct {
	Name string
}
