package util

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// NotYetImplemented Helper method to panic and trace to source method for unimplemented code
func NotYetImplemented(method string) {
	panic((method + " not yet implemented"))
}

// IsFlagRaised Helper method to check if a flag has been raised
func IsFlagRaised(flag_name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flag_name {
			found = true
		}
	})
	return found
}

// FatalErrorCheck Helper method to crash if errors exist
func FatalErrorCheck(err error, verbose bool) {
	if err != nil {
		if verbose {
			log.Print("\nERROR Fatal: " + err.Error() + "\n")
		} else {
			log.Print(err)
		}
		log.Fatal(err)
	}
}

// NonfatalErrorCheck Helper method to output present errors but not crash
func NonfatalErrorCheck(err error, verbose bool) {
	if err != nil {
		if verbose {
			log.Print("\nERROR Nonfatal: " + err.Error() + "\n")
		} else {
			log.Print(err)
		}
	}
}

// ReadFile Helper method to read a file and return it as a string
func ReadFile(infile string) string {
	// Open the current file and generate reader
	f, err := os.Open(infile)
	FatalErrorCheck(err, true)
	defer f.Close()

	// Read the current file
	content, err := io.ReadAll(f)
	// Crash if error
	FatalErrorCheck(err, true)

	return string(content)
}

// WriteFile Helper method to write a file
func WriteFile(outfile, data string) {
	// Create the file
	f, err := os.Create(outfile)
	// Crash if error
	FatalErrorCheck(err, true)
	// Close the file with defer
	defer f.Close()

	// Write the data
	f.WriteString(data)
}

// FormatTime generates a string of the current time for naming files
func FormatTime() string {
	// make filename based off of current time (lazy but effective)
	currTime := time.Now()
	// format is from the docs, plz don't change it bc the date is based on the underlying schema used by the Format()
	fileName := currTime.Format("2006-1-2_15-4-5")
	return fileName
}

// InjectYamlCA Method to inject a CA bundle into a YAML file
func InjectYamlCA(target, template, injectable string) {
	// Read and store file with tls data
	content := ReadFile(injectable)
	// Read and store file to inject tls into
	config := ReadFile(template)

	// Remove unnecessary prefix from content
	content = strings.TrimPrefix(content, ">> MutatingWebhookConfiguration caBundle:")
	// Indent content correctly
	content = strings.ReplaceAll(content, "\n", "\n                ")

	// Insert the content into the config file string
	config = strings.Replace(config, "caBundle: |", "caBundle: |\n                "+content, 1)

	// Now write to file
	WriteFile(target, config)
}

func WriteEvent(kind string, eventtype, reason, message string) {
	event := "kubectl apply -f - <<EOF\napiVersion: v1\nkind: Event\nmetadata:\n  name: " + reason + "\n  namespace: default\ntype: " + eventtype + "\nmessage: '" + message + "'\ninvolvedObject:\n  kind: " + kind + "\nEOF"
	exec.Command("bash", "-c", event).Run()
}
