#!/bin/bash

function help {
  echo "Usage:"
  echo "--version         Version of binary. Mandatory"
  echo "--docker          Build docker image with binary"
  echo "-p, --platform    Platform for binary (os/arch). Supported: linux|darwin/amd64|arm64. Default: linux/amd64"
  echo "-h, --help        Prints help information"
  exit 0
}

# Set defaults.
BIN=burnit
PLATFORM=linux/amd64
VERSION=""
IMAGE=0

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
      PLATFORM=$1
      shift
      ;;
    --image)
      IMAGE=1
      shift
      ;;
    esac
done

if [ -z $VERSION ]; then
  help
  exit 1
fi

platformParts=($(echo $PLATFORM | sed "s/\// /"))
os=${platformParts[0]}
arch=${platformParts[1]}
outPath=build

if [[ $os != "linux" ]] && [[ $os != "darwin" ]]; then
  echo "os: $os is not a supported operating system"
  exit 1
fi

if [ -z $arch ]; then
  arch=amd64
fi

if [[ $arch != "amd64" ]] && [[ $arch != "arm64" ]]; then
  echo "architecture: $arch is not a supported architecture"
  exit 1
fi

# Run tests.
go test ./...
if [ $? -ne 0 ]; then
  exit 1
fi

CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o $outPath/$BIN -ldflags="-w -s" -trimpath

if [[ $IMAGE -eq 1 && "$os" == "linux" ]]; then
  docker build -t $BIN:$VERSION --build-arg BIN=$BIN --platform $os/$arch .
  docker image prune -f
fi
