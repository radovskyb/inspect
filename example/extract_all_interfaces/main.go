package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/radovskyb/inspect"
)

func main() {
	// Find the current Go intallation.
	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		log.Fatalln(err)
	}

	// Get the directory of the Go standard package library.
	pkgsRoot := filepath.Join(strings.TrimSpace(string(goroot)), "src")

	// Parse all Go packages, ignoring all test files and unexported functions.
	pkgs, err := inspect.ParsePackagesFromDir(pkgsRoot, true, inspect.FuncExported)
	if err != nil {
		log.Fatalln(err)
	}

	// Delete any non-library, main package's.
	delete(pkgs, "main")

	// Package name to interface.
	ifaces := map[string][]*inspect.Interface{}

	for _, pkg := range pkgs {
		if len(pkg.Interfaces) > 0 {
			ifaces[pkg.Name] = pkg.Interfaces
		}
	}

	for pkg, pkgIfaces := range ifaces {
		fmt.Printf("\nPackage %s:\n", pkg)
		fmt.Printf("%s\n", strings.Repeat("=", len("Package")+len(pkg)+2))
		for _, iface := range pkgIfaces {
			fmt.Printf("\n\tInterface %s\n", iface.Name)
			fmt.Printf("\t%s\n", strings.Repeat("-", len("Interface")+len(iface.Name)+1))
			if len(iface.Interfaces) > 0 {
				fmt.Println("\tImplements:")
				for _, ifc := range iface.Interfaces {
					fmt.Printf("\t\t%s\n", ifc)
				}
			}
			if len(iface.Methods) > 0 {
				fmt.Println("\tMethods:")
				for _, mth := range iface.Methods {
					fmt.Printf("\t\t%s\n", mth)
				}
			}
			fmt.Println()
		}
	}
}
