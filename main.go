package main

import (
	"flag"
	"github.com/usfca-cs490/admissions-webhook/pkg/kind"
	"github.com/usfca-cs490/admissions-webhook/pkg/util"
	"github.com/usfca-cs490/admissions-webhook/pkg/webhook"
)

/* Startup method */
func main () {
	// Cluster name flag
	flag.Bool("create", false, "create a kind cluster")
	// Info flag
	flag.Bool("info", false, "get cluster info")
	// Interface flag
	flag.Bool("dashboard", false, "launch cluster dashboard")
	// Deploy webhook flag
	flag.Bool("deploy", false, "apply admissions webhook to cluster")
    // Print webhook pod status flag
    flag.Bool("hook", false, "start webhook, should only be called by Docker container")
	// Shutdown flag
    flag.Bool("status", false, "print out description of webhook pod")
    // Build hook flag, should only be called by Docker container
	flag.Bool("shutdown", false, "shutdown the cluster")

	// Check for flags
	flag.Parse()

    // Create cluster with argued name
	if util.IsFlagRaised("create") {
		kind.CreateCluster()
	}

    // Show kind info command output
    if util.IsFlagRaised("info") {
        kind.Info()
	}

    // Launch kind cluster interface
    if util.IsFlagRaised("interface") {
		util.NotYetImplemented("dashboard")
	}
    
    // Apply webhook to cluster
    if util.IsFlagRaised("deploy") {
        // Build and load docker image
        kind.BuildLoadHookImage("the-captains-hook", "latest", ".")
        // Gen certs
        kind.GenCerts()
        // Apply configs
        kind.ApplyConfig("./pkg/cluster-config")
        kind.ApplyConfig("./pkg/webhook/deploy-rules")
    }
    
    // Build the webhook
    if util.IsFlagRaised("hook") {
        webhook.Build()
	}

    // Print out the hook pod's status
    if util.IsFlagRaised("status") {
        kind.DescribeHook("the-captains-hook")
    }

    // Delete a cluster with argued name
    if util.IsFlagRaised("shutdown") {
		kind.Shutdown()
	}
}

