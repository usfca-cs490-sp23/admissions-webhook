package dashboard

import (
    "os/exec"

	"github.com/usfca-cs490/admissions-webhook/pkg/util"
)

// DashboardUpdate TODO: Return a DashboardUpdate struct with the result of checking the internals of the pod
type DashboardUpdate struct {
	// TODO: make a list of all SBOMs from the pod to add to DB / check with grype
	// TODO:
	Denied bool
}

// DashInit initiates the dashboard on user's local computer
func DashInit() {
	cmd := exec.Command("kubectl","apply", "-f", "https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml")
	// Run and handle errors
	err := cmd.Run()
	util.FatalErrorCheck(err)

	DashUser()

	proxy := exec.Command("kubectl", "proxy")
	// Run and handle errors
	proxyErr := proxy.Run()
	util.FatalErrorCheck(proxyErr)
	//go to http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ to access
}

//DashUser creates a new user account via K8's Service Account mechanism
func DashUser(/*manFile string*/){
	//this does not work
	//cmd := exec.Command("kubectl", "apply", "-f", manFile)
	defCmd := exec.Command("kubectl", "create", "serviceaccount", "dashboard", "-n", "default")
	defErr := defCmd.Run()
	util.FatalErrorCheck(defErr)

	//pathCmd := exec.Command("$(kubectl", "get", "serviceaccount", "dashboard", "-o", "jsonpath=\"{secrets[0].name}\")")
	tknCmd := exec.Command("kubectl", "get", "secret", 
	"$(kubectl", "get", "serviceaccount", "dashboard", "-o", "jsonpath=\"{secrets[0].name}\")", "-o", "jsonpath=\"{.data.token}\"", "|", "base64", "--decode")

	tknErr := tknCmd.Run()
	util.FatalErrorCheck(tknErr)

}

// BadPodDashUpdate TODO: expand this to have field values expressing that the pod could not be examined
func BadPodDashUpdate() DashboardUpdate {
	return DashboardUpdate{Denied: true}
}
