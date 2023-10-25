package templateparse

import (
	"bytes"
	"text/template"
)

func ParseTemplate(input map[string]any, outTemplate map[string]any) (map[string]any, []ParseError) {
	outputTemplate := make(map[string]any)
	var errors []ParseError
	for key, value := range outTemplate {
		v, err := ParseField(input, value)
		if err != nil {
			errors = append(errors, err...)
		}
		outputTemplate[key] = v
	}

	return outputTemplate, errors
}

func ParseField(input map[string]any, field any) (any, []ParseError) {
	switch f := field.(type) {
	case string: // All strings should be interpreted as a template
		v, err := GenerateFieldValue(input, f)
		if err != nil {
			return v, []ParseError{err}
		}
		return v, nil
	case map[string]any: // Deal with nested maps
		return ParseTemplate(input, f)
	case []any:
		return ParseArray(input, f)
	default:
		return f, nil
	}
}

func ParseArray(input map[string]any, array []any) ([]any, []ParseError) {
	outSlice := make([]any, len(array))
	var errors []ParseError
	for i, val := range array {
		v, err := ParseField(input, val)
		if err != nil {
			errors = append(errors, err...)
		}
		outSlice[i] = v
	}

	return outSlice, errors
}

func GenerateFieldValue(input map[string]any, fieldTemplate string) (string, ParseError) {
	tmpl, err := template.New("inputParser").Parse(fieldTemplate)
	if err != nil {
		return "<no value>", ParseErrorWrap("", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		return "<no value>", ParseErrorWrap("", err)
	}

	str := buf.String()

	if str == "<no value>" {
		return str, NewParseError("", "No value found")
	}

	return buf.String(), nil
}
