package dashboard

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/usfca-cs490/admissions-webhook/pkg/util"
)

// DashboardUpdate TODO: Return a DashboardUpdate struct with the result of checking the internals of the pod
type DashboardUpdate struct {
	CVEList map[string][]string
	Denied  bool
}

// DashInit initiates the dashboard on user's local computer
func DashInit() {
	cmd := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml")
	// Run and handle errors
	err := cmd.Run()
	util.FatalErrorCheck(err, false)

	//print("go to: http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/\n")
	DashUser("./pkg/dashboard/dashboard-adminuser.yaml", "./pkg/dashboard/admin-rb.yaml")
	OpenLink("http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/")

	cmd = exec.Command("kubectl", "proxy")
	// Run and handle errors
	err = cmd.Run()
	util.FatalErrorCheck(err, false)
	//go to http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ to access
}

// DashUser creates a new user account via K8's Service Account mechanism
func DashUser(adminUser string, adminRb string) {
	//create admin service account (see dashboard-adminuser.yaml)
	cmd := exec.Command("kubectl", "apply", "-f", adminUser)
	err := cmd.Run()
	util.FatalErrorCheck(err, false)

	//create cluster role binding (see admin-rb.yaml)
	cmd = exec.Command("kubectl", "apply", "-f", adminRb)
	err = cmd.Run()
	util.FatalErrorCheck(err, false)

	//name of the service account
	saName := "admin-user"

	tkn := exec.Command("kubectl", "-n", "kubernetes-dashboard", "create", "token", saName)
	print("and enter token: ")
	tkn.Stdout = os.Stdout
	tkn.Stderr = os.Stderr
	err = tkn.Run()
	util.FatalErrorCheck(err, false)

}

func OpenLink(link string) {
	var cmd string
	os := runtime.GOOS

	if os == "windows" {
		cmd = "cmd /c start"
	} else if os == "darwin" {
		cmd = "open"
	} else {
		cmd = "xdg-open"
	}
	site := exec.Command(cmd, link)
	err := site.Run()
	util.FatalErrorCheck(err, false)
}

// BadPodDashUpdate TODO: expand this to have field values expressing that the pod could not be examined
func BadPodDashUpdate() DashboardUpdate {
	return DashboardUpdate{Denied: true}
}
