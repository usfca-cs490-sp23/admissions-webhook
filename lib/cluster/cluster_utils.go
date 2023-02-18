package cluster

import (
    "os"
    "os/exec"
    "log"
    "github.com/usfca-cs490/admissions-webhook/lib/util"
)

/* Method to check kind cluster info */
func Info() {
    // Create command
    cmd := exec.Command("kubectl", "cluster-info", "--context", "kind-kind")
    // Redirect stdout and stderr to default for OS
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Run and handle errors
    if err := cmd.Run(); err != nil {
        log.Fatal("startup.go: FAILED TO GET CLUSTER INFO")
    }
}

/* Method to open Kubernetes web interface */
func Interface_() {
    util.NotYetImplemented("interface")
}

/* Method to deploy admissions control webhook */
func Deploy() {
    util.NotYetImplemented("deploy")
}

