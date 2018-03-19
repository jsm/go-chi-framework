FROM golang:1.10-alpine

# Install system dependencies
RUN apk add --update --no-cache \
    git

# Setup go directory
ENV GOPATH=/go
ENV BACKEND_PATH=$GOPATH/src/github.com/jsm/gode
RUN mkdir -p $BACKEND_PATH/build
WORKDIR $BACKEND_PATH

# Install dep
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml $BACKEND_PATH/
RUN $GOPATH/bin/dep ensure -vendor-only

# Build api
COPY . $BACKEND_PATH
RUN go build -ldflags "-X main.version=$(cat VERSION)" -o /$GOPATH/bin/api api/main.go api/router.go

ENTRYPOINT $GOPATH/bin/api
