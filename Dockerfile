FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN \
  BUILD_DT=`date +%FT%T%z` \
  && COMMIT=container \
  && VER=0.0.3 \
  && echo "BUILD build_dt=${BUILD_DT}" \
  && echo "BUILD commit=${COMMIT}" \
  && echo "BUILD version=${VER}" \
  && CGO_ENABLED=0 GOOS=linux go build -ldflags \
    "-X main.build_dt=${BUILD_DT} -X main.commit=${COMMIT} -X main.version=${VER}" \
    -o main cmd/cli/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /
CMD ["/main"]
