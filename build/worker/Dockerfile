FROM golang:1.10-stretch

RUN apt-get clean && apt-get update && apt-get install -y \
    ghostscript \
    git \
    libmagickwand-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Setup go directory
ENV GOPATH=/go
ENV BACKEND_PATH=$GOPATH/src/github.com/jsm/gode
RUN mkdir -p $BACKEND_PATH/build
WORKDIR $BACKEND_PATH

# Install dep
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml $BACKEND_PATH/
RUN $GOPATH/bin/dep ensure -vendor-only

# Build worker
COPY . $BACKEND_PATH
RUN go build -ldflags "-X main.version=$(cat VERSION)" -o /$GOPATH/bin/worker worker/main.go

ENTRYPOINT $GOPATH/bin/worker
