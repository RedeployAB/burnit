#!/bin/bash
bin=burnit
module_path=github.com/RedeployAB/$bin
build_root=build
platform=linux/amd64

version=""
image=0
archive=0
cleanup=1
skip_tests=0

for arg in "$@"
do
  case $arg in
    --version)
      shift
      version=$1
      shift
      ;;
    --platform)
      shift
      platform=$1
      shift
      ;;
    --image)
      image=1
      shift
      ;;
    --archive)
      archive=1
      shift
      ;;
    --no-cleanup)
      cleanup=0
      shift
      ;;
    --skip-tests)
      skip_tests=1
      shift
      ;;
  esac
done

cwd=$(pwd)

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

if [ $skip_tests -eq 0 ]; then
  go test ./...
  if [ $? -ne 0 ]; then
    exit 1
  fi
fi

platforms=(${platform//,/ })
for p in "${platforms[@]}"; do
  os_arch=(${p//\// })
  os=${os_arch[0]}
  arch=${os_arch[1]}

  bin_path=$build_root/$os/$arch
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build \
  -o $bin_path/$bin \
  -ldflags="-s -w -X '$module_path/internal/version.version=$version'" \
  -trimpath main.go

  if [ $archive -eq 1 ]; then
    cd $bin_path
    targz=$bin-$version-$os-$arch.tar.gz
    cp $cwd/README.md $cwd/LICENSE $cwd/LICENSE-THIRD-PARTY $cwd/NOTICE .
    tar -czf $targz *
    mv $targz $cwd/$build_root/$targz
    cd $cwd
  fi
done

if [ $image -eq 1 ]; then
  docker buildx create --name multiarch --use --bootstrap
  docker build -t $bin:$version --platform $platform .
  docker buildx rm multiarch
fi
