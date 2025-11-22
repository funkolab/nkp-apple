# nkp-apple

A command-line tool to create and manage Kubernetes NKP (Nutanix Kubernetes Platform) clusters using Apple's container runtime on macOS.

## Overview

`nkp-apple` is a wrapper around the NKP CLI that enables you to create and manage Kubernetes clusters on macOS using Apple's native container runtime. It automatically handles the creation of bootstrap clusters needed for self-managed Kubernetes deployments.

## Prerequisites

- macOS with Apple container runtime installed
- [Apple container runtime](https://github.com/apple/container) - Install and ensure it's running
- NKP CLI (`nkp`) installed and available in your PATH
- Go 1.25.4 or later (for building from source)

## Installation

### From Source

```bash
git clone https://github.com/funkolab/nkp-apple.git
cd nkp-apple
go build -o nkp-apple
sudo mv nkp-apple /usr/local/bin/
```

## Usage

### Starting Apple Container Runtime

Before using `nkp-apple`, ensure the Apple container runtime is running:

```bash
container system start
```

### Creating a Bootstrap Cluster

Create a local bootstrap cluster using Apple containers:

```bash
nkp-apple create bootstrap
```

This command will:
- Create a container named `konvoy-capi-bootstrapper-control-plane`
- Initialize a Kubernetes cluster using kubeadm
- Configure networking with CNI (pod CIDR: 10.244.0.0/16)
- Set up storage class
- Export kubeconfig to `~/.kube/konvoy-capi-bootstrapper.conf`

### Creating a Kubernetes Cluster

Create a Kubernetes cluster on supported platforms:

```bash
nkp-apple create cluster [platform] [flags]
```

Supported platforms:
- `aks` - Azure Kubernetes Service
- `aws` - Amazon Web Services
- `azure` - Azure
- `eks` - Amazon Elastic Kubernetes Service
- `gcp` - Google Cloud Platform
- `nutanix` - Nutanix
- `preprovisioned` - Pre-provisioned infrastructure
- `vsphere` - VMware vSphere

For self-managed clusters, add the `--self-managed` flag:

```bash
nkp-apple create cluster vsphere --self-managed [additional-flags]
```

This will automatically create a bootstrap cluster first, then create your target cluster.

### Deleting a Bootstrap Cluster

Remove the bootstrap cluster:

```bash
nkp-apple delete bootstrap
```

This command will:
- Stop and remove the bootstrap container
- Clean up the kubeconfig file

### Deleting a Kubernetes Cluster

Delete a Kubernetes cluster:

```bash
nkp-apple delete cluster [cluster-name] [flags]
```

## Configuration

The tool uses the following default configurations:

- **Node Name**: `konvoy-capi-bootstrapper-control-plane`
- **Node Image**: `docker.io/mesosphere/konvoy-bootstrap:v2.16.1`
- **Pod CIDR**: `10.244.0.0/16`
- **Memory**: 8GB allocated to bootstrap container
- **API Server Port**: 6443 (mapped to localhost)

## Examples

### Create a self-managed vSphere cluster

```bash
nkp-apple create cluster vsphere \
  --self-managed \
  --cluster-name=my-cluster \
  --control-plane-endpoint-host=192.168.1.100 \
  --datacenter=MyDatacenter \
  --datastore=MyDatastore \
  --folder=/MyFolder \
  --network=VM-Network \
  --resource-pool=MyResourcePool \
  --server=vcenter.example.com \
  --ssh-public-key-file=~/.ssh/id_rsa.pub \
  --template=/MyTemplates/ubuntu-2004-kube-v1.28.7 \
  --virtual-ip-interface=eth0
```

### Create and use bootstrap cluster manually

```bash
# Create bootstrap cluster
nkp-apple create bootstrap

# Use the kubeconfig
export KUBECONFIG=~/.kube/konvoy-capi-bootstrapper.conf
kubectl get nodes

# When done, clean up
nkp-apple delete bootstrap
```

## Architecture

The tool wraps the NKP CLI and manages the lifecycle of a local Kubernetes bootstrap cluster using Apple's container runtime. The bootstrap cluster is used by NKP to provision and manage target Kubernetes clusters on various platforms.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on GitHub.
