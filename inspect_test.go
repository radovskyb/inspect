package inspect

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"testing"
)

var fset = token.NewFileSet()
var file *ast.File

const (
	tf1FuncName = "ExportedFunctionOne"
	tf1FuncSig  = "func ExportedFunctionOne() string"
	tf1FuncDoc  = "I'm a comment for ExportedFunctionOne"

	tf1Path   = "testfiles/testfile1.go"
	tfPkgName = "testfiles"
)

func init() {
	var err error
	file, err = parser.ParseFile(fset, tf1Path, nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestParseFileImports(t *testing.T) {
	imports := ParseFileImports(file)

	if len(imports) != 2 {
		t.Errorf("%d imports found. expected 2", len(imports))
	}

	if imports[0] != "fmt" {
		t.Errorf("expected imports[0] to be `fmt`, got %s", imports[0])
	}
}

func TestParseFunction(t *testing.T) {
	funcs := []*Function{}

	bb := new(bytes.Buffer)
	for _, d := range file.Decls {
		if fnc, ok := d.(*ast.FuncDecl); ok {
			f := ParseFunction(fnc, fset, bb)
			if f != nil {
				funcs = append(funcs, f)
			}
		}
	}

	if funcs[0].Name != tf1FuncName {
		t.Errorf("function name incorrect, expected %s, got %s",
			tf1FuncName, funcs[0].Name)
	}

	if funcs[0].Signature != tf1FuncSig {
		t.Errorf("function signature incorrect, expected %s, got %s",
			tf1FuncSig, funcs[0].Signature)
	}

	if funcs[0].Documentation != tf1FuncDoc {
		t.Errorf("function documentation incorrect, expected %s, got %s",
			tf1FuncDoc, funcs[0].Documentation)
	}
}

func TestParseFile(t *testing.T) {
	funcs := ParseFile(file, fset)

	// Should only find exported functions.
	if len(funcs) > 1 {
		t.Errorf("expected to find 1 function, found %d", len(funcs))
	}

	if funcs[0].Name != tf1FuncName {
		t.Errorf("function name incorrect, expected %s, got %s",
			tf1FuncName, funcs[0].Name)
	}

	if funcs[0].Signature != tf1FuncSig {
		t.Errorf("function signature incorrect, expected %s, got %s",
			tf1FuncSig, funcs[0].Signature)
	}

	if funcs[0].Documentation != tf1FuncDoc {
		t.Errorf("function documentation incorrect, expected %s, got %s",
			tf1FuncDoc, funcs[0].Documentation)
	}
}

func TestParsePackage(t *testing.T) {
	pkgs, err := parser.ParseDir(fset, "testfiles", FilterIgnoreTests, parser.ParseComments)
	if err != nil {
		t.Error(err)
	}

	if len(pkgs) != 1 {
		t.Errorf("expected 1 package, found %d", len(pkgs))
	}

	if _, exists := pkgs[tfPkgName]; !exists {
		t.Errorf("package %s not found", tfPkgName)
	}

	if len(pkgs[tfPkgName].Files) != 2 {
		t.Errorf("expected 2 package files, found %d", len(pkgs[tfPkgName].Files))
	}

	pkg := ParsePackage(pkgs[tfPkgName], fset)

	if pkg.Name != tfPkgName {
		t.Errorf("expected package name %s, got %s", pkgs[tfPkgName].Name)
	}

	if len(pkg.Imports) != 3 {
		t.Errorf("expected to find 3 imports, found %d", len(pkg.Imports))
	}

	if len(pkg.Funcs) != 2 {
		t.Errorf("expected to find 2 functions, found %d", len(pkg.Funcs))
	}
}
