# Admissions Webhook
CS 490 Senior Team Project

## Installation

1. Clone the repo

2. (Optional) Test the installation

`go test ./lib/tests -v`

## Structure

```
.
├── README.md
├── go.mod
├── go.sum
├── main.go
├── pkg
│   ├── cluster
│   │   ├── cluster_utils.go
│   │   ├── shutdown.go
│   │   └── startup.go
│   ├── dashboard
│   │   └── dashboard.go
│   ├── tls
│   │   ├── new_cert_sript.sh
│   │   └── tls.go
│   ├── util
│   │   └── util.go
│   └── webhook
│       ├── build.go
│       ├── config
│       │   ├── app.ns.yaml
│       │   └── validating-config.yaml
│       ├── secrets
│       │   ├── cab.txt
│       │   ├── cert.txt
│       │   ├── fakeCa.txt
│       │   └── key.txt
│       └── validate.go
└── tests
    ├── startup_test.go
    └── tls_test.go
```

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

