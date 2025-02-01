#!/bin/bash

cd ./stout-content-go
go build .
nohup ./content &

cd ../stout-idp-go
go build .
nohup ./idp &