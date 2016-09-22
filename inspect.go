package inspect

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

var fset = token.NewFileSet()

// A Function contains information about a function.
type Function struct {
	Package   string
	Name      string
	Signature string
	Doc       string
}

// String prints a function's information.
func (f Function) String() string {
	separator := strings.Repeat("-", len(f.Signature))
	str := f.Signature + "\n"
	if f.Doc != "" {
		str += separator + "\n" + f.Doc + "\n"
	}
	return str
}

// IsExported is a wrapper around ast.IsExported and returns a true or false
// value based on whether the current function is exported or not.
func (f Function) IsExported() bool {
	return ast.IsExported(f.Name)
}

// Functions is a map that stores function names as keys and Function
// objects as values.
type Functions map[string]*Function

// A File contains information about a Go source file. It is also a
// wrapper around *ast.File.
type File struct {
	*ast.File
	Imports []string
	Functions
}

// NewFile creates a File object from either a file, via the filename parameter,
// or source code, via the src parameter.
//
// If src != nil, NewFile parses the source from src.
func NewFile(filename string, src interface{}) (*File, error) {
	// Parse the Go source from either filename or src.
	parsed, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Return a new File with it's fields set appropriately.
	return &File{
		File:      parsed,
		Imports:   InspectImports(parsed),   // Get the file's imports.
		Functions: InspectFunctions(parsed), // Get the file's functions.
	}, nil
}

// InspectFunctions generates a Functions map from an *ast.File object.
func InspectFunctions(file *ast.File) map[string]*Function {
	functions := make(map[string]*Function)

	bb := new(bytes.Buffer)
	ast.Inspect(file, func(n ast.Node) bool {
		bb.Reset()
		if fnc, ok := n.(*ast.FuncDecl); ok {
			if fnc.Body == nil {
				return false
			}
			if err := printer.Fprint(bb, fset, fnc); err != nil {
				return false
			}

			var startPos, toRead int
			if fnc.Doc.Text() != "" {
				startPos = int(fnc.Type.Pos() - fnc.Doc.Pos())
			}
			toRead = int(fnc.Body.Pos() - fnc.Type.Pos() - 1)

			functions[fnc.Name.Name] = &Function{
				Package:   file.Name.Name,
				Name:      fnc.Name.Name,
				Signature: bb.String()[startPos : startPos+toRead],
				Doc:       strings.TrimSpace(fnc.Doc.Text()),
			}
		}
		return true
	})

	return functions
}

// InspectImports generates a list of imports from an *ast.File object.
func InspectImports(file *ast.File) []string {
	imports := []string{}

	// Append the file's imports to the imports string slice.
	for _, i := range file.Imports {
		imports = append(imports, strings.Trim(i.Path.Value, "\""))
	}

	return imports
}
