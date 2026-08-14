FROM golang:1.26 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o idp-service ./cmd

FROM golang:1.26

EXPOSE 8080

RUN apt-get update && apt-get install -y openssl curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/idp-service .

USER 1000:1000

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./idp-service"]