package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
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

func isNumerical(temp []string) bool {
	for _, t := range temp {
		if _, err := strconv.Atoi(t); err != nil {
			return false
		}
	}
	return true
}

func stringInQuotes(temp []string) string {
	for i, _ := range temp { //nolint
		temp[i] = fmt.Sprintf("%q", temp[i])
	}
	return strings.Join(temp, ",")
}

func parseTag(tag string, fieldType string) error {
	switch {
	// If string, check length, regexp and variants
	case (fieldType == "string") || (fieldType == "[]string"):
		switch {
		// Length
		case len(lenRegexp.FindStringSubmatch(tag)) != 0:
			fLen = strings.Split(lenRegexp.FindStringSubmatch(tag)[0], ":")[1]
			_, err := strconv.Atoi(fLen)
			if err != nil {
				return fmt.Errorf("len field should be int") //nolint
			}
		// Regexp
		case len(regexpRegexp.FindStringSubmatch(tag)) != 0:
			fRegexp = fmt.Sprintf(`%q`, strings.Split(regexpRegexp.FindStringSubmatch(tag)[0], ":")[1])
		// Variants
		case len(inRegexp.FindStringSubmatch(tag)) != 0:
			temp := strings.Split(strings.Split(inRegexp.FindStringSubmatch(tag)[0], ":")[1], ",")
			fIn = stringInQuotes(temp)
		}
	// If int, check min/max and variants
	case (fieldType == "int") || (fieldType == "[]int"):
		switch {
		// Max and/or min
		case len(maxRegexp.FindStringSubmatch(tag)) != 0:
			fMax = strings.Split(maxRegexp.FindStringSubmatch(tag)[0], ":")[1]
			if _, err := strconv.Atoi(fMax); err != nil {
				return fmt.Errorf("max field should be int") //nolint
			}

			if min := minRegexp.FindStringSubmatch(tag); len(min) != 0 {
				fMin = strings.Split(min[0], ":")[1]
				if _, err := strconv.Atoi(fMin); err != nil {
					return fmt.Errorf("min field should be int") //nolint
				}
			}
		// Min
		case len(minRegexp.FindStringSubmatch(tag)) != 0:
			fMin = strings.Split(minRegexp.FindStringSubmatch(tag)[0], ":")[1]
			if _, err := strconv.Atoi(fMin); err != nil {
				return fmt.Errorf("min field should be int") //nolint
			}
		// Variants
		case len(inRegexp.FindStringSubmatch(tag)) != 0:
			temp := strings.Split(strings.Split(inRegexp.FindStringSubmatch(tag)[0], ":")[1], ",")
			if !isNumerical(temp) {
				return fmt.Errorf("variants for int should be int") //nolint
			}
			fIn = strings.Join(temp, ",")
		}
	default:
		return fmt.Errorf("unsupported type %s", fieldType) //nolint
	}
	return nil
}

func parseStructs(fset *token.FileSet) error {
	structsToValidate = []TemplateStruct{}
	for n, s := range Structs {
		tempFields := []TemplateField{}
		fields := s.Type.(*ast.StructType).Fields.List

		for _, f := range fields {
			var typeName bytes.Buffer
			_ = printer.Fprint(&typeName, fset, f.Type)
			fieldType, err := normalizeType(typeName.String())
			if err != nil {
				return err
			}

			if f.Tag == nil {
				clean()
			} else {
				err = parseTag(f.Tag.Value, fieldType)
				if err != nil {
					return err
				}
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
	return nil
}
