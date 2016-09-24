package testfiles

import (
	"fmt"
	"strings"
)

// I'm a comment for ExportedFunctionThree
func ExportedFunctionThree() string {
	fmt.Println(strings.Compare("a", "b"))
	return fmt.Sprint("ExportedFunctionThree")
}
