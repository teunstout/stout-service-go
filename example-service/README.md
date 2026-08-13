> This is a scaffold/template for starting a new Go service in this style, not a running service — it isn't built by `docker-compose.yaml` and its `main.go` just prints a hello-world message. Copy this folder as a starting point rather than trying to curl it.

# Folder structure

```bash
myapp/
│
├── cmd/ # Main applications of this project
│ └── myapp/ # `main` package for the app
│ └── main.go
│
├── internal/ # Private application and library code
│ ├── auth/ # Example module (internal logic, not imported externally)
│ │ └── service.go
│ ├── user/ # Another internal package
│ │ └── handler.go
│ ├── domain/ # Domain entities and value objects
│ ├── usecase/ # Application use cases
│ ├── handler/ # REST API handlers
│ │ └── example_handler.go
│ ├── repository/ # Database interaction
│ │ └── example_repo.go
│ ├── client/ # External service clients
│ │ └── example_client.go
│ └── config/ # Configuration loading, env handling
│ └── config.go
│
├── pkg/ # Public libraries (can be imported by other projects)
│ └── logger/ # Custom logger package
│ └── logger.go
│
├── api/ # API contracts: Protobuf, OpenAPI, GraphQL schemas
│ └── openapi.yaml
│
├── web/ # Web UI, static assets, templates, frontend code
│ └── static/
│
├── scripts/ # Helpful dev scripts (build, test, etc.)
│ └── build.sh
│
├── migrations/ # DB migrations (sql or tools like goose/migrate)
│ └── 001_init.sql
│
├── deployments/ # Kubernetes, Docker Compose, Terraform, etc.
│ └── docker/
│ └── Dockerfile
│
├── test/ # Integration or e2e tests
│ └── user_test.go
│
├── go.mod
├── go.sum
└── README.md
```

## A Bit More Detail

cmd/myapp/ – The entry point of your app. You can have multiple apps (e.g., CLI tool, server).

internal/ – Encapsulation: other projects can’t import these. Good place for domain logic, services, etc.

pkg/ – Reusable utilities or packages you want to be importable.

api/ – Use it for gRPC protobufs, OpenAPI specs, or GraphQL schemas.

deployments/ – Helm charts, K8s manifests, Terraform, etc.

test/ – External, black-box tests or integration tests.

migrations/ – If using SQL DBs, keep schema changes here.

# Domain-Driven Design (DDD) in Example Service

## Overview
Domain-Driven Design (DDD) is a software design approach that emphasizes collaboration between technical and domain experts to model complex business domains. It focuses on creating a shared understanding of the domain and aligning the software design with the business needs.

### Key Concepts in DDD
1. **Entities**: Objects with a unique identity that persists over time.
2. **Value Objects**: Immutable objects that describe aspects of the domain without a unique identity.
3. **Aggregates**: A cluster of domain objects treated as a single unit for data changes.
4. **Repositories**: Abstractions for accessing and persisting aggregates.
5. **Domain Events**: Events that represent significant occurrences in the domain.
6. **Bounded Contexts**: Logical boundaries within the domain where a particular model is defined and applicable.

## Differences Between DDD and Clean Architecture
While both DDD and Clean Architecture aim to create maintainable and scalable systems, they differ in focus:

| Aspect                | Domain-Driven Design (DDD)                     | Clean Architecture                          |
|-----------------------|-----------------------------------------------|--------------------------------------------|
| **Focus**            | Domain modeling and business logic alignment | Separation of concerns and independence    |
| **Core Components**  | Entities, Value Objects, Aggregates, etc.     | Layers: Use Cases, Adapters, Entities      |
| **Boundaries**       | Bounded Contexts                              | Layered boundaries                         |
| **Implementation**   | Domain-centric                                | Structure-centric                          |

## How Components Work Together
In DDD, components are designed to work together in a way that reflects the business domain:

1. **Domain Layer**:
   - Contains the core business logic, including entities, value objects, and domain events.
   - This layer is independent of external systems like databases or APIs.

2. **Application Layer (Use Cases)**:
   - Coordinates the application logic by invoking domain logic and interacting with repositories.
   - Acts as a bridge between the domain layer and the infrastructure layer.

3. **Infrastructure Layer**:
   - Implements the technical details, such as database access, external API calls, and message queues.
   - Provides implementations for the abstractions defined in the domain layer.

4. **Presentation Layer (Handlers)**:
   - Handles user input (e.g., HTTP requests) and delegates the processing to the application layer.
   - Converts the results into a format suitable for the user (e.g., JSON responses).

### Example Workflow
1. A user sends an HTTP request to the application.
2. The handler in the presentation layer validates the request and calls the appropriate use case.
3. The use case interacts with the domain layer to execute business logic.
4. The use case uses repositories to persist or retrieve data.
5. The handler returns the response to the user.

## Additional Go Examples for DDD

### Domain Event Example
```go
// File: internal/domain/events.go
package domain

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
```

### Aggregate Example
```go
// File: internal/domain/aggregate.go
package domain

type Aggregate struct {
	ID      string
	Entities []Entity
}

type Entity struct {
	ID   string
	Name string
}

func NewAggregate(id string, entities []Entity) *Aggregate {
	return &Aggregate{
		ID:      id,
		Entities: entities,
	}
}
```

### Value Object Example
```go
// File: internal/domain/value_object.go
package domain

type ValueObject struct {
	Value string
}

func NewValueObject(value string) *ValueObject {
	return &ValueObject{Value: value}
}
```

### Service Example
```go
// File: internal/domain/service.go
package domain

type DomainService struct {}

func (s *DomainService) PerformBusinessLogic(input string) string {
	// Example business logic
	return "Processed: " + input
}
```

## Conclusion
By adopting DDD, the `example-service` can better align with business needs, improve maintainability, and handle complexity effectively. The examples provided demonstrate how to implement key DDD concepts in Go.

### Example Files

#### REST API Example
```go
// File: internal/handler/example_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"example-service/internal/domain"
	"example-service/internal/usecase"
)

type ExampleHandler struct {
	usecase *usecase.ExampleUsecase
}

func NewExampleHandler(u *usecase.ExampleUsecase) *ExampleHandler {
	return &ExampleHandler{usecase: u}
}

func (h *ExampleHandler) HandleExample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input domain.ExampleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.usecase.ProcessExample(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
```

#### Database Interaction Example
```go
// File: internal/repository/example_repo.go
package repository

import (
	"context"
	"example-service/internal/domain"
	"github.com/jackc/pgx/v4"
)

type ExampleRepository struct {
	conn *pgx.Conn
}

func NewExampleRepository(conn *pgx.Conn) *ExampleRepository {
	return &ExampleRepository{conn: conn}
}

func (r *ExampleRepository) SaveExample(ctx context.Context, example domain.Example) error {
	_, err := r.conn.Exec(ctx, "INSERT INTO examples (id, name) VALUES ($1, $2)", example.ID, example.Name)
	return err
}
```

#### Client Call Example
```go
// File: internal/client/example_client.go
package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ExampleClient struct {
	baseURL string
}

func NewExampleClient(baseURL string) *ExampleClient {
	return &ExampleClient{baseURL: baseURL}
}

func (c *ExampleClient) CallExternalService(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.baseURL+"/external", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
```

### Comments
- **REST API**: The handler decouples HTTP logic from business logic by delegating to a use case.
- **Database**: The repository abstracts database operations, ensuring the domain model remains unaffected by persistence details.
- **Client Call**: The client encapsulates external service calls, promoting reusability and testability.
