package domain

import "fmt"

// DomainEvent represents a significant occurrence in the domain.
type DomainEvent struct {
	EventType string
	Payload   interface{}
}

// NewDomainEvent creates a new domain event.
func NewDomainEvent(eventType string, payload interface{}) *DomainEvent {
	return &DomainEvent{
		EventType: eventType,
		Payload:   payload,
	}
}

// Example usage of DomainEvent.
func ExampleDomainEvent() {
	event := NewDomainEvent("UserCreated", map[string]string{"username": "johndoe"})
	fmt.Printf("Event Type: %s, Payload: %v\n", event.EventType, event.Payload)
}
