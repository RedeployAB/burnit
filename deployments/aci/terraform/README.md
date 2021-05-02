# Azure Container Instances with Terraform

This Terraform module contains definitions for deploying the application environment
in Azure with Azure Container instances (container groups) and Azure Redis Cache.

When Azure Container instances supports proper port mapping,
the included `nginx` container will be rendered reduntant and this
module will be updated.

## Prerequisites

* `terraform` - At least `v0.14.0`
* `azure-cli` - At least `v2.14.0`

The recommended setup for this is to have a pre-existing resource group with an Azure Container registry with the `burnit` images uploaded.
In a production scenario a real TLS certificate should be used instead of this example's self signed one. Additionally a DNS record
should be registered together with the FQDN specified in the certificate. 

## Provision environnment

**Steps**:

- [Create Resource Group with Azure Container Registry](#create-resource-group-with-azure-container-registry)
- [Build and push `burnit` images](#build-and-push-burnit-images)
- [Create self-signed certificates](#create-self-signed-certificates)
- [Provision environment with Terraform](#provision-environment-with-terraform)

### Create Resource Group with Azure Container Registry
```sh
# A future version of this will include the use of a Service Principal
# for access to the Container Registry from the Container Group.
acr_rg_name=<acr-resource-group-name>
acr_name=<acr-name>
acr_url=$acr_url.azurecr.io
location=<location>

az group create --name $acr_rg_name --location $location
az acr create \
  --resource-group $acr_rg_name \
  --name $acr_name \
  --location $location \
  --sku Basic \
  --admin-enabled
```
### Build and push `burnit` images

(This will be scripted, but manual example for emphasis)

```sh
burnitgen_version=<version>
burnitdb_version=<version>
burnitgw_version=<version>

burnitgen_image=burnitgen:$burnitgen_version
burnitdb_image=burnitdb:$burnitdb_version
burnitgw_image=burnitgw:$burnitgw_version
# From the project root.
cd burnitgen
./build --version $burnitgen_version --docker
docker tag $burnitgen_image $acr_url/$burnitgen_image
docker push $acr_url/$burnitgen_image
cd ..

cd burnitdb
./build --version $burnitdb_version --docker
docker tag $burnitdb_image $acr_url/$burnitdb_image
docker push $acr_url/$burnitdb_image
cd ..

cd burnitdb
./build --version $burnitgw_version --docker
docker tag $burnitgw_image $acr_url/$burnitgw_image
docker push $acr_url/$burnitgw_image
cd ..
```

### Create self-signed certificates

As described above this is just for the example. In a production scenario a proper certificate should be used.

```sh
# The name of .key and .crt kan be whatever, the terraform definitions gives
# them the correct names inside the volume.
openssl req -x509 -nodes -newkey rsa:2048 -keyout cert.key -out cert.crt -days 3650
```

Copy these to a location where you can specify their path for the input parameters
in the Terraform steps. (Don't forget to edit .gitignore so any certificates or other
sensitive information is commited).

### Provision environment with Terraform

## Usage

Minimum `terraform.tfvars` example (fill out as necessary):
```hcl
# Resource group variables
resource_group_name   = "<resource-group-name>"

# All resources variables
location = "<location>"
tags = {}

# Redis variables
redis_cache_name = "<redis-cache-name>"

# ACR and registry variables
acr_name                = "<acr-name>"
acr_resource_group_name = "<acr-resource-group-name>"

# Service variables
ssl_certificate_path = "<path-to-cert>"
ssl_key_path         = "<path-to-key>"
nginx_config_path      = "./conf/nginx.conf"
```

Or without `.tfvars`:

```
terraform plan -out=deploy.plan \
  -var resource_group_name="<resource-group-name>" \
  -var location="<location>" \
  -var redis_cache_name="<redis-cache-name>" \
  -var acr_name="<acr-name>" \
  -var acr_resource_group_name="<acr-resource-group-name>" \
  -var ssl_certificate_path="<path-to-cert>" \
  -var ssl_key_path="<path-to-key>" \
  -var nginx_config_path="./conf/nginx.conf"
```

For more available settings end methods of deployment see: `variables.tf`.
