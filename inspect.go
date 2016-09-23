package inspect

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ErrParseDir is an error for when there is an error parsing a directory.
var ErrParseDir = errors.New("error: parsing directory")

// FilterIgnoreTests is a filter function for parser.ParseDir
// that ignores all test files.
var FilterIgnoreTests = func(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

// A Package describes a package.
//
// A Package contains a package name and a slice
// of all of the Function's that the package contains.
type Package struct {
	Name  string      `json:"-"`
	Funcs []*Function `json:",omitempty"`
}

// A Function describes a function.
//
// A Function contains a function name, function signature
// and also the function's documentation.
type Function struct {
	Name          string `json:"Name"`
	Signature     string `json:"Sig"`
	Documentation string `json:"Doc,omitempty"`
}

// ParsePackagesFromDir parses all packages in a directory.
//
// If ignoreTests is true, all test files will be ignored.
//
// If there are directories contained within dir, ParsePackagesFromDir
// attempts to traverse into those directories as well.
//
// If an error occurs whilst traversing the nested directories,
// ParsePackagesFromDir will return a map containing any correctly
// passed packages and the error that occured.
func ParsePackagesFromDir(dir string, ignoreTests bool) (map[string]*Package, error) {
	fset := token.NewFileSet()

	pkgs := make(map[string]*Package)

	var filter func(os.FileInfo) bool
	if ignoreTests {
		filter = FilterIgnoreTests
	}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() || strings.HasPrefix(path, filepath.Join(dir, "cmd")) {
			return nil
		}

		parsed, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, pkg := range parsed {
			p := ParsePackage(pkg, fset)
			if _, exists := pkgs[pkg.Name]; exists {
				pkgs[pkg.Name].Funcs = append(pkgs[pkg.Name].Funcs, p.Funcs...)
			} else {
				pkgs[pkg.Name] = p
			}
		}

		return nil
	}); err != nil {
		return pkgs, err
	}

	return pkgs, nil
}

// ParsePackage returns a *Package generated from an *ast.Package.
func ParsePackage(pkg *ast.Package, fset *token.FileSet) *Package {
	p := &Package{Name: pkg.Name, Funcs: []*Function{}}

	bb := new(bytes.Buffer)
	ast.Inspect(pkg, func(n ast.Node) bool {
		bb.Reset()
		if fnc, ok := n.(*ast.FuncDecl); ok {
			f := ParseFunction(fnc, fset, bb)
			if f != nil {
				p.Funcs = append(p.Funcs, f)
			}
		}
		return true
	})

	return p
}

// ParseFile returns a []*Function's generated from an *ast.File.
func ParseFile(file *ast.File, fset *token.FileSet) []*Function {
	funcs := []*Function{}

	bb := new(bytes.Buffer)
	ast.Inspect(file, func(n ast.Node) bool {
		bb.Reset()
		if fnc, ok := n.(*ast.FuncDecl); ok {
			f := ParseFunction(fnc, fset, bb)
			if f != nil {
				funcs = append(funcs, f)
			}
		}
		return true
	})

	return funcs
}

// ParseFunction returns a *Function's generated from an *ast.FuncDecl.
func ParseFunction(fnc *ast.FuncDecl, fset *token.FileSet, bb *bytes.Buffer) *Function {
	f := &Function{Name: fnc.Name.Name}

	// Skip the function if it's unexported.
	if !fnc.Name.IsExported() {
		return nil
	}

	fnc.Body = nil

	if err := printer.Fprint(bb, fset, fnc); err != nil {
		return nil
	}

	var startPos int
	if fnc.Doc.Text() != "" {
		startPos = int(fnc.Type.Pos() - fnc.Doc.Pos())
	}

	f.Signature = bb.String()[startPos:]
	f.Documentation = strings.TrimSpace(fnc.Doc.Text())

	return f
}

// ParseImports generates a list of imports from an *ast.File object.
func ParseImports(file *ast.File) []string {
	imports := []string{}

	// Append the file's imports to the imports string slice.
	for _, i := range file.Imports {
		imports = append(imports, strings.Trim(i.Path.Value, "\""))
	}

	return imports
}
