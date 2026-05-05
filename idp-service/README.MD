# IDP Service

## Quickstart

```bash
cd idp-service

# Windows
$env:CONNECTION_STRING="user=golang password=golang host=localhost port=5432 dbname=postgres sslmode=disable"; go run ./cmd/main.go

# Linux
CONNECTION_STRING="user=golang password=golang host=localhost port=5432 dbname=postgres sslmode=disable" go run ./cmd/main.go
```

### Curl's

#### Register

```bash
curl -i -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"golang","password":"golang"}'
```

#### Login

```bash
curl -i -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"golang","password":"golang"}'
```

## Reference

<https://auth0.com/docs/secure/application-credentials/generate-rsa-key-pair>

This package is structures following:

- <https://github.com/golang-standards/project-layout/blob/master/README.md>
- <https://go.dev/doc/modules/layout>
