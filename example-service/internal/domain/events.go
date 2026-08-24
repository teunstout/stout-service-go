package domain

import "fmt"

type DomainEvent struct {
	EventType string
	Payload   interface{}
}

func NewDomainEvent(eventType string, payload interface{}) *DomainEvent {
	return &DomainEvent{
		EventType: eventType,
		Payload:   payload,
	}
}

func ExampleDomainEvent() {
	event := NewDomainEvent("UserCreated", map[string]string{"username": "johndoe"})
	fmt.Printf("Event Type: %s, Payload: %v\n", event.EventType, event.Payload)
}
