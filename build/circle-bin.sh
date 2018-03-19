#! /bin/bash

set -ex

SERVICE=$1
BUILD_IMAGE_TAG=$SERVICE-build:$CIRCLE_SHA1
BIN_IMAGE_TAG=$SERVICE:$CIRCLE_SHA1

docker build --rm=false -f build/$SERVICE/Dockerfile -t $BUILD_IMAGE_TAG .
docker run -v `pwd`/tmp:/hosttmp --entrypoint=cp $BUILD_IMAGE_TAG /go/bin/$SERVICE /hosttmp/$SERVICE
docker build --rm=false -f build/$SERVICE/bin.Dockerfile -t $BIN_IMAGE_TAG .
