package audit

import (
    "os"
    "time"
    "os/exec"
	"github.com/usfca-cs490/admissions-webhook/pkg/util"
	"github.com/usfca-cs490/admissions-webhook/pkg/kind"
)

// Audit audits a cluster by sending in a kubeaudit image
func Audit() {
    // Create the command to add the auditor pod to the cluster
    kind.AddPod("./pkg/audit/auditor.yaml")

    // Wait for five seconds for the container to create
    time.Sleep(5 * time.Second)

    // Get the kubeaudit pod's full name
    pods := kind.GetPodsStruct("kind-control-plane")
    kubeaudit_pod := kind.FindPod(pods, "kubeaudit")

    // Create the command to add the auditor pod to the cluster
    cmd := exec.Command("kubectl", "logs", string(kubeaudit_pod.Name))
    
    // Redirect output to stdout
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)

    // Delete all kubeaudit pods so that it can be rerun later
    kind.DeleteItem("pod", string(kubeaudit_pod.Name))
    kind.DeleteItem("serviceaccount", "kubeaudit")
    kind.DeleteItem("job.batch", "kubeaudit")
}

