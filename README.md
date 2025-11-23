# nkp-apple

A command-line tool to create and manage Kubernetes NKP (Nutanix Kubernetes Platform) clusters using Apple container runtime on macOS.

## Overview

`nkp-apple` is a wrapper around the NKP CLI that enables you to create and manage Kubernetes clusters on macOS using Apple's native container runtime. It automatically handles the creation of bootstrap clusters needed for self-managed Kubernetes deployments.

## Prerequisites

- macOS with Apple container runtime installed
- [Apple container runtime](https://github.com/apple/container) - Install and ensure it's running
- NKP CLI (`nkp`) installed and available in your PATH
- Go 1.25.4 or later (for building from source)

## Installation

There are several installation options:

- As Homebrew or Linuxbrew package
- Manual installation

After installing, the tool will be available as `nkp-apple`.

### Homebrew Package

You can install with [Homebrew](https://brew.sh)

```sh
brew install funkolab/tap/nkp-apple
```

Keep up-to-date with `brew upgrade nkp-apple` (or `brew upgrade` to upgrade everything)

### Manual

 - Download your corresponding [release](https://github.com/funkolab/nkp-apple/releases)
 - Install the binary somewhere in your PATH (`/usr/local/bin` for example)
 - use it with `nkp-apple`

***MacOS X notes for security error***

 Depending on your OS settings, when installing the binary manually you must run the following command:
 `xattr -r -d com.apple.quarantine /usr/local/bin/nkp-apple`

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
- Export kubeconfig to `~/.kube/config`

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

Other flags are passed directly with the nkp command in the background.


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
nkp-apple delete cluster -c [cluster-name] --kubeconfig xxx [--self-managed]
```

## Configuration

The tool uses the following default configurations:

- **Node Name**: `konvoy-capi-bootstrapper-control-plane`
- **Node Image**: `docker.io/mesosphere/konvoy-bootstrap:nkp-version`
- **Pod CIDR**: `10.244.0.0/16`
- **Memory**: 8GB allocated to bootstrap container
- **API Server Port**: 6443 (mapped to localhost)

## Architecture

The tool wraps the NKP CLI and manages the lifecycle of a local Kubernetes bootstrap cluster using Apple's container runtime. The bootstrap cluster is used by NKP to provision and manage target Kubernetes clusters on various platforms.

## Contributing

Contributions and feedback are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on GitHub.
