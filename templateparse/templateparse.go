package templateparse

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"github.com/TheLeeeo/grpc-hole/fieldselector"
)

func ParseTemplate(l fieldselector.Selection, input map[string]any, outTemplate map[string]any) (map[string]any, []ParseError) {
	outputTemplate := make(map[string]any)
	var errors []ParseError
	for key, value := range outTemplate {
		fmt.Println("Parsing key", key)
		v, err := ParseField(l.AppendField(key), input, value)
		if err != nil {
			errors = append(errors, err...)
		}
		outputTemplate[key] = v
	}

	return outputTemplate, errors
}

func ParseField(l fieldselector.Selection, input map[string]any, field any) (any, []ParseError) {
	switch f := field.(type) {
	case string: // All strings should be interpreted as a template
		fmt.Println("Parsing string", f)
		v, err := GenerateFieldValue(l, input, f)
		if err != nil {
			return v, []ParseError{err}
		}
		return v, nil
	case map[string]any: // Deal with nested maps
		return ParseTemplate(l, input, f)
	default:
		if reflect.TypeOf(f).Kind() == reflect.Slice || reflect.TypeOf(f).Kind() == reflect.Array {
			fmt.Println("Parsing array")
			return ParseArray(l, input, convertToArray(f))
		}
		return f, nil
	}
}
func convertToArray(input interface{}) []any {
	arr := reflect.ValueOf(input)
	out := make([]any, arr.Len())
	for i := 0; i < arr.Len(); i++ {
		out[i] = arr.Index(i).Interface()
	}
	return out
}

func ParseArray(l fieldselector.Selection, input map[string]any, array []any) ([]any, []ParseError) {
	outSlice := make([]any, len(array))
	var errors []ParseError
	for i, val := range array {
		v, err := ParseField(l.SetIndex(i), input, val)
		if err != nil {
			errors = append(errors, err...)
		}
		outSlice[i] = v
	}

	return outSlice, errors
}

func GenerateFieldValue(l fieldselector.Selection, input map[string]any, fieldTemplate string) (string, ParseError) {
	tmpl, err := template.New("inputParser").Parse(fieldTemplate)
	if err != nil {
		return "<no value>", ParseErrorWrap(l, err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		return "<no value>", ParseErrorWrap(l, err)
	}

	str := buf.String()

	if str == "<no value>" {
		return str, NewParseError(l, "No value found")
	}

	return buf.String(), nil
}
