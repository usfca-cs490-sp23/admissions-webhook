package webhook

import (
	"encoding/json"
	"github.com/usfca-cs490/admissions-webhook/pkg/dashboard"
	"github.com/usfca-cs490/admissions-webhook/pkg/util"
	"log"
	"os/exec"

	corev1 "k8s.io/api/core/v1"
)

// Evaluation special struct that acts as the top level of the json parsing structure for reading grype results
type Evaluation struct {
	Matches []Match `json:"matches"`
}
type Match struct {
	Vulnerability Vulns `json:"vulnerability"`
}

// Vulns TODO: do the rest of the fileds if need be
type Vulns struct {
	ID         string `json:"id"`
	DataSource string `json:"dataSource"`
	NameSpace  string `json:"namespace"`
	Severity   string `json:"severity"`
}

// generateSBOM generates an SBOM from an image and stores it in an argued path (currently: ./pkg/sboms/<filename>)
func generateSBOM(outfile, image string) {
	// Create and run command
    log.Print("validate.go: GenerateSBOM -> creating SBOM for " + image + "...")
	out, err := exec.Command("syft", "-o", "json", image).Output()
	// Crash if there are any errors
	util.FatalErrorCheck(err, true)
	log.Print("validate.go: GenerateSBOM -> created SBOM for " + image + " and stored at " + outfile)

	// Write output to file
	util.WriteFile(outfile, string(out))
}

// runGrypeOnSingleImage takes a string representing an sbom in json format, runs the grype command,
// writes the grype evaluation results to a json file, reads the evaluation into the special struct,
// checks if that image breaks the security rules (has "Critical" rated CVEs),
// and returns false if the rules are not broken
func evaluateImage(sbomFile string, imageName string) bool {
	timeName := util.FormatTime()
	//EX: evals/nginx_eval_2023-3-20_17-57-50.json
	outFile := "evals/" + imageName + "_eval_" + timeName + ".json"

	// run grype command
	givenSBOM := "sbom:./" + sbomFile
	// To scan an SBOM: grype sbom:./example.json
    log.Print("validate.go: evaluateImage -> running grype on SBOM at " + givenSBOM)
	out, err := exec.Command("grype", givenSBOM, "-o", "json").Output()
	// Crash if there are any errors
	util.FatalErrorCheck(err, true)
    log.Print("validate.go: evaluateImage -> grype ran successfully and stored output in " + outFile)

	// Write output to file
	util.WriteFile(outFile, string(out))

	// read the file back in as a string
	rawContent := util.ReadFile(outFile)

	// create and populate a struct that is tailor-made for the json structure output by grype
	var result Evaluation
	_ = json.Unmarshal([]byte(rawContent), &result)

	// result is now a list of matches that can be iterated through
	for i := 0; i < len(result.Matches); i++ {
		// This rule can be changed easily in the future via yaml/json rules that can be read in from config
		if result.Matches[i].Vulnerability.Severity == "Critical" {
			return true
		}
	}
	return false
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
	failure := false
	for range containers {
		imageSlice[counter] = containers[counter].Image

		// TODO: don't do any this here, this is just the proof of concept
		timeRaw := util.FormatTime()
		// EX: sboms/nginx_sbom_2023-3-20_17-57-50.json
		sbomName := "sboms/" + containers[counter].Image + "_sbom_" + timeRaw + ".json"
		generateSBOM(sbomName, containers[counter].Image)
		// we could immediately stop and break here, but I think it's worth checking all the images in the pod
		// TODO: change this so that it makes a list of the images that failed and try to come up with a way that can say what CVEs caused the failure? idk
		grypeRes := evaluateImage(sbomName, containers[counter].Image)
        log.Print("validate.go: checkPodImages -> successfully evaluated vulnerabilities")
		if grypeRes == true {
			failure = true
		}

		counter++
	}

	singleImage := pod.Spec.Containers[0].Image

	// TODO: pass each image to the DB interface from here and receive the grype results

	// currently rejects any pod with an image containing a Critical level CVE
	if failure {
        log.Print("validate.go: checkPodImages -> " + singleImage + " was denied")
		return dashboard.DashboardUpdate{Denied: true}, nil
	}

    log.Print("validate.go: checkPodImages -> " + singleImage + " was accepted")
	return dashboard.DashboardUpdate{Denied: false}, nil
}
