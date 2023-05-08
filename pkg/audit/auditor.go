package audit

import (
	"bytes"
	"fmt"
	"github.com/usfca-cs490/admissions-webhook/pkg/kind"
	"github.com/usfca-cs490/admissions-webhook/pkg/util"
	"os"
	"os/exec"
	"strings"
	"time"
)

// filterAudit filters out audit results to only include annotations from namespace=apps
func filterAudit(results string) {
	// Split the results based on headings
	resultSplit := strings.Split(results, "---------------- Results for ---------------")

	fmt.Print("[36m---------------- Results for ---------------[0m")

	// Loop through results
	for _, note := range resultSplit[1:] {
		// Only print the result if the namespace is apps
		if strings.Contains(note, "namespace: apps") {
			fmt.Println(note)
		}
	}

	fmt.Println("\n[36m--------[ FINISHED AUDITING CLUSTER ]-------[0m")
}

// Audit audits a cluster by sending in a kubeaudit image
func Audit() {
	// Create the command to add the auditor pod to the cluster
	kind.AddPod("./pkg/audit/auditor.yaml")

	// Wait for five seconds for the container to create
	time.Sleep(5 * time.Second)

	// Get the kubeaudit pod's full name
	pods := kind.GetPodsStruct("kind-control-plane")
	kubeauditPod := kind.FindPod(pods, "kubeaudit")

	// Create the command to add the auditor pod to the cluster
	cmd := exec.Command("kubectl", "logs", string(kubeauditPod.Name))

	var results bytes.Buffer

	// Redirect output to results for filtering
	cmd.Stdout = &results
	cmd.Stderr = os.Stderr

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)

	filterAudit(results.String())
	// Delete all kubeaudit pods so that it can be rerun later
	kind.DeleteItem("pod", string(kubeauditPod.Name))
	kind.DeleteItem("serviceaccount", "kubeaudit")
	kind.DeleteItem("job.batch", "kubeaudit")
}
