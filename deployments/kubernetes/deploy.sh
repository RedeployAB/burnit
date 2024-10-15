#!/bin/bash

NAMESPACE=burnit

for arg in "$@"
do
  case $arg in
    --burnit-config)
      shift
      BURNIT_CONFIG=$1
      shift
      ;;
  esac
done

if [[ -z $BURNIT_CONFIG ]]; then
  echo "config file path for burnit must be provided"
  exit 1
fi

MANIFESTS=`dirname "$0"`

kubectl create namespace $NAMESPACE

kubectl create secret generic burnit-config \
  --from-file=$BURNIT_CONFIG \
  --namespace $NAMESPACE

kubectl apply -f $MANIFESTS/burnit/deployment.yaml --namespace $NAMESPACE
kubectl apply -f $MANIFESTS/burnit/service.yaml --namespace $NAMESPACE
