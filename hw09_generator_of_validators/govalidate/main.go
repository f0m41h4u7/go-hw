package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"regexp"
)

var (
	// Types are non-struct user-defined types
	Types   map[string]string
	Structs map[string]*ast.TypeSpec

	structsToValidate []TemplateStruct

	needsValidation = regexp.MustCompile(`.*validate:"(len|regexp|min|max|in):.+"`)
	pathRegexp      = regexp.MustCompile(`[^\/]+\.go`)
)

func main() {
	Types = make(map[string]string)
	Structs = make(map[string]*ast.TypeSpec)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(fmt.Errorf("wrong argument: %w", err))
	}

	for _, f := range node.Decls {
		genD, ok := f.(*ast.GenDecl)
		if !ok {
			fmt.Printf("SKIP %T is not *ast.GenDecl\n", f)
			continue
		}
		for _, spec := range genD.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				fmt.Printf("SKIP %T is not *ast.TypeSpec\n", spec)
				continue
			}

			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				var typeName bytes.Buffer
				err := printer.Fprint(&typeName, fset, currType.Type)
				if err != nil {
					log.Fatalf("failed printing %s", err)
				}
				Types[currType.Name.String()] = typeName.String()
				continue
			}

			fields := currStruct.Fields.List

			for _, f := range fields {
				var tag string
				if f.Tag == nil {
					tag = ""
				} else {
					tag = f.Tag.Value
				}

				valid := needsValidation.FindStringSubmatch(tag)
				if len(valid) != 0 {
					Structs[currType.Name.String()] = currType
				}
			}
		}
	}

	if len(Structs) == 0 {
		log.Fatalf("No validation tags were found\n")
	}

	parseStructs(fset)

	path := pathRegexp.Split(os.Args[1], -1)[0] + "model_validation.go"
	generate(path)
}
