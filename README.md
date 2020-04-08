# sqs-to-hec

[![Build Status - master](https://travis-ci.com/djschaap/sqs-to-hec.svg?branch=master)](https://travis-ci.com/djschaap/sqs-to-hec)

## Overview

sqs-to-hec reads messages from an AWS SQS queue and submits them to
Splunk HTTP Event Collector (HEC).

## Build/Release

### Compile Locally

```bash
go fmt cmd/cli/cli.go
go fmt pkg/receivesqs/receivesqs.go
go fmt pkg/sendhec/sendhec.go
go mod tidy
go test all
# commit any changes
BUILD_DT=`date +%FT%T%z`
COMMIT=`git rev-parse --short HEAD`
FULL_COMMIT=`git log -1`
VER=0.0.0
go build -ldflags \
  "-X main.build_dt=${BUILD_DT} -X main.commit=${COMMIT} -X main.version=${VER}" \
  cmd/cli/cli.go
```

### Build Container

```bash
docker build .
```

### Run Container

```bash
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
docker run -d \
  -e AWS_ACCESS_KEY_ID -e AWS_REGION -e AWS_SECRET_ACCESS_KEY \
  -e SRC_QUEUE=https://sqs.us-east-1.amazonaws.com/ACCOUNT/QUEUE_NAME \
  -e HEC_URL=https://splunk.example.com:8088 \
  -e HEC_TOKEN=00000000-0000-0000-0000-000000000000 \
  sqs-to-hec
```
