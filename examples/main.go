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
			fmt.Println(f)
		}
	}
	fmt.Println("--------------------\n")

	// Print all of the file's exported function's names.
	fmt.Println("Exported functions:")
	fmt.Println("--------------------")
	for _, f := range file.Functions {
		if f.IsExported() {
			fmt.Println(f)
		}
	}
}
