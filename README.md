# Admissions Webhook
CS 490 Senior Team Project

## Usage

This interface currently only supports reading a basic config for launching a cluster. The only supported option is `name`. For example, config.txt may look like this:

```
name test-cluster
```

To create a cluster, simply run:

`go run startup.go from [config file]`


To shutdown the cluster, run:

`go run startup.go shutdown [name]`
