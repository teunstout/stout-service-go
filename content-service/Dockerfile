FROM golang:1.24 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o content-service

FROM golang:1.24

EXPOSE 8080

WORKDIR /app
COPY --from=builder /build/content-service .
COPY --from=builder /build/internal/app/assets ./assets

USER 1000:1000

CMD ["./content-service"]