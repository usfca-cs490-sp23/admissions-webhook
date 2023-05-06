package util

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	corev1 "k8s.io/api/core/v1"
)

// NotYetImplemented Helper method to panic and trace to source method for unimplemented code
func NotYetImplemented(method string) {
	panic((method + " not yet implemented"))
}

// IsFlagRaised Helper method to check if a flag has been raised
func IsFlagRaised(flag_name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flag_name {
			found = true
		}
	})
	return found
}

// FatalErrorCheck Helper method to crash if errors exist
func FatalErrorCheck(err error, verbose bool) {
	if err != nil {
		if verbose {
			log.Print("\nERROR Fatal: " + err.Error() + "\n")
		} else {
			log.Print(err)
		}
		log.Fatal(err)
	}
}

// NonfatalErrorCheck Helper method to output present errors but not crash
func NonfatalErrorCheck(err error, verbose bool) {
	if err != nil {
		if verbose {
			log.Print("\nERROR Nonfatal: " + err.Error() + "\n")
		} else {
			log.Print(err)
		}
	}
}

// ReadFile Helper method to read a file and return it as a string
func ReadFile(infile string) string {
	// Open the current file and generate reader
	f, err := os.Open(infile)
	FatalErrorCheck(err, true)
	defer f.Close()

	// Read the current file
	content, err := io.ReadAll(f)
	// Crash if error
	FatalErrorCheck(err, true)

	return string(content)
}

// WriteFile Helper method to write a file
func WriteFile(outfile, data string) {
	// Create the file
	f, err := os.Create(outfile)
	// Crash if error
	FatalErrorCheck(err, true)
	// Close the file with defer
	defer f.Close()

	// Write the data
	f.WriteString(data)
}

// FormatTime generates a string of the current time for naming files
func FormatTime() string {
	// make filename based off of current time (lazy but effective)
	currTime := time.Now()
	// format is from the docs, plz don't change it bc the date is based on the underlying schema used by the Format()
	fileName := currTime.Format("2006-1-2_15-4-5")
	return fileName
}

// InjectYamlCA Method to inject a CA bundle into a YAML file
func InjectYamlCA(target, template, injectable string) {
	// Read and store file with tls data
	content := ReadFile(injectable)
	// Read and store file to inject tls into
	config := ReadFile(template)

	// Remove unnecessary prefix from content
	content = strings.TrimPrefix(content, ">> MutatingWebhookConfiguration caBundle:")
	// Indent content correctly
	content = strings.ReplaceAll(content, "\n", "\n                ")

	// Insert the content into the config file string
	config = strings.Replace(config, "caBundle: |", "caBundle: |\n                "+content, 1)

	// Now write to file
	WriteFile(target, config)
}

func ChangeConfig(level string, path string) {
	content := ReadFile(path)
	r, _ := regexp.Compile("\"severity_limit\": [^\n]*")
	new_level := "\"severity_limit\": \"" + level + "\","

	content = r.ReplaceAllString(content, new_level)

	WriteFile(path, content)
}

func WritePodEvent(podName string, reason bool, message map[string][]string) {
	var newMess string
	var cveMessage string
	for key, val := range message {
		// Convert each key/value pair in m to a string
		newMess = fmt.Sprintf("%s=\"%s\"", key, val)
		cveMessage = cveMessage + newMess
	}

	var eventReason string
	var eventType string
	// if true then the pod was denied
	if reason {
		eventReason = "Pod Denied"
		eventType = "Warning"
	} else {
		eventReason = "Pod accepted"
		eventType = "Normal"
	}

	// set a name that will change even for duplicate pods
	eventName := podName + FormatTime()

	event := &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      eventName,
			Namespace: "apps",
		},
		InvolvedObject: corev1.ObjectReference{Namespace: "apps"},
		Reason:         eventReason,
		Message:        cveMessage,
		FirstTimestamp: metav1.Time{
			Time: time.Now(),
		},
		LastTimestamp: metav1.Time{
			Time: time.Now(),
		},
		Count: 1,
		Type:  eventType,
		Source: corev1.EventSource{
			Component: "the-captains-hook",
		},
	}

	// set the api version (doing here bc I keep getting warnings when I put it in the struct)
	event.APIVersion = "v1"

	var config *rest.Config
	var err error
	// Load kubeconfig from $HOME/.kube/config or in-cluster configuration
	if _, err = os.Stat(os.Getenv("HOME") + "/.kube/config"); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	} else {
		config, err = rest.InClusterConfig()
	}
	NonfatalErrorCheck(err, true)

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	NonfatalErrorCheck(err, true)

	// Send the event to the dashboard
	_, err = clientset.CoreV1().Events("apps").Create(context.Background(), event, metav1.CreateOptions{})
	NonfatalErrorCheck(err, true)
	log.Print("Event sent\n")
}

// WriteRedeployEvent sends an event with the results of a cluster re-validation
func WriteRedeployEvent(reason string, evictionList []string) {
	evictedPods := ""
	for _, evicted := range evictionList {
		evictedPods = evictedPods + evicted
	}

	var eventReason = reason
	eventType := "Normal"

	// set a name that will change even for duplicate pods
	eventName := "" + FormatTime()

	event := &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      eventName,
			Namespace: "apps",
		},
		InvolvedObject: corev1.ObjectReference{Namespace: "apps"},
		Reason:         eventReason,
		Message:        evictedPods,
		FirstTimestamp: metav1.Time{
			Time: time.Now(),
		},
		LastTimestamp: metav1.Time{
			Time: time.Now(),
		},
		Count: 1,
		Type:  eventType,
		Source: corev1.EventSource{
			Component: "the-captains-hook",
		},
	}

	// set the api version (doing here bc I keep getting warnings when I put it in the struct)
	event.APIVersion = "v1"

	var config *rest.Config
	var err error
	// Load kubeconfig from $HOME/.kube/config or in-cluster configuration
	if _, err = os.Stat(os.Getenv("HOME") + "/.kube/config"); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	} else {
		config, err = rest.InClusterConfig()
	}
	NonfatalErrorCheck(err, true)

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	NonfatalErrorCheck(err, true)

	// Send the event to the dashboard
	_, err = clientset.CoreV1().Events("apps").Create(context.Background(), event, metav1.CreateOptions{})
	NonfatalErrorCheck(err, true)
	log.Print("Event sent\n")
}
