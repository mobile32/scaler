#!/bin/bash

cd src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o main main.go && zip deployment.zip main