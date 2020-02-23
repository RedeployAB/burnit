#!/bin/bash

if [ -z "$1" ]; then
    echo "Specify a version as first argument, as 'x.x.x'."
    exit 1
fi

VERSION=$1
OS=linux
BIN=burnitdb
# Test and build with make.
make VERSION=$VERSION release -j3 || { echo 'FAILURE: test/build failed'; exit 1;}
# Rename and copy file
mkdir -p release/bin
cp release/$BIN-$VERSION-$OS-amd64 release/bin/$BIN
# Create docker image.
docker build -t $BIN:$VERSION --build-arg VERSION=$VERSION .
# Clean up dangling images.
echo ""
echo "Deletibng dangling images..."
docker image prune -f
echo ""
echo "Build complete."
