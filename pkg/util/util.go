package util

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
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
func FatalErrorCheck(err error) {
	if err != nil {
		log.Print(err)
		log.Print("\nERROR: " + err.Error() + "\n")
		log.Fatal(err)
	}
}

// NonfatalErrorCheck Helper method to output present errors but not crash
func NonfatalErrorCheck(err error) {
	if err != nil {
		//log.Print(err)
		log.Print("\nERROR Nonfatal: " + err.Error() + "\n")
	}
}

// ReadFile Helper method to read a file and return it as a string
func ReadFile(infile string) string {
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

// WriteFile Helper method to write a file
func WriteFile(outfile, data string) {
	// Create the file
	f, err := os.Create(outfile)
	// Crash if error
	FatalErrorCheck(err)
	// Close the file with defer
	defer f.Close()

	// Write the data
	f.WriteString(data)
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
