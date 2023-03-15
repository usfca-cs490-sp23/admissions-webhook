# Admissions Webhook
CS 490 Senior Team Project

## Installation

1. Clone the repo

2. (Optional) Test the installation

`go test ./tests/ -v`


## Usage

This interface provides a wrapper for kind, so that it is easy to create and destroy a cluster. This can be done with the `-create` flag, so the following will create a cluster that the hook can easily access:

`go run main.go -create`

From here, applying the webhook is as easy as running the following:

`go run main.go -deploy`

These can be combined to create the cluster and then deploy the hook to it with `go run main.go -create -deploy`

Shutting down the cluster can be done with another wrapper method:

`go run main.go -shutdown`

A full list of functionalities can be seen by running `go run main.go -h`

## Structure

```
.
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── pkg
│   ├── cluster-config
│   │   ├── app.ns.yaml
│   │   ├── validating.config.template.yaml
│   │   └── validating.config.yaml
│   ├── dashboard
│   │   └── dashboard.go
│   ├── kind
│   │   ├── cluster_utils.go
│   │   └── kind.cluster.yaml
│   ├── tls
│   │   └── gen_certs.sh
│   ├── util
│   │   └── util.go
│   └── webhook
│       ├── build.go
│       ├── deploy-rules
│       │   ├── webhook.deploy.yaml
│       │   ├── webhook.svc.yaml
│       │   └── webhook.tls.secret.yaml
│       └── validate.go
├── README.md
└── tests
    └── startup_test.go
```
