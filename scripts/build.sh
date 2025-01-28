#!/bin/bash
bin=burnit
module_path=github.com/RedeployAB/$bin
os=linux
arch=amd64
build_root=build

version=""
container=0
archive=0
cleanup=0

for arg in "$@"
do
  case $arg in
    --version)
      shift
      version=$1
      shift
      ;;
    --os)
      shift
      os=$1
      shift
      ;;
    --arch)
      shift
      arch=$1
      shift
      ;;
    --image)
      container=1
      shift
      ;;
    --archive)
      archive=1
      shift
      ;;
    --cleanup)
      cleanup=1
      shift
      ;;
  esac
done

bin_path=$build_root/$os/$arch

if [ -z $bin ]; then
  bin=$(echo $(basename $(pwd)))
fi

if [ -z $version ]; then
  echo "A version must be specified."
  exit 1
fi

if [ -z $build_root ]; then
  exit 1
fi

if [ $cleanup -eq 1 ] && [ -d $build_root ]; then
  rm -rf $build_root/*
fi

mkdir -p $build_root

go test ./...
if [ $? -ne 0 ]; then
  exit 1
fi

CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build \
  -o $bin_path/$bin \
  -ldflags="-s -w -X '$module_path/internal/version.version=$version'" \
  -trimpath main.go

if [ $container -eq 1 ] && [ "$os" == "linux" ]; then
  docker buildx create --name multiarch --use --bootstrap
  docker build -t $bin:$version --platform $os/$arch .
  docker buildx rm multiarch
fi

if [ $archive -eq 1 ]; then
  cwd=$(pwd)
  cp LICENSE LICENSE-THIRD-PARTY.md README.md $bin_path
  cd $bin_path
  targz=$bin-$version-$os-$arch.tar.gz
  tar -czf $targs *
  mv $targz $build_root/$targz
  cd $cwd
fi
