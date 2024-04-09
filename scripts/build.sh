#!/usr/bin/env bash

GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap ../main.go

