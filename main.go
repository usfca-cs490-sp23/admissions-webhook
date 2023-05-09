package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/usfca-cs490/admissions-webhook/pkg/audit"
	"github.com/usfca-cs490/admissions-webhook/pkg/cluster"
	"github.com/usfca-cs490/admissions-webhook/pkg/dashboard"
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
	flag.Bool("reconfigure", false, "reconfigure the cluster")
	// Update severity level flag
	severity := flag.String("severity", "", "update severity level to one of following: critical, high, medium, low, negligible")
	// Audit the cluster
	flag.Bool("audit", false, "audit the cluster for vulnerabilities")
	// Display logs
	flag.Bool("logstream", false, "stream webhook logs to terminal")
	// Show all pods in kind-control-plane node
	flag.Bool("pods", false, "show all pods in the kind-control-plane node")
	// Add a pod
	pod_config_path := flag.String("add", "./pkg/cluster/test-pods/hello-good.yaml", "attempt to add a pod to the cluster")
	// Shutdown flag
	flag.Bool("status", false, "print out description of webhook pod")
	// Build hook flag, should only be called by Docker container
	flag.Bool("shutdown", false, "shutdown the cluster")

	// Check for flags
	flag.Parse()

	// Create cluster with argued name
	if util.IsFlagRaised("create") {
		cluster.CreateCluster()
	}

	// Show kind info command output
	if util.IsFlagRaised("info") {
		cluster.Info()
	}

	// Launch kind cluster interface
	if util.IsFlagRaised("dashboard") {
		dashboard.DashInit()
	}

	// Apply webhook to cluster
	if util.IsFlagRaised("deploy") {
		// Build and load docker image
		cluster.BuildLoadHookImage("the-captains-hook", "latest", ".")
		// Gen certs
		cluster.GenCerts()
		// Apply configs
		cluster.ApplyConfig("./pkg/cluster-config")
		cluster.ApplyConfig("./pkg/webhook/deploy-rules")
	}

	// Change severity level
	if util.IsFlagRaised("severity") {
		// If argued level is valid, reconfigure the severity in admission_policy file
		if _, ok := cluster.SeverityLvls[strings.ToLower(*severity)]; ok {
			cluster.ChangeConfig(strings.ToUpper(string((*severity)[0]))+string((*severity)[1:]), "./pkg/webhook/admission_policy.json")
		} else {
			fmt.Println("severity: invalid param \"" + *severity + "\". Please select either critical, high, medium, low, or negligible.")
		}
	}

	// Reconfigure the cluster
	if util.IsFlagRaised("reconfigure") {
		// Get the webhook's pod's full name
		pods := cluster.GetPodsStruct("kind-control-plane")
		hookPod := cluster.FindPod(pods, "the-captains-hook")
		hookPodName := string(hookPod.Name)
		// copy the policy from local files into the webhook container
		cluster.CopyPolicy(hookPodName)
		// run new policy on all pods and remove any that break security, return list of names of the evicted
		// do that by requesting a dummy pod that the webhook has special cases to handle
		cluster.AddPod("./pkg/webhook/review-dummy.yaml")
		// now remove that erroneous pod
		cluster.DeletePod("apps", "review-dummy")
	}

	// Audit the cluster using kubeaudit
	if util.IsFlagRaised("audit") {
		audit.Audit()
	}

	// Stream cluster logs to terminal
	if util.IsFlagRaised("logstream") {
		// Stream the logs
		cluster.StreamLogs("the-captains-hook")
	}

	// Show all pods in kind-control-plane node
	if util.IsFlagRaised("pods") {
		fmt.Println(string(cluster.GetPods("kind-control-plane")))
	}

	// Add a pod
	if util.IsFlagRaised("add") {
		cluster.AddPod(*pod_config_path)
	}

	// Build the webhook
	if util.IsFlagRaised("hook") {
		webhook.Build()
	}

	// Print out the hook pod's status
	if util.IsFlagRaised("status") {
		cluster.DescribeHook("the-captains-hook")
	}

	// Delete a cluster with argued name
	if util.IsFlagRaised("shutdown") {
		cluster.Shutdown()
	}
}
