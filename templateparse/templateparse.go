package templateparse

import (
	"bytes"
	"fmt"
	"text/template"
)

func ParseTemplate(input map[string]any, outTemplate map[string]any) (map[string]any, error) {
	outputTemplate := make(map[string]any)
	for key, value := range outTemplate {
		v, err := ParseField(input, value)
		if err != nil {
			return nil, err
		}
		outputTemplate[key] = v
	}

	return outputTemplate, nil
}

func ParseField(input map[string]any, field any) (any, error) {
	var out any
	var err error
	switch val := field.(type) {
	case string: // All strings should be interpreted as a template
		out, err = GenerateFieldValue(input, val)
		if err != nil {
			return nil, err
		}
	case map[string]any: // Deal with nested maps
		out, err = ParseTemplate(input, val)
		if err != nil {
			return nil, err
		}
	case []any:
		outSlice := make([]any, len(val))
		for i, v := range val {
			outSlice[i], err = ParseField(input, v)
			if err != nil {
				return nil, err
			}
		}
		out = outSlice
	default:
		out = val
	}

	return out, nil
}

// TODO: Return a interface{} instead of a string
func GenerateFieldValue(input map[string]any, fieldTemplate string) (string, error) {
	tmpl, err := template.New("inputParser").Parse(fieldTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return "", err
	}

	stringReturn := buf.String()

	return stringReturn, nil
}
