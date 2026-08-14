# Jisho Service

The Jisho Service is a microservice designed to interact with the Jisho API, a popular online Japanese-English dictionary. This service provides a clean and efficient way to query the Jisho API and retrieve dictionary data, making it a valuable tool for applications that require Japanese language support.

Every request must carry a JWT issued by [idp-service](../idp-service/README.MD). The token is verified locally against a shared RSA public key (see [Quickstart](#quickstart)) — this service never calls back into idp-service to check it.

## Quickstart

```bash
cd jisho-service

# Windows
$env:CONNECTION_STRING="user=golang password=golang host=localhost port=5432 dbname=production sslmode=disable"; go run ./cmd/main.go

# Linux
CONNECTION_STRING="user=golang password=golang host=localhost port=5432 dbname=production sslmode=disable" go run ./cmd/main.go
```

If `CONNECTION_STRING` isn't set, it falls back to a hardcoded local default (`dbname=postgres`) — set it explicitly to match whichever database `init.sql` actually ran against (`production` in the root `docker-compose.yaml`).

The service also needs the RSA **public** key used to verify JWTs signed by idp-service: it reads `/app/keys/app.rsa.pub` (as mounted in `docker-compose.yaml` from the `rsa-key-generator` init container), falling back to `../app.rsa.pub` (repo root) when run locally outside Docker. The service listens on `http://localhost:8080`.

### Via the nginx gateway

When running the full stack with `docker compose up`, `jisho-service` isn't exposed on the host directly — reach it through nginx on `http://localhost:8080`, prefixed with `/jisho-service` (e.g. `http://localhost:8080/jisho-service/v1/search`). The curl below uses the direct-to-service URL, i.e. running the service standalone as shown above.

### Curl's

These are curls used for most of the endpoints. Please be aware that these might not be up to date after updating usecases in the app.

#### Search

Requires a `Bearer` JWT from idp-service's `/v1/login` (its `sub` claim is used as the member id when recording search history). Returns `400` if `keyword` is missing/empty, and `401` if the token is missing, malformed, or fails RS256 verification.

```bash
curl -i -X GET "http://localhost:8080/v1/search?keyword=犬" \
  -H "Authorization: Bearer <jwt-from-idp-login>"
```

Response body is the raw [Jisho API](https://jisho.org/) response, re-shaped as:

```json
{
  "meta": { "status": 200 },
  "data": [
    {
      "slug": "犬",
      "is_common": true,
      "tags": [],
      "jlpt": ["jlpt-n5"],
      "japanese": [{ "word": "犬", "reading": "いぬ" }],
      "senses": [{ "english_definitions": ["dog"], "parts_of_speech": ["Noun"], "tags": [], "see_also": [] }]
    }
  ]
}
```

## Clean Architecture Overview

The Jisho Service is built using the principles of Clean Architecture. Clean Architecture emphasizes separation of concerns and independence of frameworks, databases, and external services. This approach ensures that the core business logic remains independent and testable, while the infrastructure and external dependencies can be easily replaced or modified.

### Key Principles of Clean Architecture:

1. **Independence**: The core logic is independent of frameworks, databases, and external services.
2. **Testability**: The architecture is designed to make unit testing straightforward.
3. **Separation of Concerns**: Each layer has a specific responsibility, reducing coupling and increasing maintainability.

## Folder Structure and Relationships

The Jisho Service follows a modular folder structure to align with Clean Architecture principles. Below is an explanation of the folders and their relationships:

### 1. `cmd/`

- Contains the entry point of the application (`main.go`).
- Responsible for initializing the application and wiring up dependencies.

### 2. `internal/`

- The core of the service, divided into several subfolders:

#### a. `adapters/`

- Contains implementations of interfaces for external systems.
- Example: `logger/` includes `logrus_logger.go`, which provides logging functionality using the Logrus library.

#### b. `domain/`

- Defines core business entities and interfaces.
- Example: `logger.go` defines the logging interface used across the service.

#### c. `external/`

- Handles communication with external services or APIs.
- Example: `jishoClient/` includes `jisho_client.go` and `response.go` for interacting with the Jisho API.

#### d. `handler/`

- Contains HTTP handlers for processing incoming requests and returning responses.
- Example: `jisho_handler.go` handles requests related to Jisho API queries.

#### e. `middleware/`

- Implements middleware for cross-cutting concerns like authentication and logging.
- Example: `auth.go` provides authentication middleware.

#### f. `repository/`

- Manages data persistence and retrieval.
- Example: `jisho_repo.go` handles data storage and retrieval logic for Jisho-related data.

#### g. `usecase/`

- Contains the core business logic and application use cases.
- Example: `jisho_usecase.go` implements the main logic for querying the Jisho API and processing results.

## Differences with Domain-Driven Design (DDD)

While Clean Architecture and Domain-Driven Design (DDD) share some similarities, they have distinct focuses and approaches:

1. **Focus**:

   - Clean Architecture focuses on the separation of concerns and independence of the core logic from external systems.
   - DDD emphasizes modeling the domain and solving complex business problems by aligning the software design with the business domain.

2. **Core Components**:

   - Clean Architecture organizes the code into layers (e.g., use cases, adapters, and entities) to ensure independence and testability.
   - DDD introduces concepts like Aggregates, Entities, Value Objects, Repositories, and Domain Events to model the domain effectively.

3. **Boundaries**:

   - Clean Architecture defines boundaries between layers to isolate the core logic.
   - DDD defines bounded contexts to separate different parts of the domain and their respective models.

4. **Implementation**:
   - Clean Architecture is more prescriptive about the structure and flow of dependencies.
   - DDD is more focused on understanding and modeling the domain, leaving implementation details more flexible.

## How to Use

See [Quickstart](#quickstart) above for the environment variables and command to run the service, and [Curl's](#curls) for example requests.
