#!/bin/bash

function help {
  echo "Usage:"
  echo "--version         Version of binary. Mandatory"
  echo "--docker          Build docker image with binary"
  echo "-p, --platform    Platform for binary. Supported: linux|darwin. Default: linux"
  echo "-h, --help        Prints help information"
  exit 0
}

# Set defaults.
BIN=burnitgw
ARCH=amd64
OS=linux
VERSION=""
CONTAINER=0

# Handle incoming parameters.
for arg in "$@"
do
  case $arg in
    -h|--help)
      help
      ;;
    --version)
      shift
      VERSION=$1
      shift
      ;;
    -p|--platform)
      shift
      if [[ $1 != "linux" && $1 != "darwin" ]]; then
        echo "$1 is not a supported platform"
        exit 1
      fi
      OS=$1
      shift
      ;;
    --docker)
      CONTAINER=1
      shift
      ;;
    esac
done

if [ -z $VERSION ]; then
  help
  exit 1
fi

# Run tests.
go test ./...
if [ $? -ne 0 ]; then
  exit 1
fi

bin_full_name=$BIN-$VERSION-$OS-$ARCH
bin_path=release/$OS/bin

mkdir -p release
GOOS=$OS GOARCH=$ARCH go build -o release/$bin_full_name -ldflags="-w -s" -trimpath cmd/$BIN/main.go

if [[ $CONTAINER -eq 1 && "$OS" == "linux" ]]; then
  mkdir -p $bin_path
  cp release/$bin_full_name $bin_path/$BIN
  docker build -t $BIN:$VERSION --build-arg VERSION=$VERSION .
  docker image prune -f
fi