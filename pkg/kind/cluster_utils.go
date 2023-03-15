package kind

import (
	"fmt"
	"github.com/usfca-cs490/admissions-webhook/pkg/util"
	"log"
	"os"
	"os/exec"
)

// CreateCluster Method to create a cluster using kind
func CreateCluster() {
	// Create command
	cmd := exec.Command("kind", "create", "cluster")
	// Redirect stdout and stderr to default for OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and handle errors
	err := cmd.Run()
	util.FatalErrorCheck(err)
}

// Shutdown Method to shutdown a cluster
func Shutdown() {
	// Create command
	cmd := exec.Command("kind", "delete", "cluster")
	// Redirect stdout and stderr to default for OS
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and handle errors
	if err := cmd.Run(); err != nil {
		log.Fatal("startup.go: FAILED TO GET CLUSTER INFO")
	}
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
	util.FatalErrorCheck(err)

	// Status print
	fmt.Println("Loading Docker image", (image_name + ":" + version))

	// Create command
	cmd = exec.Command("kind", "load", "docker-image", (image_name + ":" + version))

	// Run and handle errors
	err = cmd.Run()
	// Crash if error
	util.FatalErrorCheck(err)
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
    util.FatalErrorCheck(err)

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
	util.NonfatalErrorCheck(err)
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
	util.NonfatalErrorCheck(err)
}
