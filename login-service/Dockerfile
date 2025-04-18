FROM golang:1.24 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o login-service

FROM golang:1.24

EXPOSE 8080

RUN apt-get update && apt-get install -y openssl && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/login-service .

USER 1000:1000

CMD ["./login-service"]