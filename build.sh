#!/bin/bash

BUILD_TIME=$(date +%Y-%m-%d_%H_%M)
APP_VERSION=v1.11.2
CGO_ENABLED=0 go build -v -o oaplus -trimpath -ldflags "-X 'main.BuildTime=${BUILD_TIME}' -X 'main.Version=${APP_VERSION}' -X 'main.DbFlag=false' " .
