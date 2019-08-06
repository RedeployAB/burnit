#!/bin/bash

if [ -z "$1" ]; then
    echo "Specify a version as first argument, as 'x.x.x'."
    exit 1
fi

VERSION=$1
OS=linux
# Test and build with make.
make VERSION=$VERSION release -j3 || { echo 'FAILURE: test/build failed'; exit 1;}
# Rename and copy file
mkdir -p release/bin
cp release/secretgen-$VERSION-$OS-amd64 release/bin/secretgen
# Create docker image.
docker build -t secretgen:$VERSION --build-arg VERSION=$VERSION .
