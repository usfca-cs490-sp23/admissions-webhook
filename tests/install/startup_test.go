package tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

/* Test case to check that all external dependencies are present on the system */
func TestDependencies(t *testing.T) {
	// Array of dependencies
	dependencies := []string{"kind", "kubectl", "syft", "grype", "openssl", "docker", "xclip", "pbcopy"}

	// Loop through every dependency in array
	for _, dependency := range dependencies {
		// Construct a string that is returned from calling an unknown package in a shell
		invalidString := fmt.Sprintf("%s not found", dependency)
		// Run a command with just the package and -h flag
		output, err := exec.Command("which", dependency).Output()
		// If there is an error, print it to stderr and fail the test
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			t.Errorf("Error while running TestDependencies on: %s", dependency)
		}

		// If the command output contains "[package]: command not found", fail the test
		if string(output) == invalidString {
			t.Errorf("Error while running TestDependencies: %s not installed or present on path", dependency)
		}
	}
}
