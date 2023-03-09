package util

import (
    "fmt"
    "os"
    "flag"
    "math/rand"
    "time"
)

func Usage() {
    // Print message to stdout
    fmt.Println(UsageString())
    // Terminate program processing
    os.Exit(0)
}

/* Method to print out correct usage method and exit the prgram */
func UsageString() string {
    return ("\nUsage: go run main.go COMMAND\n\nA program to build a kubernetes " +
    "cluster using kind, and then apply a security admission webhook\n\nCommands:\n" +
    "\t-c [config_file_path] \tCreate cluster from a config file\n" +
    "\t-cluster\t\tCreate a kind cluster\n\t-info\t\t\tGet cluster information\n" +
    "\t-interface\t\tOpen Kubernetes web interface\n\t-deploy\t\t\tDeploy admissions " +
    "controller webhook\n\t-shutdown [name]\tShutdown a kubernetes cluster\n\n" +
    "For more information, please visit https://github.com/usfca-cs490-sp23/admissions-webhook")
}

/* Helper method to panic and trace to source method for unimplemented code */
func NotYetImplemented(method string) {
    panic((method + " not yet implemented"))
}

func IsFlagRaised(flag_name string) bool {
    found := false
    flag.Visit(func(f *flag.Flag) {
        if f.Name == flag_name {
            found = true
        }
    })
    return found
}

/* Helper method to generate a random string */
func RandomString (length int) string {
    // Generate seed
    rand.Seed(time.Now().UnixNano())
    // Make a byte array with length + 2
    byteArr := make([]byte, length+2)
    // Write random bytes into the byte array
    rand.Read(byteArr)
    // Return the random string
    return fmt.Sprintf("%x", byteArr)[2 : length+2]
}

