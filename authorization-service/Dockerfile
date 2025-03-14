FROM golang:1.24 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o authorization-service

FROM golang:1.24

WORKDIR /app
COPY --from=builder /build/authorization-service .

USER 1000:1000

CMD ["./authorization-service"]