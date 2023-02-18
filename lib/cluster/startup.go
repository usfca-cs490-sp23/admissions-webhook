package cluster

import (
    "fmt"
    "os"
    "os/exec"
    "bufio"
    "strings"
    "log"
)

/* Method to parse a startup file and execute with parameters */
func Startup (path string) {
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
    CreateCluster(name)
}

/* Method to create a cluster using kind */
func CreateCluster(name string) {
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
