package inspect

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"strings"
)

// A Function contains a function name and it's documentation text.
type Function struct {
	Package   string
	Name      string
	Doc       string
	Signature string
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
	if src == nil {
		slurp, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		src = bytes.NewReader(slurp)
	}

	// Parse the Go source from either filename or src.
	parsed, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Return a new File with it's fields set appropriately.
	functions, err := InspectFunctions(parsed, src.(io.ReadSeeker))
	if err != nil {
		return nil, err
	}

	return &File{
		File:      parsed,
		Imports:   InspectImports(parsed), // Get the file's imports.
		Functions: functions,              // Get the file's functions.
	}, nil
}

// InspectFunctions generates a Functions map from an *ast.File object and an
// io.ReadSeaker containing the file's bytes.
func InspectFunctions(file *ast.File, fileReader io.ReadSeeker) (map[string]*Function, error) {
	functions := make(map[string]*Function)

	var bb = new(bytes.Buffer)
	for _, d := range file.Decls {
		if f, okay := d.(*ast.FuncDecl); okay {
			if f.Body != nil {
				sigStart := int64(f.Pos() - 1)
				sigEnd := int64(f.Body.Lbrace - 2)

				toRead := sigEnd - sigStart

				// Go to the start of the function declaration.
				_, err := fileReader.Seek(sigStart, io.SeekStart)
				if err != nil {
					return nil, err
				}

				// Make sure bb is empty.
				bb.Reset()

				// Read toRead number of bytes from fileReader to bb.
				_, err = io.CopyN(bb, fileReader, toRead)
				if err != nil && err != io.EOF {
					return nil, err
				}

				functions[f.Name.String()] = &Function{
					file.Name.String(),
					f.Name.String(),
					strings.TrimSpace(f.Doc.Text()),
					bb.String(),
				}
			}
		}
	}

	return functions, nil
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
