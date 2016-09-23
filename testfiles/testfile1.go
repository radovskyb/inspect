package testfiles

import "fmt"

// I'm a comment for ExportedFunctionOne
func ExportedFunctionOne() string {
	return fmt.Sprint("ExportedFunctionOne")
}

// I'm a comment for unexportedFunctionOne
func unexportedFunctionOne() string {
	return fmt.Sprint("unexportedFunctionOne")
}
