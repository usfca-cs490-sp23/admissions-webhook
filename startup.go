package main

import (
    "fmt"
    "os"
    "os/exec"
    "bufio"
    "strings"
    "github.com/usfca-cs490/admissions-webhook/lib/util"
    "log"
)

/* Method to parse a startup file and execute with parameters */
func startup (path string) {
    // Open up the file and store it
    file, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
    }

    // Close the file
    defer file.Close()

    // Initialize parameter fields
    var name string

    // Parse the file for valid fields
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fields := strings.Split(scanner.Text(), " ")
        for i, ele := range fields {
            switch ele {
            case "name":    // name flag
                name = fields[i+1]
            }
        }
    }

    // Call createCluster on a new thread
    createCluster(name)
}

/* Method to create a cluster using kind */
func createCluster(name string) {
    // Create command
    cmd := exec.Command("kind", "create", "cluster", "--name", name)
    // Redirect stdout and stderr to default for OS
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Run and handle errors
    if err := cmd.Run(); err != nil {
        log.Fatal("startup.go: FAILED TO CREATE CLUSTER")
    }
}

/* Method to check kind cluster info */
func info() {
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
func interface_() {
    util.NotYetImplemented("interface")
}

/* Method to deploy admissions control webhook */
func deploy() {
    util.NotYetImplemented("deploy")
}

func shutdown(name string) {
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

/* Main method */
func main() {
    // Usage check
    if len(os.Args) < 2 {
        util.IncorrectUsage()
    }

    // Check second command line argument
    switch os.Args[1] {
    case "from":
        if len(os.Args) < 3 {
            util.IncorrectUsage()
        } else {
            startup(os.Args[2])
        }
    case "cluster":
        util.NotYetImplemented("create cluster")
    case "info":
        info()
    case "interface":
        interface_()
    case "deploy":
        deploy()
    case "shutdown":
        if len(os.Args) < 3 {
            util.IncorrectUsage()
        } else {
            shutdown(os.Args[2])
        }
    default:
        util.IncorrectUsage()
    }
}

