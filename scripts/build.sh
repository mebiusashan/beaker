#! /bin/bash

WORKSPACE=`pwd`
cd $WORKSPACE
cd cmd

cd cli
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../beaker_cli_linux main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../../beaker_cli_mac main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../../beaker_cli_windows main.go

cd ../server 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../beaker_server main.go

cd ../admin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../beaker_admin main.go
