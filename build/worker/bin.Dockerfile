FROM golang:1.10-stretch

RUN apt-get clean && apt-get update && apt-get install -y \
    ghostscript \
    libmagickwand-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Setup go directory
ENV GOPATH=/go

COPY tmp/worker $GOPATH/bin/worker
ENTRYPOINT $GOPATH/bin/worker
