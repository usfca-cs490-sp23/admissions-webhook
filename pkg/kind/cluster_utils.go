package kind

import (
	"fmt"
	"github.com/usfca-cs490/admissions-webhook/pkg/util"
	"log"
	"os"
	"os/exec"
    "strings"
)

// GetPodName get an argued pod's name
func GetPodName(pod_name string) string {
    // Call kubectl describe on the argued pod name
    hook_desc, err := exec.Command("kubectl", "describe", "pod", pod_name).Output()
    // Crash if there is an error
    util.FatalErrorCheck(err, true)

    // Extract just the name from the description
    hook_desc_str := string(hook_desc)
    hook_desc_str = hook_desc_str[0:strings.Index(hook_desc_str, "\n")]
    hook_desc_str = strings.Trim(strings.TrimPrefix(hook_desc_str, "Name:"), " ")

    // Return the name
    return hook_desc_str
}

// StreamLogs send the logs to stdout
func StreamLogs(pod_name string) {
    // Create a logging command with kubectl
    cmd := exec.Command("kubectl", "logs", "-l", string("app=" + pod_name), "-f")

    // Redirect output to terminal
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Run the command and handle errors
    err := cmd.Run()
    util.FatalErrorCheck(err, true)
}

// CreateCluster Method to create a cluster using kind
func CreateCluster() {
	// Create command
	cmd := exec.Command("kind", "create", "cluster")
	// Redirect stdout and stderr to default for OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and handle errors
	err := cmd.Run()
	util.FatalErrorCheck(err, false)
}

// Shutdown Method to shutdown a cluster
func Shutdown() {
	// Create command
	cmd := exec.Command("kind", "delete", "cluster")
	// Redirect stdout and stderr to default for OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and handle errors
    err := cmd.Run()
    util.FatalErrorCheck(err, true)
}

// Info Method to check kind cluster info
func Info() {
	// Create command
	cmd := exec.Command("kubectl", "cluster-info", "--context", "kind-kind")
	// Redirect stdout and stderr to default for OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and handle errors
	if err := cmd.Run(); err != nil {
		log.Fatal("startup.go: FAILED TO GET CLUSTER INFO")
	}
}

// BuildLoadHookImage Methoed to build an image from a specified Dockerfile
func BuildLoadHookImage(image_name, version, dfile_path string) {
	// Status print
	fmt.Println("Building Docker image", (image_name + ":" + version), "from Dockerfile at", dfile_path)

	// Create command
	os.Setenv("DOCKER_BUILDKIT", "1")
	cmd := exec.Command("docker", "build", "-t", (image_name + ":" + version), dfile_path)
	cmd.Stderr = os.Stderr

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.FatalErrorCheck(err, true)

	// Status print
	fmt.Println("Loading Docker image", (image_name + ":" + version))

	// Create command
	cmd = exec.Command("kind", "load", "docker-image", (image_name + ":" + version))

	// Run and handle errors
	err = cmd.Run()
	// Crash if error
	util.FatalErrorCheck(err, true)
}

// GenCerts Method to generate TLS certifications and cluster configs
func GenCerts() {
    // Status print
    fmt.Println("Generating TLS certificate, key, and CA bundle and injecting into configuration files")

    // Create command
    cmd := exec.Command("/bin/sh", "./pkg/tls/gen_certs.sh")

    // Run and handle errors
    err := cmd.Run()
    // Crash if error
    util.FatalErrorCheck(err, true)

    // Inject CA Bundle into validating.config.yaml
    util.InjectYamlCA("./pkg/cluster-config/validating.config.yaml",
    "./pkg/cluster-config/validating.config.template.yaml", "./pkg/tls/cab64.crt")
}

// ApplyConfig Method to apply a configuration YAML file to a cluster using kubectl
func ApplyConfig(config_file string) {
	// Status print
	fmt.Println("Applying config", config_file, "to cluster")

	// Create command
	cmd := exec.Command("kubectl", "apply", "-f", config_file)

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)
}

// DescribeHook Method to get hook pod data
func DescribeHook(hook_name string) {
	// Create command
	cmd := exec.Command("kubectl", "describe", "pod", hook_name)
	// Redirect stdout
	cmd.Stdout = os.Stdout

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)
}

