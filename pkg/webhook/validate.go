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

// Vulns the fields for the vulnerability data
type Vulns struct {
	ID         string `json:"id"`
	DataSource string `json:"dataSource"`
	NameSpace  string `json:"namespace"`
	Severity   string `json:"severity"`
}

// Policy is a special struct to read the admission_policy file that sets the rules for the webhook
type Policy struct {
	SeverityLimit string   `json:"severity_limit"`
	Whitelist     []string `json:"id_whitelist"`
}

// convertSeverityString takes a string form of a severity for admission policy and converts it to the comparable integer form
func convertSeverityString(severity string) int {
	sevVal := 0
	if severity == "Negligible" {
		sevVal = 0
	} else if severity == "Low" {
		sevVal = 1
	} else if severity == "Medium" {
		sevVal = 2
	} else if severity == "High" {
		sevVal = 3
	} else if severity == "Critical" {
		sevVal = 4
	} else {
		// if it is "Unknown" then assume the worst,
		// or "Setup" is for pre-approved pods as part of the cluster starting
		sevVal = 5
	}
	return sevVal
}

// compareSeverity takes the given CVE severity and checks if it falls within the given limit as dictated by the
// admission policy, and if the CVE is within the limit, return true
func compareSeverity(givenSeverity string, limit int) bool {
	givenInt := convertSeverityString(givenSeverity)

	if givenInt >= limit {
		return false
	}
	return true
}

// ConstructPolicy reads in the admission_policy.json file and parses it into usable data via the Policy struct
func ConstructPolicy(policyFile string) (int, map[string]int) {
	// read the file back in as a string
	rawContent := util.ReadFile(policyFile)

	// read the admission policy into the custom struct
	var policyInfo Policy
	_ = json.Unmarshal([]byte(rawContent), &policyInfo)

	// get the severity limit and convert to int for easier comparisons
	rawSeverity := policyInfo.SeverityLimit
	severityLimit := convertSeverityString(rawSeverity)
	// make a map to speed up search time later
	whiteListMap := make(map[string]int)
	for _, id := range policyInfo.Whitelist {
		whiteListMap[id] = 1
	}

	// return the two pieces of info needed to enforce the policy
	return severityLimit, whiteListMap
}

// generateSBOM generates an SBOM from an image and stores it in an argued path
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
func evaluateImage(sbomFile string, imageName string, severityLimit int, whiteListMap map[string]int) []string {
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

	// list to store each CVE that breaks policy from this image
	var cveList []string

	// result is now a list of matches that can be iterated through
	for i := 0; i < len(result.Matches); i++ {
		// get the current CVE id and severity level
		currID := result.Matches[i].Vulnerability.ID
		currSeverity := result.Matches[i].Vulnerability.Severity
		// if not within the limit, check the whitelist
		if !compareSeverity(currSeverity, severityLimit) {
			// don't care about the value, just whether or not the id exists in the white-list (present is a bool)
			_, present := whiteListMap[currID]
			// if not in the whitelist
			if !present {
				// add to the list of CVEs that break policy
				cveList = append(cveList, currID)
			}
		}
	}
	return cveList
}

// checkPodImages pulls out all images from a pod struct and sends them to the DB interface,
// which then checks if an SBOM exists for each (if not, then sends the image to syft) and then,
// based off the result of grype (which should return to this function) and says what CVEs
// exist within each image, and if any of those CVEs are unacceptable, the whole pod is Denied
func checkPodImages(pod *corev1.Pod) (dashboard.DashboardUpdate, error) {
	containers := pod.Spec.Containers
	// get the number of images
	sliceSize := len(containers)
	// setup an empty slice to hold each image
	imageSlice := make([]string, sliceSize)

	// get the security policy here (any amount of repeated calculations someone would argue, I argue this adds per pod flexibility, and NO I don't care who says to change it
	severityLimit, whiteListMap := ConstructPolicy("webhook/admission_policy.json")

	counter := 0
	// extract all images and store in the list
	failure := false
	// make a map[string][]string to store image name as key, and it's cveList as value
	imageCVEMap := make(map[string][]string)
	for range containers {
		imageSlice[counter] = containers[counter].Image

		// get the time to make the file names unique
		timeRaw := util.FormatTime()
		// EX: sboms/nginx_sbom_2023-3-20_17-57-50.json
		sbomName := "sboms/" + containers[counter].Image + "_sbom_" + timeRaw + ".json"
		generateSBOM(sbomName, containers[counter].Image)

		grypeRes := evaluateImage(sbomName, containers[counter].Image, severityLimit, whiteListMap)
		log.Print("validate.go: checkPodImages -> successfully evaluated vulnerabilities")

		// if any CVE's broke policy then add this image to the map of bad images
		if len(grypeRes) > 0 {
			imageCVEMap[containers[counter].Image] = grypeRes
			failure = true
		}

		counter++
	}

	// get the pod name for the print-out
	podName := pod.ObjectMeta.Name

	// currently rejects any pod with an image containing a Critical level CVE
	if failure {
		log.Print("validate.go: checkPodImages -> " + podName + " was denied")
		return dashboard.DashboardUpdate{Denied: true, CVEList: imageCVEMap}, nil
	}

	log.Print("validate.go: checkPodImages -> " + podName + " was accepted")
	return dashboard.DashboardUpdate{Denied: false, CVEList: imageCVEMap}, nil
}
