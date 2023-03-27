package testys

import (
    "testing"
    "strings"
	"github.com/usfca-cs490/admissions-webhook/pkg/kind"
)

// TestWebhook test case for two good pods and one bad one
func TestWebhook(t * testing.T) {
    // THe pods to test (name of their config file without path or extension
    test_pods := []string{"hello-good", "alpine-good", "nginx-bad"}

    // Loop through pods
    for i, test_pod := range test_pods {
        // Add the pod to the cluster
        kind.AddPod(string("../../pkg/kind/test-pods/" + test_pod + ".yaml"))

        // Get a list of all the pods in the cluster
        pods := kind.GetPods("kind-control-plane")
       
        // Check that all of the good pods are in the cluster
        if i != 2 {
            if !strings.Contains(pods, test_pod) {
                t.Errorf("%s should have been admitted but was not!", test_pod)
            }
        } else { // Check that the bad pod is not in the cluster
            if strings.Contains(pods, test_pod) {
                t.Errorf("%s should not have been admitted but was!", test_pod)
            }
        }
    }
}
