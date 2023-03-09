package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
    "github.com/usfca-cs490/admissions-webhook/pkg/dashboard"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/* Build the webhook */
func Build() {
	// handle our core application
	http.HandleFunc("/validate-pods", ValidatePod)

	// start the server
	// listens to clear text http on port 8080 unless TLS env var is set to "true" (which it should be)
	if os.Getenv("TLS") == "true" {
		// TODO: see if these are based off of the k8s admin setup or something else
		cert := "/etc/admission-webhook/tls/tls.crt"
		key := "/etc/admission-webhook/tls/tls.key"

		// it should ALWAYS be true for us I believe
		logrus.Fatal(http.ListenAndServeTLS(":443", cert, key, nil))
	} else {
		logrus.Print("Listening on port 8080...")
		logrus.Fatal(http.ListenAndServe(":8080", nil))
	}
}

// ValidatePod validates an admission request and then writes an admission review to response writer
func ValidatePod(w http.ResponseWriter, r *http.Request) {
	// special logging stuff
	logger := logrus.WithField("uri", r.RequestURI)
	logger.Debug("received validation request")

	// get the information from the request
	in, err := reviewAdmission(*r)
	// if there was an error, handle it here
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: use DashUpdate for passing to whatever code will eventually get written to show results on the K8s Dashboard
	// out, dashUpdate, err := ValidatePodReview(in.Request)
	// TODO: use above line once k8s dashboard interaction is needed
	out, _, err := ValidatePodReview(in.Request)
	if err != nil {
		e := fmt.Sprintf("could not generate admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	// set the response's header type to json
	w.Header().Set("Content-Type", "application/json")
	// takes the admision review struct and turns it into json
	jout, err := json.Marshal(out)
	// if Marshal() fails, log the error and exits the function
	if err != nil {
		e := fmt.Sprintf("could not parse admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	logger.Debug("sending response")
	logger.Debugf("%s", jout)
	// same as sprinf in C
	fmt.Fprintf(w, "%s", jout)
}

// Either returns an admission review struct or an error
// parseRequest extracts an AdmissionReview from an http.Request if possible
func reviewAdmission(r http.Request) (*admissionv1.AdmissionReview, error) {
	// Check if the given content is JSON, and if not, then return nil and log what kind of content was given
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}

	// create a buffer to read all the data from the body of the request
	// then denote that buffer as byte data
	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	// if the body is empty, then return nil and log the error
	if len(body) == 0 {
		return nil, fmt.Errorf("admission request body is empty")
	}

	// create an AdmissionReview object to store the data in
	var a admissionv1.AdmissionReview

	// Unmarshal() takes a byte array (in this case body) and an empty pointer (a) to store the results into
	// if it cannot be unmarshalled then return nil and log the error
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	// if the AdmissionReview's request field is empty, retun nil and log the error
	if a.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}

	// return the AdmissionReview struct and nil for error
	return &a, nil
}

// ValidatePodReview Take a K8s admission request and return a review struct based on,
// whether or not it is accepted into the cluster
func ValidatePodReview(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, dashboard.DashboardUpdate, error) {
	pod, err := extractPod(request)
	if err != nil {
		//e := fmt.Sprintf("could not parse pod in admission review request: %v", err)
		return nil, dashboard.BadPodDashUpdate(), err
	}

	//v := validation.NewValidator(a.Logger)
	podDecision, err := checkPodImages(pod)
	if err != nil {
		//e := fmt.Sprintf("could not validate pod: %v", err)
		return nil, dashboard.BadPodDashUpdate(), err
	}

	// if the pod is scanned and allowed, then return this review struct
	if !podDecision.Denied {
		return &admissionv1.AdmissionReview{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AdmissionReview",
				APIVersion: "admission.k8s.io/v1",
			},
			Response: &admissionv1.AdmissionResponse{
				UID:     request.UID,
				Allowed: true,
				Result: &metav1.Status{
					Code:    http.StatusAccepted,
					Message: "Pod scanned and admitted",
				},
			},
		}, podDecision, nil
	} else {
		// if the pod is reviewed and disalloed CVEs are found, return this rejection review
		return &admissionv1.AdmissionReview{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AdmissionReview",
				APIVersion: "admission.k8s.io/v1",
			},
			Response: &admissionv1.AdmissionResponse{
				UID:     request.UID,
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusForbidden,
					Message: "Pod scanned and denied",
				},
			},
		}, podDecision, nil
	}
}

// extractPod given an admission request, extract and return a Pod
func extractPod(request *admissionv1.AdmissionRequest) (*corev1.Pod, error) {
	pod := corev1.Pod{}
	// if the pod ain't right, throw an error and return
	if err := json.Unmarshal(request.Object.Raw, &pod); err != nil {
		return nil, err
	}

	// otherwise return the
	return &pod, nil
}

