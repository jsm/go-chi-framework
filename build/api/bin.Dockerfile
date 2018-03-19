FROM golang:1.10-alpine

# Update system dependencies
RUN apk update --no-cache

# Setup go directory
ENV GOPATH=/go

COPY tmp/api $GOPATH/bin/api
ENTRYPOINT $GOPATH/bin/api
