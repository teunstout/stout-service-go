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

```bash
curl -i -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"golang","password":"golang"}'
```

#### Restricted call

```bash
curl -i -X GET http://localhost:8080/v1/endpoints/restricted/example \
  -H "Content-Type: application/json" \
  -H "Authorization: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzgyMDEzNjYsImlhdCI6MTc3ODExNDk2NiwiaXNzIjoibG9naW4tc2VydmljZSIsImp0aSI6IjJlOTg0YTUyLTg4MjQtNDQ5OS05ZjMwLWRjN2I4OTZkOTJhMCIsInN1YiI6MX0.dWgN8FOientvbUpqHQ5fY1zIZzLGOimO5w2jzvD6SelR0Bkhb0MId98LWIrRp6nPaeWUe1fLbJop5JkLBVMtWUOExcqV35gvKdYrkLj0OT1aslnrVKWQqq0AlbmlnQFYWyXo-Psmmy7Ks8aRpt4IVWgdkvrqu8TxTpu4hfoai5lBXpltw9yT05ppjUWoub2ZWltEW8evD0WsTkYe_kxgRikvqCHFW2ONSh7OWznqL5y2RUYIPfHkbVlVpcOH7ICnOoiwmL6y7at3vXMVVl9WR4TdbHGG9X_woj_omlA34f68QQIfNEusZPMdR-HUYSan24U0HcIrPXUak-I2EM8fm9frG5g-llt1ZZwfHdA_YGdoBKx6eVpXHDnzVUJbTygi8VayYI2GStcTwCxS4TnshZ_DyjKGe2XSS3atfNJNkalagamyUYZRuNkEClQuoxx1xLROIIIC4jwjvULZYZgdLdF9Zm5evhQgg6s0iSiR6oetve4vIF0YxTTN1cE5rpdh" \
  -d '{"username":"golang","password":"golang"}'
```

#### JWKS

```bash
curl -i -X GET http://localhost:8080/.well-known/jwks.json
```

## Reference

<https://auth0.com/docs/secure/application-credentials/generate-rsa-key-pair>

This package is structures following:

- <https://github.com/golang-standards/project-layout/blob/master/README.md>
- <https://go.dev/doc/modules/layout>
