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
├── lib
│   ├── cluster
│   │   ├── cluster_utils.go
│   │   ├── shutdown.go
│   │   └── startup.go
│   ├── keygen
│   │   └── keygen.go
│   ├── tests
│   │   ├── startup_test.go
│   │   └── tls_test.go
│   └── util
│       └── util.go
├── main.go
└── webhook
    └── secrets
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

