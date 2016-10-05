# inspect
`inspect` is a small Go package that makes it easy to inspect information about `Go` source code.

# Installation

```shell
go get github.com/radovskyb/inspect
```

# Example

#### Encode all standard library package's information into a single json file.

```go
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
	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		log.Fatalln(err)
	}

	// Get the directory of the Go standard package library.
	pkgsRoot := filepath.Join(strings.TrimSpace(string(goroot)), "src")

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
```
