# Content Service

The Content Service serves static content over HTTP. Today that's a single endpoint that streams a CV/résumé PDF, but it's the place new static/content-delivery endpoints would go.

## Quickstart

```bash
cd content-service

go run .
```

The service listens on `http://localhost:8080` and expects to be run from a directory where `./assets/Teun-Johán-Stout.pdf` resolves relative to the working directory (this is why the Dockerfile copies the `assets` folder next to the binary — see below).

No environment variables are required to run this service. `docker-compose.yaml` passes it a `CONNECTION_STRING` for consistency with the other services, but the service does not currently read it or talk to Postgres.

### Docker

```bash
docker build -t content-service ./content-service
docker run -p 8080:8080 content-service
```

### Curl's

These are curls used for most of the endpoints. Please be aware that these might not be up to date after updating usecases in the app.

#### Download CV

Streams the CV as `application/pdf`. No authentication required.

```bash
curl -i -X GET http://localhost:8080/v1/content/cv --output cv.pdf
```

Through the nginx gateway (see [nginx](../nginx/README.md)):

```bash
curl -i -X GET http://localhost:8080/content-service/v1/content/cv --output cv.pdf
```

## Structure

```text
content-service/
├── main.go                      # entrypoint, calls app.NewApp()
└── internal/app/
    ├── app.go                   # route registration + handlers
    └── assets/                  # static files served by the service
```
