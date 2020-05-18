package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Identify tags
	lenRegexp    = regexp.MustCompile(`len:[0-9]+`)
	regexpRegexp = regexp.MustCompile(`regexp:[^\"]+`)
	minRegexp    = regexp.MustCompile(`min:[0-9]+`)
	maxRegexp    = regexp.MustCompile(`max:[0-9]+`)
	inRegexp     = regexp.MustCompile(`in:[^\"]+`)

	// Validation tags
	fRegexp string
	fLen    string
	fMax    string
	fMin    string
	fIn     string
)

func clean() {
	fRegexp = ""
	fLen = ""
	fMax = ""
	fMin = ""
	fIn = ""
}

func setInTag(tag string) {
	temp := strings.Split(strings.Split(inRegexp.FindStringSubmatch(tag)[0], ":")[1], ",")
	// Store strings with quotes and numbers without quotes
	isString := false
	for i, t := range temp {
		_, err := strconv.Atoi(t)
		if (err != nil) || isString {
			isString = true
			temp[i] = fmt.Sprintf("%q", temp[i])
		}
	}
	fIn = strings.Join(temp, ",")
}

func parseTag(tag string, fieldType string) {
	switch {
	// If string, check length, regexp and variants
	case (fieldType == "string") || (fieldType == "[]string"):
		switch {
		// Length
		case len(lenRegexp.FindStringSubmatch(tag)) != 0:
			fLen = strings.Split(lenRegexp.FindStringSubmatch(tag)[0], ":")[1]
			_, err := strconv.Atoi(fLen)
			if err != nil {
				log.Fatal("Len field should be int\n")
			}
		// Regexp
		case len(regexpRegexp.FindStringSubmatch(tag)) != 0:
			fRegexp = strings.Split(regexpRegexp.FindStringSubmatch(tag)[0], ":")[1]
			fRegexp = fmt.Sprintf(`%q`, fRegexp)
		// Variants
		case len(inRegexp.FindStringSubmatch(tag)) != 0:
			setInTag(tag)
		}
	// If int, check min/max and variants
	case (fieldType == "int") || (fieldType == "[]int"):
		switch {
		// Max and/or min
		case len(maxRegexp.FindStringSubmatch(tag)) != 0:
			fMax = strings.Split(maxRegexp.FindStringSubmatch(tag)[0], ":")[1]
			_, err := strconv.Atoi(fMax)
			if err != nil {
				log.Fatal("Max field should be int\n")
			}

			if min := minRegexp.FindStringSubmatch(tag); len(min) != 0 {
				fMin = strings.Split(min[0], ":")[1]
				_, err := strconv.Atoi(fMin)
				if err != nil {
					log.Fatal("Min field should be int\n")
				}
			}
		// Min
		case len(minRegexp.FindStringSubmatch(tag)) != 0:
			fMin = strings.Split(minRegexp.FindStringSubmatch(tag)[0], ":")[1]
			_, err := strconv.Atoi(fMin)
			if err != nil {
				log.Fatal("Min field should be int\n")
			}
		// Variants
		case len(inRegexp.FindStringSubmatch(tag)) != 0:
			setInTag(tag)
		}
	default:
		log.Fatalf("Unsupported type %s", fieldType)
	}
}

func parseStructs(fset *token.FileSet) {
	structsToValidate = []TemplateStruct{}
	for n, s := range Structs {
		tempFields := []TemplateField{}
		fields := s.Type.(*ast.StructType).Fields.List

		for _, f := range fields {
			var typeName bytes.Buffer
			_ = printer.Fprint(&typeName, fset, f.Type)
			fieldType, err := normalizeType(typeName.String())
			if err != nil {
				log.Fatal(err)
			}

			if f.Tag == nil {
				clean()
			} else {
				parseTag(f.Tag.Value, fieldType)
			}

			tempFields = append(tempFields, TemplateField{
				Name:   f.Names[0].String(),
				Type:   fieldType,
				Regexp: fRegexp,
				Len:    fLen,
				Max:    fMax,
				Min:    fMin,
				In:     fIn,
			})
			clean()
		}

		structsToValidate = append(structsToValidate, TemplateStruct{
			Name:   n,
			Fields: tempFields,
		})
	}
}
