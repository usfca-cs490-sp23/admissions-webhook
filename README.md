# Admissions Webhook
CS 490 Senior Team Project

## Installation

1. Clone the repo

2. (Optional) Test the installation

`go test ./lib/tests -v`

## Usage

This interface currently only supports reading a basic config for launching a cluster. The only supported option is `name`. For example, config.txt may look like this:

```
name test-cluster
```

To create a cluster, simply run:

`go run main.go -cluster [name]`


To create a cluster with additional configuration, simply run:

`go run main.go -c [config file]`


To shutdown the cluster, run:

`go run main.go shutdown [name]`

