package domain

import "fmt"

// Aggregate represents a cluster of domain objects treated as a single unit.
type Aggregate struct {
	ID       string
	Entities []Entity
}

// Entity represents a part of the aggregate.
type Entity struct {
	ID   string
	Name string
}

// NewAggregate creates a new aggregate.
func NewAggregate(id string, entities []Entity) *Aggregate {
	return &Aggregate{
		ID:       id,
		Entities: entities,
	}
}

// Example usage of Aggregate.
func ExampleAggregate() {
	entities := []Entity{
		{ID: "1", Name: "Entity1"},
		{ID: "2", Name: "Entity2"},
	}
	aggregate := NewAggregate("123", entities)
	fmt.Printf("Aggregate ID: %s, Entities: %v\n", aggregate.ID, aggregate.Entities)
}
