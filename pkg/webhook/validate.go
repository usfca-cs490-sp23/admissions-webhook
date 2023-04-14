package webhook

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/usfca-cs490/admissions-webhook/pkg/dashboard"
	"github.com/usfca-cs490/admissions-webhook/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"log"
	"os/exec"
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

type Validator struct {
	Severity  int
	WhiteList map[string]struct{}
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
func ConstructPolicy(policyFile string) *Validator {
	// read the file back in as a string
	rawContent := util.ReadFile(policyFile)

	// read the admission policy into the custom struct
	var policyInfo Policy
	_ = json.Unmarshal([]byte(rawContent), &policyInfo)

	// get the severity limit and convert to int for easier comparisons
	rawSeverity := policyInfo.SeverityLimit
	severityLimit := convertSeverityString(rawSeverity)
	// make a map to speed up search time later (essentailly making a set)
	whiteListMap := make(map[string]struct{})
	for _, id := range policyInfo.Whitelist {
		whiteListMap[id] = struct{}{}
	}

	// return the two pieces of info needed to enforce the policy in a Validator literal
	return &Validator{severityLimit, whiteListMap}
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
func evaluateImage(sbomFile string, imageName string, severityLimit int, whiteListMap map[string]struct{}) []string {
	// run grype command
	givenSBOM := "sbom:./" + sbomFile
	// To scan an SBOM: grype sbom:./example.json
	log.Print("validate.go: evaluateImage -> running grype on SBOM at " + givenSBOM)
	out, err := exec.Command("grype", givenSBOM, "-o", "json").Output()
	// Crash if there are any errors
	util.FatalErrorCheck(err, true)
	log.Print("validate.go: evaluateImage -> grype ran successfully")

	// Got rid of extra file IO here, big win

	// create and populate a struct that is tailor-made for the json structure output by grype
	var result Evaluation
	_ = json.Unmarshal(out, &result)

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

// dbLookup retrieves an SBOM from the database if it exists
func dbLookup(dbKey string, monthNum int) {
	client := redis.NewClient(&redis.Options{
		Addr:     "10.244.0.5:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Just using this bc it is apparently the 'default' for starting contexts, if there is a better one plz fix
	ctx := context.Background()

	// Get the direct value from
	val, err := client.Get(ctx, dbKey).Result()

	// TODO: Look for current first, then if not present look for old, if old present remove, then generate new (and put later)
	// key did not exist
	if err == redis.Nil {
		if val == "" {
			return
		}
	}

	if err != nil {
		// TODO: remove panic
		panic(err)
	}

}

// TODO: parse the file name so that the month is included with the image name, so at least once a month a fresh SBOM gets made (this avoids the issue of image versioning)
// dbStore stores an SBOM given it's file name and ...??????
func dbStore() {
	client := redis.NewClient(&redis.Options{
		Addr:     "10.244.0.5:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Just using this bc it is apparently the 'default' for starting contexts, if there is a better one plz fix
	ctx := context.Background()

	// This is an example of how to store values into the database
	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		// TODO: remove panic
		panic(err)
	}
}

// checkPodImages pulls out all images from a pod struct and sends them to the DB interface,
// which then checks if an SBOM exists for each (if not, then sends the image to syft) and then,
// based off the result of grype (which should return to this function) and says what CVEs
// exist within each image, and if any of those CVEs are unacceptable, the whole pod is Denied
func (v *Validator) checkPodImages(pod *corev1.Pod) (dashboard.DashboardUpdate, error) {
	containers := pod.Spec.Containers

	failure := false
	// make a map[string][]string to store image name as key, and it's cveList as value
	imageCVEMap := make(map[string][]string)
	for _, container := range containers {
		// get the time to make the file names unique
		timeRaw := util.FormatTime()

		//TODO: parse name, try to grab from db, if fail or different month, make new sb and send to db
		//monthString := strings.Split(timeRaw, "-")[1]
		//monthInt, _ := strconv.Atoi(monthString)
		//sbomDbName := containers[counter].Image + monthString
		// dbLookup(sbomDbName, monthInt)
		//                      year-month-day
		// EX: sboms/nginx_sbom_2023-3-20_17-57-50.json
		// TODO: fix for all file names
		//sbomName := fmt.Sprintf("sboms/%s_sbom_%s.json", container.Image, timeRaw)
		sbomName := "sboms/" + container.Image + "_sbom_" + timeRaw + ".json"
		generateSBOM(sbomName, container.Image)

		grypeRes := evaluateImage(sbomName, container.Image, v.Severity, v.WhiteList)
		log.Print("validate.go: checkPodImages -> successfully evaluated vulnerabilities")

		// if any CVE's broke policy then add this image to the map of bad images
		if len(grypeRes) > 0 {
			imageCVEMap[container.Image] = grypeRes
			failure = true
		}
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
