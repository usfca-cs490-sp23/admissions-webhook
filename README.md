# Admissions Webhook

An ecosystem that uses a Command Line Interface to create a [Kind](https://kind.sigs.k8s.io/) cluster, deploy a Webhook to it, and provide easy access to the Kubernetes Dashboard. The ecosystem has a number of security auditing and reporting features built in. The purpose of the webhook is to generate a Software Bill of Materials (SBOM) for each image uploaded to the cluster, and then checking it for Common Vulnerabilities and Exposures (CVEs). If there any CVEs found that are rated at a severity at or higher than the tolerance specified in the admission policy, the image will not be allowed in. 

## Installation

### Dependencies

This project requires the following dependencies:

- [Kubernetes](https://kubernetes.io/releases/download/)

- [Kind](https://kind.sigs.k8s.io/)

- [Kubectl](https://kubernetes.io/docs/tasks/tools/)

- [Docker](https://www.docker.com/products/docker-desktop/)

- [Openssl](https://github.com/openssl/openssl)

- Clipboard Tool
    - xclip (Linux)
    - pbcopy (macOS)
    - clip (Windows)

NOTE: An internet connection is required to run the webhook

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

To add a pod is a simple wrapper on the standard applying functionality:

`go run main.go -add <path-to-.yaml-file>`


### Further Functionalities

| Flag           | Arg (blank if bool)   | Description                 |
|----------------|-----------------------|-----------------------------|
| `-h`           |                       | get list of flags |
| `-add`         | path to yaml file     | attempt to add a pod to the cluster |
| `-audit`       |                       | audit the cluster for vulnerabilities |
| `-create`      |                       | create a kind cluster |
| `-dashboard`   |                       | launch cluster dashboard  |
| `-deploy`      |                       | apply admissions webhook to cluster |
| `-info`        |                       | get cluster info |
| `-logstream`   |                       | stream webhook logs to terminal |
| `-pods`        |                       | show all pods in the kind-control-plane node |
| `-reconfigure` |                       | reconfigure the cluster |
| `-severity`    | level                 | update severity level to one of following: critical, high, medium, low, negligible |
| `-shutdown`    |                       | shutdown the cluster |
| `-status`      |                       | print out description of webhook pod |


The wehbook can be tested with three test pods stored in `pkg/cluster/test-pods/` by running:

`go test ./tests/webhook/ -v`

The webhook can be audited with 

Starting up the K8's dashboard is done with: `go run main.go -dashboard` and then following the steps provided in the terminal. NOTE: The dashboard can take up to 30 seconds to generate, so be patient if the page does not load immediately. If nothing pops up automatically, try refreshing.

A full list of functionalities can be seen by running `go run main.go -h`


## Configuring Admission Policy

In order to quickly reconfigure the admission severity level, the `-severity [level]` flag fan be used, and then the `-reconfigure` flag can push the new level to the webhook. The severity options are as follows:

- Critical
- High
- Medium
- Low
- Negligible

The two flags can be run in conjunction like this: `go run main.go -severity High -reconfigure`.

NOTE: The severity flag will only accept valid levels, and is not case-sensitive, although modifying the actual admission policy file stored in `./pkg/webhook/admission_policy.json` has no such validity or case checks, so if you must modify this file directly, be sure to follow the detailed instructions included in `./pkg/webhook/policy_accepted_parameters.txt`.

If needed, CVEs can be added to the whitelist in the admission policy file, but be sure to consult the instructions in `./pkg/webhook/policy_accepted_parameters.txt` before modifying the file.


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
│   ├── cluster
│   │   ├── cluster_utils.go
│   │   ├── kind.cluster.yaml
│   │   └── test-pods
│   │       ├── alpine-good.yaml
│   │       ├── hello-good.yaml
│   │       └── nginx-fail.yaml
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
│       ├── review-dummy.yaml
│       └── validate.go
├── README.md
└── tests
    ├── install
    │   └── startup_test.go
    └── webhook
        └── webhook_test.go

16 directories, 35 files
```
