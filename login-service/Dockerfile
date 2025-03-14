FROM golang:1.24 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o login-service

FROM golang:1.24

RUN apt-get update && apt-get install -y openssl && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/login-service .

RUN openssl genrsa -out app.rsa 3072 && \
    openssl rsa -in app.rsa -pubout > app.rsa.pub && \
    chown 1000:1000 app.rsa app.rsa.pub

USER 1000:1000

CMD ["./login-service"]