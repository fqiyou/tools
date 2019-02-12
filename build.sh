#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/run/newsearn_v1 src/run/newsearn/main.go 

