package webhook

import (
	"fmt"
    "os/exec"
    "github.com/usfca-cs490/admissions-webhook/pkg/dashboard"
    "github.com/usfca-cs490/admissions-webhook/pkg/util"

	corev1 "k8s.io/api/core/v1"
)

// GenerateSBOM generates an SBOM from an image and stores it in an argued path
func GenerateSBOM(outfile, image string) {
	// Create and run command
    out, err := exec.Command("syft","-o", "json", image).Output()
    // Crash if there are any errors
    util.FatalErrorCheck(err)

    // Write output to file
    util.WriteFile(outfile, string(out))
}

// TODO: add all functionality
// checkPodImages pulls out all images from a pod struct and sends them to the DB interface,
// which then checks if an SBOM exists for each (if not, then sends the image to syft) and then,
// based off the result of grype (which should return to this function) and says what CVEs
// exist within each image, and if any of those CVEs are unacceptable, the whole pod is Denied
func checkPodImages(pod *corev1.Pod) (dashboard.DashboardUpdate, error) {
	// TODO: pod.Spec.ImagePullSecrets // should allow to get all images from a pod? (maybe just secret ones?)
	// get the list of all given containers in this pod
	containers := pod.Spec.Containers
	// get the number of images
	sliceSize := len(containers)
	// setup an empty slice to hold each image
	imageSlice := make([]string, sliceSize)

	// TODO: this is a shit way of doing this, but I just want to see it work right now, clean this up later
	counter := 0
	// extract all images and store in the list
	for range containers {
		imageSlice[counter] = containers[counter].Image
		counter++
	}

	singleImage := pod.Spec.Containers[0].Image
	fmt.Println(singleImage)

	// TODO: pass each image to the DB interface from here and receive the grype results

	// currently allows any pod into cluster
	return dashboard.DashboardUpdate{Denied: false}, nil
}

