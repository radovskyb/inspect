# inspect
`inspect` is a small Go library that makes it easy to inspect information about `Go` source code.

`inspect` is a small Go library which wraps functionality from standard `Go` packages such as `go/ast`, `go/parser` and `go/token`, to make it easy to inspect information about `Go` source code.

### Example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/radovskyb/inspect"
)

func main() {
	// The filepath of the file to inspect.
	filename := "/usr/local/go/src/net/http/client.go"

	// Create a new inspect.File object from the file located at filename.
	file, err := inspect.NewFile(filename, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the file's package name.
	fmt.Printf("Package %q:\n", file.Package)
	fmt.Println("--------------------\n")

	// Print the file's imports.
	fmt.Println("Imports:")
	fmt.Println("--------------------")
	for _, i := range file.Imports {
		fmt.Println(i)
	}
	fmt.Println("--------------------\n")

	// Print all of the file's unexported function's names.
	fmt.Println("Unexported functions:")
	fmt.Println("--------------------")
	for _, f := range file.Functions {
		if !f.IsExported() {
			fmt.Println(f.Name)
		}
	}
	fmt.Println("--------------------\n")

	// Print all of the file's exported function's names.
	fmt.Println("Exported functions:")
	fmt.Println("--------------------")
	for _, f := range file.Functions {
		if f.IsExported() {
			fmt.Println(f.Name)
		}
	}
	fmt.Println("--------------------\n")

	// Print the documentation for the function named `Get`.
	fmt.Println(file.Functions["Get"].Doc)
}
```
