# sqs-to-hec

[![Build Status - master](https://travis-ci.com/djschaap/sqs-to-hec.svg?branch=master)](https://travis-ci.com/djschaap/sqs-to-hec)

Project Home: https://github.com/djschaap/sqs-to-hec

Docker Hub: https://cloud.docker.com/repository/docker/djschaap/sqs-to-hec

## Overview

sqs-to-hec reads messages from an AWS SQS queue and submits them to
Splunk HTTP Event Collector (HEC).

## Build/Release

### Compile Locally

```bash
go fmt cmd/cli/main.go
go fmt pkg/receivesqs/receivesqs.go
go fmt pkg/sendhec/sendhec.go
go mod tidy
go test ./...
# commit any changes
BUILD_DT=`date +%FT%T%z`
COMMIT=`git rev-parse --short HEAD`
FULL_COMMIT=`git log -1`
VER=0.0.0
go build -ldflags \
  "-X main.build_dt=${BUILD_DT} -X main.commit=${COMMIT} -X main.version=${VER}" \
  cmd/cli/main.go
```

### Build Container (Manually)

```bash
docker build -t sqs-to-hec .
```

### Publish Container (Manually)

```bash
CONTAINER_REGISTRY=ACCOUNT.dkr.ecr.REGION.amazonaws.com
VER=0.0.0
# If AWS ECR, log in to registry every 12 hours:
# aws ecr get-login-password \
#  | docker login -u AWS --password-stdin ${CONTAINER_REGISTRY}
docker tag sqs-to-hec ${CONTAINER_REGISTRY}/sqs-to-hec:${VER}
docker push ${CONTAINER_REGISTRY}/sqs-to-hec:${VER}
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
