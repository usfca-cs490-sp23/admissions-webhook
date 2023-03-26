package dashboard

import (
	"os"
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
	cmd := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml")
	// Run and handle errors
	err := cmd.Run()
	util.FatalErrorCheck(err)

	print("go to: http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/\n")

	DashUser("./pkg/dashboard/dashboard-adminuser.yaml", "./pkg/dashboard/admin-rb.yaml")

	cmd = exec.Command("kubectl", "proxy")
	// Run and handle errors
	err = cmd.Run()
	util.FatalErrorCheck(err)
	//go to http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ to access
}

// DashUser creates a new user account via K8's Service Account mechanism
func DashUser(adminUser string, adminRb string) {
	//create admin service account (see dashboard-adminuser.yaml)
	cmd := exec.Command("kubectl", "apply", "-f", adminUser)
	err := cmd.Run()
	util.FatalErrorCheck(err)

	//create cluster role binding (see admin-rb.yaml)
	cmd = exec.Command("kubectl", "apply", "-f", adminRb)
	err = cmd.Run()
	util.FatalErrorCheck(err)

	//name of the service account
	saName := "admin-user"

	tkn := exec.Command("kubectl", "-n", "kubernetes-dashboard", "create", "token", saName)
	print("and enter token: ")
	tkn.Stdout = os.Stdout
	tkn.Stderr = os.Stderr
	err = tkn.Run()
	util.FatalErrorCheck(err)

}

// BadPodDashUpdate TODO: expand this to have field values expressing that the pod could not be examined
func BadPodDashUpdate() DashboardUpdate {
	return DashboardUpdate{Denied: true}
}
