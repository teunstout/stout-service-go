#!/bin/bash

docker run \
    --network teunstout \
    -e CONNECTION_STRING='user=golang password=golang host=http://postgres port=5432 dbname=development sslmode=disable' \
    login-service


cd ./authorization-service
go build .
docker build -t authorization-service .

cd ./content-service
go build .
docker build -t content-service .

cd ./jisho-service
go build .
docker build -t jisho-service .

cd ./login-service
go build .
docker build -t jisho-service .

cd ./nginx
go build .
docker build -t nginx .

docker compose up
