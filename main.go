package main

import (
    "os"
    "flag"
    "github.com/usfca-cs490/admissions-webhook/lib/cluster"
    "github.com/usfca-cs490/admissions-webhook/lib/util"
)

/* Main method */
func main() {
    // Config file flag
    config_file := flag.String("c", "", "path-to-config-file")
    // Cluster name flag
    cluster_name := flag.String("cluster", "cluster", "cluster-name")
    // Info flag
    flag.Bool("info", false, "get cluster info")
    // Interface flag
    flag.Bool("interface", false, "launch cluster interface")
    // Deploy webhook flag
    flag.Bool("deploy", false, "apply admissions webhook to cluster")
    // Shutdown flag
    shutdown := flag.String("shutdown", "cluster", "name-of-cluster")

    // Check for flags
    flag.Parse()

    // Usage check
    if len(os.Args) < 2 {
        util.Usage()
    }

    // Check second command line argument
    if util.IsFlagRaised("c") {                 // Build cluster and webhook from config file
        cluster.Startup(*config_file)
    } else if util.IsFlagRaised("cluster") {    // Create cluster with argued name
        cluster.CreateCluster(*cluster_name)
    } else if util.IsFlagRaised("info") {       // Show kind info command output
       util.NotYetImplemented("info")
    } else if util.IsFlagRaised("interface") {  // Launch kind cluster interface
        util.NotYetImplemented("interface")
    } else if util.IsFlagRaised("deploy") {     // Apply webhook to cluster
        util.NotYetImplemented("deploy")
    } else if util.IsFlagRaised("shutdown") {   // Delete a cluster with argued name
        cluster.Shutdown(*shutdown)
    }
}
