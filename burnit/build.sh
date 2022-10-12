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
      PLATFORM=$1
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

platformParts=($(echo $PLATFORM | sed "s/\// /"))
os=${platformParts[0]}
arch=${platformParts[1]}

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

if [[ $CONTAINER -eq 1 && "$os" == "linux" ]]; then
  mkdir build
  cp -r ../common build
  docker build -t $BIN:$VERSION --build-arg OS=$os --build-arg ARCH=$arch --build-arg BIN=$BIN --build-arg VERSION=$VERSION --platform $os/$arch .
  docker image prune -f
  rm -rf build
else
  bin_full_name=$BIN-$VERSION-$os-$arch
  bin_path=release/$os/bin
  mkdir -p release
  GOOS=$os GOARCH=$arch go build -o $bin_path/$bin_full_name -ldflags="-w -s" -trimpath main.go
fi
