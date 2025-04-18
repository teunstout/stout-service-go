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
