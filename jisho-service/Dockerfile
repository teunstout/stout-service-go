FROM golang:1.24 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o jisho-service ./cmd

FROM golang:1.24

EXPOSE 8080

WORKDIR /app
COPY --from=builder /build/jisho-service .

# Ensure the keys are accessible in the container
USER 1000:1000

CMD ["./jisho-service"]