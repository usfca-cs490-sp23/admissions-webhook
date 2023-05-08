package dashboard

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/usfca-cs490/admissions-webhook/pkg/util"
)

// DashboardUpdate struct for storing results that are eventually sent in an event to the dashboard
type DashboardUpdate struct {
	CVEList map[string][]string
	Denied  bool
	PodName string
}

// DashInit initiates the dashboard on user's local computer
func DashInit() {
	// Execute dashboard command
	cmd := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml")
	// Run and handle errors
	err := cmd.Run()
	util.NonfatalErrorCheck(err, false)

	// Send user instructions to command line
	fmt.Println("go to: http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/")
	DashUser("./pkg/dashboard/dashboard-adminuser.yaml", "./pkg/dashboard/admin-rb.yaml")
	OpenLink("http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/")
	RunDashboard()
}

// DashUser creates a new user account via K8's Service Account mechanism
func DashUser(adminUser, adminRb string) {
	// Create admin service account (see dashboard-adminuser.yaml)
	cmd := exec.Command("kubectl", "apply", "-f", adminUser)
	err := cmd.Run()
	util.NonfatalErrorCheck(err, false)

	// Create cluster role binding (see admin-rb.yaml)
	cmd = exec.Command("kubectl", "apply", "-f", adminRb)
	err = cmd.Run()
	util.NonfatalErrorCheck(err, false)

	// Name of the service account
	saName := "admin-user"

	// Create a token
	tkn := exec.Command("kubectl", "-n", "kubernetes-dashboard", "create", "token", saName)
	// Give user to token, along with usage instructions
	fmt.Print("and enter token: ")
	out, err := tkn.CombinedOutput()
	fmt.Print(string(out))
	util.NonfatalErrorCheck(err, false)
	CopyTkn(string(out))
}

// CopyTkn copies the bearer token (for login) to the user's clipboard
func CopyTkn(code string) {
	os := runtime.GOOS
	var command string
	if os == "windows" {
		command = "clip"
	} else if os == "darwin" {
		command = "pbcopy"
	} else if os == "linux" {
		command = "xclip -sel clip"
	}
	code = "echo \"" + code + "\" | " + command
	exec.Command("bash", "-c", code).Run()
}

// RunDashboard initiates the dashboard on http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/
func RunDashboard() {
	cmd := exec.Command("kubectl", "proxy")
	// Run and handle errors
	err := cmd.Run()
	util.NonfatalErrorCheck(err, false)
}

// OpenLink contains cross-os compatibility to open the dashboard link in the user's preferred browser
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
	util.NonfatalErrorCheck(err, false)
}

// ProcessEvent takes a filled out Dash struct and puts it's info into a dashboard event
func ProcessEvent(update DashboardUpdate, eventTypePod bool, infoList []string) {
	if eventTypePod == true {
		util.WritePodEvent(update.PodName, update.Denied, update.CVEList)
	} else {
		util.WriteRedeployEvent("Redeploy event", infoList)
	}

}

// BadPodDashUpdate a failover func for when things really go sideways
func BadPodDashUpdate() DashboardUpdate {
	return DashboardUpdate{Denied: true}
}
