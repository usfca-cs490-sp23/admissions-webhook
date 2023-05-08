package kind

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/usfca-cs490/admissions-webhook/pkg/util"
)

// Pod struct for wide output of kubectl get pods
type Pod struct {
	Namespace       []byte
	Name            []byte
	Ready           []byte
	Status          []byte
	Restarts        []byte
	Age             []byte
	Ip              []byte
	Node            []byte
	Nominated_node  []byte
	Readiness_gates []byte
}

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
	cmd := exec.Command("kubectl", "logs", "-l", string("app="+pod_name), "-f")

	// Redirect output to terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and handle errors
	err := cmd.Run()
	util.FatalErrorCheck(err, true)
}

// GetPods gets the pods in the kind-control-plane
func GetPods(node_name string) []byte {
	// Create a logging command with kubectl
	out, err := exec.Command("kubectl", "get", "pods", "--all-namespaces", "-o", "wide",
		"--field-selector", string("spec.nodeName="+node_name)).Output()

	// Run the command and handle errors
	util.NonfatalErrorCheck(err, true)

	return out
}

// GetPodsStruct gets the pod data from kubectl and stores the results in an array of structs
func GetPodsStruct(node_name string) []Pod {
	// Get the pods using kubectl wide option
	getPods := string(GetPods(node_name))
	podStats := strings.Split(getPods, "\n")

	// Initialize a new list of Pods
	var pods []Pod = make([]Pod, len(podStats)-1)

	// Initialize the counter of pods in the array
	podCtr := 0

	// Loop through the matches (ignoring the header)
	for _, stat := range podStats[1:] {
		// Get fields between whitespaces
		fields := strings.Fields(stat)
		// If there are not ten fields, this is incomplete data, do not add it to struct
		if len(fields) != 10 {
			continue
		}

		// Store the relevant fields in a struct within the pods array
		pods[podCtr] = Pod{
			Namespace:       []byte(fields[0]),
			Name:            []byte(fields[1]),
			Ready:           []byte(fields[2]),
			Status:          []byte(fields[3]),
			Restarts:        []byte(fields[4]),
			Age:             []byte(fields[5]),
			Ip:              []byte(fields[6]),
			Node:            []byte(fields[7]),
			Nominated_node:  []byte(fields[8]),
			Readiness_gates: []byte(fields[9]),
		}
		podCtr++
	}

	return pods
}

// FindPod finds a Pod in a list of Pod structs by its name (first find)
func FindPod(pods []Pod, target_name string) *Pod {
	// Loop through pods
	for i := range pods {
		// If there is a match, return it
		if len(string(pods[i].Name)) >= len(target_name) {
			if string(pods[i].Name[0:len(target_name)]) == target_name {
				return &pods[i]
			}
		}
	}
	return nil
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

// BuildLoadHookImage Method to build an image from a specified Dockerfile
func BuildLoadHookImage(image_name, version, dfile_path string) {
	// Status print
	fmt.Println("Building Docker image", (image_name + ":" + version), "from Dockerfile at", dfile_path)

	// Create command
	os.Setenv("DOCKER_BUILDKIT", "1")
	// Get time for benchmarking build time
	loadHookStart := time.Now()
	cmd := exec.Command("docker", "build", "-t", (image_name + ":" + version), dfile_path)
	cmd.Stderr = os.Stderr

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.FatalErrorCheck(err, true)

	// Get build time benchmark benchmark
	buildTime := (time.Since(loadHookStart))

	// Status print
	fmt.Println("Loading Docker image", (image_name + ":" + version))

	// Create command
	cmd = exec.Command("kind", "load", "docker-image", (image_name + ":" + version))

	// Run and handle errors
	err = cmd.Run()
	// Crash if error
	util.FatalErrorCheck(err, true)

	// Using hidden policy file trickery to let the redis pod in at startup
	// read in each file and store the data
	userContents := util.ReadFile("./pkg/webhook/admission_policy.json")
	defaultContents := util.ReadFile("./pkg/webhook/.default_policy.json")

	// Write the default data to the user file to allow redis into the cluster
	util.WriteFile("./pkg/webhook/admission_policy.json", defaultContents)

	// If the amount of time to build is less than 12 seconds, wait until 12 have elapsed
	// 	to allow sufficient time for default namespace to be setup before redis enters cluster
	t, _ := time.ParseDuration("12s")
	time.Sleep(time.Duration(t.Nanoseconds() - buildTime.Nanoseconds()))

	// Configure and apply redis
	CreateConfigMap("./pkg/webhook/database/redis-config.yaml")
	AddPod("./pkg/webhook/database/redis-pod.yaml")

	// Now write back the user info
	util.WriteFile("./pkg/webhook/admission_policy.json", userContents)

	// Not actually making a config map, its a service, but its the same command
	CreateConfigMap("./pkg/webhook/database/redis-service-config.yaml")

	// Set up cluster role so that events are properly sent
	CreateAdminRole("dashboard", "default:default")
	// Set up cluster role for kubeaoudit
	CreateAdminRole("daemonsets.apps", "default:kubeaudit")
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

// AddPod Method to add a pod
func AddPod(pod_config_path string) {
	// Create command
	cmd := exec.Command("kubectl", "apply", "-f", pod_config_path)
	// Redirect stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and handle errors
	err := cmd.Run()
	//if error create custom event
	if err != nil {
		fmt.Printf("pod/%s denied\n", pod_config_path)
	}
}

// CreateConfigMap method to create a new ConfigMap for a pod
func CreateConfigMap(config_path string) {
	// Create command
	cmd := exec.Command("kubectl", "create", "-f", config_path)
	// Redirect stdout
	cmd.Stdout = os.Stdout

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)
}

// DeleteItem takes an item type and its name in order to delete it from a cluster
func DeleteItem(type_, name string) {
	// Create command
	cmd := exec.Command("kubectl", "delete", type_, name)
	// Redirect stdout
	cmd.Stdout = os.Stdout

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)
}

// CreateAdminRole method to create a clusterrolebinding within the cluster with admin level priveleges
func CreateAdminRole(name, serviceacc string) {
	// kubectl create clusterrolebinding dashboard --clusterrole=cluster-admin --serviceaccount=default:default
	// Create command
	cmd := exec.Command("kubectl", "create", "clusterrolebinding", name, "--clusterrole=cluster-admin", "--serviceaccount="+serviceacc)
	// Redirect stdout
	cmd.Stdout = os.Stdout

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)
}

// DeletePod takes a pod namespace and its name in order to delete it from a cluster
func DeletePod(namespace string, name string) {
	// Create command
	cmd := exec.Command("kubectl", "delete", "pod", "-n", namespace, name)
	// Redirect stdout
	cmd.Stdout = os.Stdout

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)
}

// CopyPolicy copy the policy to the webhook so it can be applied
func CopyPolicy(hookName string) {
	// final command should look like: kubectl cp ./pkg/webhook/admission_policy.json default/the-captains-hook-646c87d54-nlqqx:webhook/admission_policy.json
	// Create command
	target := "default/" + hookName + ":webhook/admission_policy.json"
	cmd := exec.Command("kubectl", "cp", "./pkg/webhook/admission_policy.json", target)
	// Redirect stdout
	cmd.Stdout = os.Stdout

	// Run and handle errors
	err := cmd.Run()
	// Crash if error
	util.NonfatalErrorCheck(err, true)
}
