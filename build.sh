#!/bin/sh
#https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
#
# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/s3master s3master.go

env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/s3master main.go

aws s3 cp ./build/s3master s3://eu-west-1-convergeassets/asset-json/cmd/
