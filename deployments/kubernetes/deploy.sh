#!/bin/bash

NAMESPACE=burnit

for arg in "$@"
do
  case $arg in
    --burnitdb-config)
      shift
      BURNITDB_CONFIG=$1
      shift
      ;;
    --burnitgw-config)
      shift
      BURNITGW_CONFIG=$1
      shift
      ;;
  esac
done

if [[ -z $BURNITDB_CONFIG || -z $BURNITGW_CONFIG ]]; then
  echo "config file path for both burnitdb and burnitgw must be provided"
  exit 1
fi

MANIFESTS=`dirname "$0"`

kubectl create namespace $NAMESPACE

kubectl create secret generic burnitdb-config \
  --from-file=$BURNITDB_CONFIG \
  --namespace $NAMESPACE

kubectl create secret generic burnitgw-config \
  --from-file=$BURNITGW_CONFIG \
  --namespace $NAMESPACE

kubectl apply -f $MANIFESTS/burnitgen/deployment.yaml --namespace $NAMESPACE
kubectl apply -f $MANIFESTS/burnitgen/service.yaml --namespace $NAMESPACE

kubectl apply -f $MANIFESTS/burnitdb/deployment.yaml --namespace $NAMESPACE
kubectl apply -f $MANIFESTS/burnitdb/service.yaml --namespace $NAMESPACE

kubectl apply -f $MANIFESTS/burnitgw/deployment.yaml --namespace $NAMESPACE
kubectl apply -f $MANIFESTS/burnitgw/service.yaml --namespace $NAMESPACE
