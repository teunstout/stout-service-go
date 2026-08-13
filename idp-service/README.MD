# IDP Service

This idp service is used for my service. The idea of go is to write it as explicit as possible but also because I like the go language.

## Quickstart

To start the service please use the below commands. Sometimes windows blocks a program, you can go to your windows defender, apps and browser control, smart app control settings, and turn off enhanced protection.

```bash
cd idp-service

# Windows
$env:CONNECTION_STRING="user=golang password=golang host=localhost port=5432 dbname=postgres sslmode=disable"; go run ./cmd/main.go

# Linux
CONNECTION_STRING="user=golang password=golang host=localhost port=5432 dbname=postgres sslmode=disable" go run ./cmd/main.go
```

### Curl's

These are curls used for most of the endpoints. Please be aware that these might not be up to date after updating usecases in the app.

#### Register

```bash
curl -i -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"golang","password":"golang"}'
```

#### Login

Responds with `Set-Cookie: session_token=...` (HttpOnly) and a JSON body of `{"jwt": "...", "csrfToken": "..."}`.

```bash
curl -i -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"golang","password":"golang"}'
```

#### Authenticate (session check)

Validates the `session_token` cookie against the `x-csrf-token` header, as set/returned by `/v1/login`. Intended to be called by other services (e.g. nginx `auth_request`) to check whether a browser session is still valid.

```bash
curl -i -X GET http://localhost:8080/v1/auth \
  -H "Cookie: session_token=<session-token-from-login>" \
  -H "x-csrf-token: <csrf-token-from-login>"
```

#### Logout

Deletes the session/CSRF tokens and clears the `session_token` cookie. Add `?all_devices=true` to log out every session for the account.

```bash
curl -i -X GET http://localhost:8080/v1/logout \
  -H "Cookie: session_token=<session-token-from-login>" \
  -H "x-csrf-token: <csrf-token-from-login>"
```

#### Restricted call

The JWT returned by `/v1/login` must be sent as a `Bearer` token.

```bash
curl -i -X GET http://localhost:8080/v1/endpoints/restricted/example \
  -H "Authorization: Bearer <jwt-from-login>"
```

#### JWKS

```bash
curl -i -X GET http://localhost:8080/.well-known/jwks.json
```

#### Health

```bash
curl -i -X GET http://localhost:8080/health
```

## Reference

<https://auth0.com/docs/secure/application-credentials/generate-rsa-key-pair>

This package is structures following:

- <https://github.com/golang-standards/project-layout/blob/master/README.md>
- <https://go.dev/doc/modules/layout>
