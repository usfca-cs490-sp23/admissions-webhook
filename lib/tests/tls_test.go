package tests

import (
    "testing"
	"github.com/usfca-cs490/admissions-webhook/lib/keygen"
    "io/ioutil"
    "log"
    "os"
    "regexp"
)

var (
    // List of files to test
    test_files = []string {
                    "./test_pem_cert.txt",
                    "./test_pem_key.txt",
                    "./test_ca_bundle.txt",
                }
)

/* Helper method to handle errors if they exist */
func errorCheck(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

/* Method to create a file from the list */
func createFiles() {
    // Generate PEM data
    pemCert, pemKey, caBundle := keygen.CreatePEMs()

    // Store PEM data and file to write to after encoding in a map
    data := map[string][]byte {
        test_files[0]: pemCert,
        test_files[1]: pemKey,
        test_files[2]: caBundle,
    }

    // Convert from PEM to base64 and write to files in data map
    keygen.ConvertPEMToB64(data)
}

/* Method to delete all of the files created for testing */
func deleteFiles() {
    for _, file := range test_files {
        err := os.Remove(file)
        errorCheck(err) 
    }
}

/* Test case for generating and encoding TLS certs, keys, and ca bundle */
func TestTLSFiles(t *testing.T) {
    // Create the files
    createFiles()
    // Loop through the files
    for _, file := range test_files {
        // Read the files
        content, err := ioutil.ReadFile(file)
        // Handle any file reading related errors    
        errorCheck(err) 

        // Evaluate file content with regex
        match, err := regexp.Match("^([a-zA-Z0-9+/=]*)$", content)
        // Handle any regex related errors
        errorCheck(err) 

        // If the data did nost match regex or file is empty, fail test case
        if !match || len(content) < 2 {
            t.Errorf("Error while running TestFiles: " + file + " is not in base64 format")
        }
    }
    // Delete test files
    deleteFiles()
}
