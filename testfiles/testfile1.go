package testfiles

import (
	"fmt"
	"reflect"
)

// I'm a comment for ExportedFunctionOne
func ExportedFunctionOne() string {
	fmt.Println(reflect.TypeOf(0))
	return fmt.Sprint("ExportedFunctionOne")
}

// I'm a comment for unexportedFunctionOne
func unexportedFunctionOne() string {
	return fmt.Sprint("unexportedFunctionOne")
}
