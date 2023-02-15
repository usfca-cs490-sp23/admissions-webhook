package util

import (
    "fmt"
    "os"
)

/* Method to print out correct usage method and exit the prgram */
func IncorrectUsage() {
    // Print message to stdout
    fmt.Println("\nUsage: go run builder.go COMMAND\n\nA program to build a kubernetes",
    "cluster using kind, and then apply a security admission webhook\n\nCommands:\n",
    "\tcluster\t\tCreate a kind cluster\n\tinfo\t\tGet cluster information\n",
    "\tinterface\tOpen Kubernetes web interface\n\tdeploy\t\tDeploy admissions",
    "controller webhook")
    // Terminate program processing
    os.Exit(0)
}

/* Helper method to panic and trace to source method for unimplemented code */
func NotYetImplemented(method string) {
    panic((method + " not yet implemented"))
}

