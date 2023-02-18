package cluster

import (
    "os"
    "os/exec"
    "log"
)

// Method to shutdown a cluster
func Shutdown(name string) {
    // Create command
    cmd := exec.Command("kind", "delete", "cluster", "--name", name)
    // Redirect stdout and stderr to default for OS
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Run and handle errors
    if err := cmd.Run(); err != nil {
        log.Fatal("startup.go: FAILED TO GET CLUSTER INFO")
    }
}
