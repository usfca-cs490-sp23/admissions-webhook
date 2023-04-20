# Admissions Webhook

A CLI tool for creating a [Kind](https://kind.sigs.k8s.io/) cluster and deploying a Webhook to it. The purpose of the webhook is to generate a Software Bill of Materials (SBOM) for each image uploaded to the cluster, and then checking it for critical Common Vulnerabilities and Exposures (CVEs). If there are any CVEs that could endanger the cluster, the image will not be allowed in. 

## Installation

### Dependencies

This project requires the following dependencies:

- [Kubernetes](https://kubernetes.io/releases/download/)

- [Kind](https://kind.sigs.k8s.io/)

- [Kubectl](https://kubernetes.io/docs/tasks/tools/)

- [Docker](https://www.docker.com/products/docker-desktop/)

- [Syft](https://github.com/anchore/syft)

- [Grype](https://github.com/anchore/grype)

- [Openssl](https://github.com/openssl/openssl)

- Clipboard Tool
    - xclip (Linux)
    - pbcopy (macOS)
    - clip (Windows)

### Installing

1. Clone the repo (that's it!)

2. (Optional) Test the installation

`go test ./tests/install/ -v`


## Usage

This interface provides a wrapper for kind, so that it is easy to create and destroy a cluster. This can be done with the `-create` flag, so the following will create a cluster that the hook can easily access:

`go run main.go -create`

From here, applying the webhook is as easy as running the following:

`go run main.go -deploy`

These can be combined to create the cluster and then deploy the hook to it with: 

`go run main.go -create -deploy`

Shutting down the cluster can be done with another wrapper method:

`go run main.go -shutdown`

### Further Functionalities

| Flag           | Arg (blank if bool)| Description                 |
|----------------|--------------------|-----------------------------|
| `-h`           |                    | get list of flags |
| `-add`         | path to yaml file  | attempt to add a pod to the cluster (default "./pkg/kind/test-pods/hello-good.yaml") |
| `-audit`       |                    | audit the cluster for vulnerabilities |
| `-create`      |                    | create a kind cluster |
| `-dashboard`   |                    | launch cluster dashboard  |
| `-deploy`      |                    | apply admissions webhook to cluster |
| `-info`        |                    | get cluster info |
| `-logstream`   |                    | stream webhook logs to terminal |
| `-pods`        |                    | show all pods in the kind-control-plane node |
| `-reconfigure` |                    | reconfigure the cluster |
| `-shutdown`    |                    | shutdown the cluster |
| `-status`      |                    | print out description of webhook pod |

To add a pod is a simple wrapper on the standard applying functionality:

`go run main.go -add <path-to-.yaml-file>`

The wehbook can be tested with three test pods stored in `pkg/kind/test-pods/` by running:

`go test ./tests/webhook/ -v`

The webhook can be audited with 

Starting up the K8's dashboard is done with: `go run main.go -dashboard` and then following the steps given once it is started

A full list of functionalities can be seen by running `go run main.go -h`

## Structure

```
.
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── pkg
│   ├── audit
│   │   ├── auditor.go
│   │   └── auditor.yaml
│   ├── cluster-config
│   │   ├── app.ns.yaml
│   │   ├── validating.config.template.yaml
│   │   └── validating.config.yaml
│   ├── dashboard
│   │   ├── admin-rb.yaml
│   │   ├── dashboard-adminuser.yaml
│   │   └── dashboard.go
│   ├── evals
│   │   └── test.json
│   ├── kind
│   │   ├── cluster_utils.go
│   │   ├── kind.cluster.yaml
│   │   └── test-pods
│   │       ├── alpine-good.yaml
│   │       ├── hello-good.yaml
│   │       └── nginx-fail.yaml
│   ├── sboms
│   │   └── test.json
│   ├── tls
│   │   └── gen_certs.sh
│   ├── util
│   │   └── util.go
│   └── webhook
│       ├── admission_policy.json
│       ├── build.go
│       ├── database
│       │   ├── redis-config.yaml
│       │   ├── redis-pod.yaml
│       │   └── redis-service-config.yaml
│       ├── deploy-rules
│       │   ├── webhook.deploy.yaml
│       │   ├── webhook.svc.yaml
│       │   └── webhook.tls.secret.yaml
│       ├── policy_accepted_parameters.txt
│       └── validate.go
├── README.md
└── tests
    ├── install
    │   └── startup_test.go
    └── webhook
        └── webhook_test.go

16 directories, 34 files
```
