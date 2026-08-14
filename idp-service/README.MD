# IDP Service

This idp service is used for my service. The idea of go is to write it as explicit as possible but also because I like the go language.

It owns accounts and sessions, issues signed JWTs (RS256) on login, and exposes a JWKS endpoint so other services can verify those JWTs without calling back into this service (see [jisho-service](../jisho-service/README.MD), which does exactly that).

## Quickstart

To start the service please use the below commands. Sometimes windows blocks a program, you can go to your windows defender, apps and browser control, smart app control settings, and turn off enhanced protection.

The `dbname` below matches `POSTGRES_DB` in the root `docker-compose.yaml` (`production`) — if you point at a database `init.sql` hasn't run against, the `accounts`/`session_tokens`/`csrf_tokens` tables won't exist yet.

```bash
cd idp-service

# Windows
$env:CONNECTION_STRING="user=golang password=golang host=localhost port=5432 dbname=production sslmode=disable"; go run ./cmd/main.go

# Linux
CONNECTION_STRING="user=golang password=golang host=localhost port=5432 dbname=production sslmode=disable" go run ./cmd/main.go
```

The service reads its RSA keypair from `/app/keys/app.rsa` / `/app/keys/app.rsa.pub` (as mounted in `docker-compose.yaml` from the `rsa-key-generator` init container), falling back to `../app.rsa` / `../app.rsa.pub` (repo root) when run locally outside Docker.

### Via the nginx gateway

When running the full stack with `docker compose up`, `idp-service` isn't exposed on the host directly — reach it through nginx on `http://localhost:8080`, prefixed with `/idp-service` (e.g. `http://localhost:8080/idp-service/v1/login`). The curls below use direct-to-service URLs (`http://localhost:8080/v1/...`), i.e. running the service standalone as shown above.

### Curl's

These are curls used for most of the endpoints. Please be aware that these might not be up to date after updating usecases in the app.

#### Register

Returns `201` on success, `400` for a missing/invalid body, or `409` (`{"error":"Account already exists"}`) if the username is taken. This is how accounts get created — there is no signup UI in the frontend by design, so creating an account means running this curl (or inserting into `accounts` directly).

```bash
curl -i -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"golang","password":"golang"}'
```

Via Docker Compose (through nginx, prefixed with `/idp-service` — see [above](#via-the-nginx-gateway)):

```bash
curl -i -X POST http://localhost:8080/idp-service/v1/register \
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
