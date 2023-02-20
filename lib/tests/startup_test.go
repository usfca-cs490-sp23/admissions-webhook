package tests

import (
    "fmt"
    "testing"
    "os"
    "os/exec"
    "strings"
    "github.com/usfca-cs490/admissions-webhook/lib/util"
)

/* Test case for no arguments */
func TestEmptyArgs(t *testing.T) {
    // Grab the expected message from util
    expected := util.UsageString()
    // Execute main.go with no argument
    actual, err := exec.Command("go", "run", "../../main.go").Output()

    // If there is an error, print it out to stderr and fail the test
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s", err)
        t.Errorf("Error while running TestEmptyArgs")
    } else if strings.Trim(string(actual), "\n") !=         /* If the expected and actual strings do not */
                strings.Trim(expected, "\n") {              /*  match, fail the test */
        t.Errorf("Actual did not match expected!")
    }
}

/* Test case to check that all external dependencies are present on the system */
func TestDependencies(t * testing.T) {
    // Array of dependencies
    dependencies := []string{"kind", "kubectl", "syft"}

    // Loop through every dependency in array
    for _, ele := range dependencies {
        // Construct a string that is returned from calling an unknown package in a shell
        invalidSubstr := fmt.Sprintf("%s: command not found", ele)
        // Run a command with just the package and -h flag
        output, err := exec.Command(ele, "-h").Output()
        // If there is an error, print it to stderr and fail the test
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s", err)
            t.Errorf("Error while running TestDependencies on: %s", ele)
        }

        // If the command output contains "[package]: command not found", fail the test
        if strings.Contains(string(output), invalidSubstr) {
            t.Errorf("Error while running TestDependencies: %s not installed or present on path", ele)
        }
    }
}

/* Test case to check container manager is present (docker or podman) */
func TestContainerProvider(t * testing.T) {
    // Array of two valid container managers (NOTE: podmans is experimental with kind)
    dependencies := []string{"docker", "podman"}

    // Loop through every dependency in array
    for _, ele := range dependencies {
        // Construct a string that is returned from calling an unknown package in a shell
        invalidSubstr := fmt.Sprintf("%s: command not found", ele)
        // Run a command with just the package and -h flag
        output, _ := exec.Command(ele, "-h").Output()

        // If the command output does not contain "[package]: command not found", then pass the test
        if !strings.Contains(string(output), invalidSubstr) {
            return
        }
    }
    // Fail test case if neither provider was found
    t.Errorf("Error while running TestContainer: neither docker nor podman installed or present on path")
}

