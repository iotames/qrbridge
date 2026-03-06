#!/bin/bash

git pull
go mod tidy

BUILD_TIME=$(date +%Y-%m-%d_%H_%M)
# APP_VERSION=$(cat version.txt)
read APP_VERSION < version.txt
echo "($APP_VERSION)"
CGO_ENABLED=0 go build -v -o oaplus -trimpath -ldflags "-X 'main.BuildTime=${BUILD_TIME}' -X 'main.Version=${APP_VERSION}' -X 'main.DbFlag=false' " .
