FROM golang:1.24 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o jisho-service

FROM golang:1.24

EXPOSE 8080

WORKDIR /app
COPY --from=builder /build/jisho-service .

USER 1000:1000

CMD ["./jisho-service"]