package util

import (
    "fmt"
    "os"
    "flag"
    "math/rand"
    "time"
    "log"
    "io"
    "strings"
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

/* Helper method to check if a flag has been raised */
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

/* Helper method to crash if errors exist */
func FatalErrorCheck(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

/* Helper method to output present errors but not crash */
func NonfatalErrorCheck(err error) {
    if err != nil {
        log.Print(err)
    }
}

/* Helper method to read a file and return it as a string */
func readFile (infile string) string {
    // Open the current file and generate reader
    f, err := os.Open(infile)
    FatalErrorCheck(err)
    defer f.Close()

    // Read the current file
    content, err := io.ReadAll(f)
    // Crash if error
    FatalErrorCheck(err) 

    return string(content)
}

/* Helper method to write a file */
func writeFile (outfile, data string) {
    // Create the file 
    f, err := os.Create(outfile)
    // Crash if error
    FatalErrorCheck(err)
    // Close the file with defer
    defer f.Close()

    // Write the data
    f.WriteString(data)
}

/* Method to inject a CA bundle into a YAML file */
func InjectYamlCA (target, template, injectable string) {
    // Read and store file with tls data
    content := readFile(injectable)
    // Read and store file to inject tls into
    config := readFile(template)

    // Remove unnecessary prefix from content
    content = strings.TrimPrefix(content, ">> MutatingWebhookConfiguration caBundle:")
    // Indent content correctly
    content = strings.ReplaceAll(content, "\n", "\n                ")

    // Insert the content into the config file string
    config = strings.Replace(config, "caBundle: |", "caBundle: |\n                " + content, 1)

    // Now write to file
    writeFile(target, config)
}

