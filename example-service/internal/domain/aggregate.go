package domain

import "fmt"

type Aggregate struct {
	ID       string
	Entities []Entity
}

type Entity struct {
	ID   string
	Name string
}

func NewAggregate(id string, entities []Entity) *Aggregate {
	return &Aggregate{
		ID:       id,
		Entities: entities,
	}
}

func ExampleAggregate() {
	entities := []Entity{
		{ID: "1", Name: "Entity1"},
		{ID: "2", Name: "Entity2"},
	}
	aggregate := NewAggregate("123", entities)
	fmt.Printf("Aggregate ID: %s, Entities: %v\n", aggregate.ID, aggregate.Entities)
}
