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
	if len(os.Args) != 2 { //nolint
		log.Fatal("Wrong argument amount")
	}
	err := runGenerator(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func runGenerator(modelsFile string) error {
	Types = make(map[string]string)
	Structs = make(map[string]*ast.TypeSpec)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, modelsFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("wrong argument: %w", err)
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
				_ = printer.Fprint(&typeName, fset, currType.Type)
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
		return fmt.Errorf("no validation tags were found")
	}

	err = parseStructs(fset)
	if err != nil {
		return err
	}

	path := pathRegexp.Split(modelsFile, -1)[0] + "model_validation.go"
	return generate(path)
}
