package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/usfca-cs490/admissions-webhook/pkg/audit"
	"github.com/usfca-cs490/admissions-webhook/pkg/dashboard"
	"github.com/usfca-cs490/admissions-webhook/pkg/kind"
	"github.com/usfca-cs490/admissions-webhook/pkg/util"
	"github.com/usfca-cs490/admissions-webhook/pkg/webhook"
)

/* Startup method */
func main() {
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
	// Reconfigure cluster flag
	reconfig_val := flag.String("reconfigure", "", "reconfigure the cluster")
	// Audit the cluster
	flag.Bool("audit", false, "audit the cluster for vulnerabilities")
	// Display logs
	flag.Bool("logstream", false, "stream webhook logs to terminal")
	// Show all pods in kind-control-plane node
	flag.Bool("pods", false, "show all pods in the kind-control-plane node")
	// Add a pod
	pod_config_path := flag.String("add", "./pkg/kind/test-pods/hello-good.yaml", "attempt to add a pod to the cluster")
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
	if util.IsFlagRaised("dashboard") {
		//util.NotYetImplemented("dashboard")
		dashboard.DashInit()
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

	// Reconfigure the cluster
	if util.IsFlagRaised("reconfigure") {
		if len(strings.TrimSpace(*reconfig_val)) != 0 {
			util.ChangeConfig(*reconfig_val, "./pkg/webhook/admission_policy.json")
		}
		// Get the webhook's pod's full name
		pods := kind.GetPodsStruct("kind-control-plane")
		hookPod := kind.FindPod(pods, "the-captains-hook")
		hookPodName := string(hookPod.Name)
		// copy the policy from local files into the webhook container
		kind.CopyPolicy(hookPodName)
		// run new policy on all pods and remove any that break security, return list of names of the evicted
		// do that by requesting a dummy pod that the webhook has special cases to handle
		kind.AddPod("./pkg/webhook/review-dummy.yaml")
		// now remove that erroneous pod
		kind.DeletePod("apps", "review-dummy")
	}

	// Audit the cluster using kubeaudit
	if util.IsFlagRaised("audit") {
		audit.Audit()
	}

	// Stream cluster logs to terminal
	if util.IsFlagRaised("logstream") {
		// Stream the logs
		kind.StreamLogs("the-captains-hook")
	}

	// Show all pods in kind-control-plane node
	if util.IsFlagRaised("pods") {
		fmt.Println(string(kind.GetPods("kind-control-plane")))
	}

	// Add a pod
	if util.IsFlagRaised("add") {
		kind.AddPod(*pod_config_path)
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
