package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/radovskyb/inspect"
)

func main() {
	// Find the current Go intallation.
	gobin, err := exec.LookPath("go")
	if err != nil {
		log.Fatalln(err)
	}

	// Get the path of the root of the Go standard library packages.
	pkgsRoot := filepath.Join(strings.TrimSuffix(gobin, filepath.Join("bin", "go")), "src")

	// Parse all Go packages, ignoring all test files and unexported functions.
	pkgs, err := inspect.ParsePackagesFromDir(pkgsRoot, true, inspect.FuncExported)
	if err != nil {
		log.Fatalln(err)
	}

	// Delete any non-library, main package's.
	delete(pkgs, "main")

	// Create a new json file to store all of Go's standard package library info.
	jsonFile, err := os.Create("packages.json")
	if err != nil {
		log.Fatalln(err)
	}

	// Create a new json encoder that writes to jsonFile and set it's
	// indentation formatting to a single tab.
	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "\t")

	// Encode all of the Package's to the json file.
	if err := encoder.Encode(pkgs); err != nil {
		log.Fatalln(err)
	}
}
