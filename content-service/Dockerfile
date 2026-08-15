FROM golang:1.26 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o content-service

FROM golang:1.26

EXPOSE 8080

WORKDIR /app
COPY --from=builder /build/content-service .

USER 1000:1000

CMD ["./content-service"]