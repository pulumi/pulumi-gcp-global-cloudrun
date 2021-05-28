# Global GCP Cloud Run

This repo is a [Pulumi Package](https://www.pulumi.com/docs/guides/pulumi-packages/) used to deploy a CloudRun service across
all CloudRun regions and then will configure a global Cloud HTTP Load Balancer with an anycast IP that will route users to the
nearest Cloud Run location.

It's written in Go, but thanks to Pulumi's multi language SDK generating capability, it create usable SDKs for all of Pulumi's [supported languages](https://www.pulumi.com/docs/intro/languages/)

> :warning: **This package is a work in progress**: Please do not use this in a production environment!

# Building and Installing

## Building from source

But if you need to build it yourself, just download this repository, [install](https://taskfile.dev/#/installation) [Task](https://taskfile.dev/):

```sh
go get github.com/go-task/task/v3/cmd/task
```

And run the following command to build and install the plugin in the correct folder (resolved automatically based on the current Operating System):

```sh
task install
```

## Install Plugin Binary

Before you begin, you'll need to install the latest version of the Pulumi Plugin using `pulumi plugin install`:

```
pulumi plugin install resource globalgcpcloudrun v0.0.1 --server https://stack72.jfrog.io/artifactory/pulumi-packages/pulumi-globalgcpcloudrun
```

This installs the plugin into `~/.pulumi/plugins`.

## Install your chosen SDK

Next, you need to install your desired language SDK using your languages package manager.

### Python

```
pip3 install stack72-pulumi-globalgcpcloudrun
```

### NodeJS

```
npm install @stack72/pulumi-globalgcpcloudrun
```

### DotNet

```
Coming Soon
```

### Go

```
go get -t github.com/stack72/pulumi-globalgcpcloudrun/sdk/go/globalgcpcloudrun
```

# Usage

Once you've installed all the dependencies, you can use the library like any other Pulumi SDK. See the [examples](examples/) directory for examples of how you might use it.

# Limitations

This module currently only works with HTTP addresses. It is planned to add HTTPS addresses and SSL configuration in a future release