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

// A FuncOption is used to decide whether exported and/or
// unexported functions get parsed.
type FuncOption int

const (
	FuncUnexported FuncOption = 1 << iota
	FuncExported

	FuncBoth = FuncExported | FuncUnexported
)

// A Package describes a package.
//
// A Package contains a package name, a slice of the package's imports,
// and also a slice of all of the Function's that the package contains.
type Package struct {
	Name    string      `json:"-"`
	Imports []string    `json:",omitempty"`
	Funcs   []*Function `json:",omitempty"`
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

// IsExported is a wrapper around ast.IsExported that returns a true or false
// value based on whether the current function is exported or not.
func (f *Function) IsExported() bool {
	return ast.IsExported(f.Name)
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
// parsed packages and the error that occurred.
func ParsePackagesFromDir(dir string, ignoreTests bool, funcOption FuncOption) (map[string]*Package, error) {
	fset := token.NewFileSet()

	pkgs := make(map[string]*Package)

	var filter func(os.FileInfo) bool
	if ignoreTests {
		filter = FilterIgnoreTests
	}

	return pkgs, filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || strings.HasPrefix(path, filepath.Join(dir, "cmd")) {
			return nil
		}

		parsed, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, pkg := range parsed {
			p := ParsePackage(fset, pkg, funcOption)
			if _, exists := pkgs[pkg.Name]; exists {
				pkgs[pkg.Name].Funcs = append(pkgs[pkg.Name].Funcs, p.Funcs...)
			} else {
				pkgs[pkg.Name] = p
			}
		}

		return nil
	})
}

// ParsePackage returns a *Package generated from an *ast.Package.
func ParsePackage(fset *token.FileSet, pkg *ast.Package, funcOption FuncOption) *Package {
	// Merge all of the package's files into a single file, and filter
	// out any import or function duplicates along the way.
	mergedFile := ast.MergePackageFiles(pkg,
		ast.FilterFuncDuplicates+ast.FilterImportDuplicates,
	)

	// Return a new Package with it's fields appropriately set.
	return &Package{
		Name:    pkg.Name,
		Funcs:   ParseFileFuncs(fset, mergedFile, funcOption),
		Imports: ParseFileImports(mergedFile),
	}
}

// ParseFileFuncs returns a []*Function's generated from an *ast.File.
func ParseFileFuncs(fset *token.FileSet, file *ast.File, funcOption FuncOption) []*Function {
	funcs := []*Function{}

	// If funcOption isn't set then parse both exported and unexported functions.
	if funcOption&FuncUnexported == 0 && funcOption&FuncExported == 0 {
		funcOption = FuncBoth
	}

	bb := new(bytes.Buffer)
	ast.Inspect(file, func(n ast.Node) bool {
		bb.Reset()
		if fnc, ok := n.(*ast.FuncDecl); ok {
			var f *Function
			if funcOption&FuncUnexported != 0 && !fnc.Name.IsExported() {
				f = ParseFunction(fset, fnc, bb)
			}
			if funcOption&FuncExported != 0 && fnc.Name.IsExported() {
				f = ParseFunction(fset, fnc, bb)
			}
			if f != nil {
				funcs = append(funcs, f)
			}
		}
		return true
	})

	return funcs
}

// ParseFunction returns a *Function's generated from an *ast.FuncDecl.
func ParseFunction(fset *token.FileSet, fnc *ast.FuncDecl, bb *bytes.Buffer) *Function {
	f := &Function{Name: fnc.Name.Name}

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

// ParseFileImports generates a list of imports from an *ast.File object.
func ParseFileImports(file *ast.File) []string {
	imports := []string{}

	// Append the file's imports to the imports string slice.
	for _, i := range file.Imports {
		imports = append(imports, strings.Trim(i.Path.Value, "\""))
	}

	return imports
}
