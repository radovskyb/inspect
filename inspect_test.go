package inspect

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"testing"
)

const testFileSrc = `package testfile

import "fmt"

// I'm a comment for ExportedFunctionOne
func ExportedFunctionOne() string {
	return fmt.Sprint("ExportedFunctionOne")
}

// I'm a comment for UnexportedFunctionOne
func unexportedFunctionOne() string {
	return fmt.Sprint("ExportedFunctionOne")
}`

var fset = token.NewFileSet()
var file *ast.File

func init() {
	var err error
	file, err = parser.ParseFile(fset, "", testFileSrc, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestParseImports(t *testing.T) {
	imports := ParseImports(file)

	if len(imports) != 1 {
		t.Errorf("%d imports found. expected 0", len(imports))
	}

	if imports[0] != "fmt" {
		t.Errorf("expected imports[0] to be `fmt`, got %s", imports[0])
	}
}

func TestParseFunction(t *testing.T) {
	newfset := token.NewFileSet()

	newfile, err := parser.ParseFile(newfset, "", testFileSrc, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}

	funcs := []*Function{}

	bb := new(bytes.Buffer)
	for _, d := range newfile.Decls {
		if fnc, ok := d.(*ast.FuncDecl); ok {
			f := ParseFunction(fnc, fset, bb)
			if f != nil {
				funcs = append(funcs, f)
			}
		}
	}

	if funcs[0].Name != "ExportedFunctionOne" {
		t.Errorf("function name incorrect, expected ExportedFunctionOne, got %s", funcs[0].Name)
	}

	signature := "func ExportedFunctionOne() string"
	if funcs[0].Signature != signature {
		t.Errorf("function signature incorrect, expected %s, got %s",
			signature, funcs[0].Signature)
	}

	documentation := "I'm a comment for ExportedFunctionOne"
	if funcs[0].Documentation != documentation {
		t.Errorf("function documentation incorrect, expected %s, got %s",
			documentation, funcs[0].Documentation)
	}
}

func TestParseFile(t *testing.T) {
	functions := ParseFile(file, fset)

	// Should only find exported functions.
	if len(functions) > 1 {
		t.Errorf("expected to find 1 function, found %d", len(functions))
	}

	if functions[0].Name != "ExportedFunctionOne" {
		t.Errorf("function name incorrect, expected ExportedFunctionOne, got %s",
			functions[0].Name)
	}

	signature := "func ExportedFunctionOne() string"
	if functions[0].Signature != signature {
		t.Errorf("function signature incorrect, expected %s, got %s",
			signature, functions[0].Signature)
	}

	documentation := "I'm a comment for ExportedFunctionOne"
	if functions[0].Documentation != documentation {
		t.Errorf("function documentation incorrect, expected %s, got %s",
			documentation, functions[0].Documentation)
	}
}
