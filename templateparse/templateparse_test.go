package templateparse

import (
	"reflect"
	"testing"
)

func Test_ParseTemplate(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": "{{.baz}}",
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": "1", // TODO: Should this be an int?
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_NonStringsInTemplate(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"foo": true,
		"baz": 2,
		"fizz": map[string]any{
			"buzz": 1.2,
		},
	}
	expected := map[string]any{
		"foo": true,
		"baz": 2,
		"fizz": map[string]any{
			"buzz": 1.2,
		},
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_Nested(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": map[string]any{
			"baz": "{{.baz}}",
		},
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"baz": "1",
		},
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_NestedArray(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": []any{
			"{{.baz}}",
		},
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": []any{
			"1",
		},
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_Float(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1.1,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": "{{.baz}}",
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": "1.1", // TODO: Should this be a float?
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_Bool(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": true,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": "{{.baz}}",
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": "true", // TODO: Should this be a bool?
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_NestedArrayWithMap(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": []any{
			map[string]any{
				"baz": "{{.baz}}",
			},
		},
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": []any{
			map[string]any{
				"baz": "1",
			},
		},
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_TemplateLogic(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": []any{
			map[string]any{
				"baz": "{{if eq .baz 1}}{{.baz}}{{end}}",
			},
		},
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": []any{
			map[string]any{
				"baz": "1",
			},
		},
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_TemplateLogic_Empty(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 2,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": []any{
			map[string]any{
				"baz": "{{if eq .baz 1}}{{.baz}}{{end}}",
			},
		},
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": []any{
			map[string]any{
				"baz": "",
			},
		},
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_InvalidTemplate(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": "{{.baz",
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": "<no value>",
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err == nil {
		t.Errorf("Expected error, got %v", actual)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_InvalidSelector(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"foo": "{{.foo}}",
		"baz": "{{.bar}}",
	}
	expected := map[string]any{
		"foo": "bar",
		"baz": "<no value>",
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_TemplateNotMatch(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": 1,
	}
	outTemplate := map[string]any{
		"bar":  "{{.foo}}",
		"fizz": "{{.buzz}}",
	}
	expected := map[string]any{
		"bar":  "bar",
		"fizz": "<no value>",
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func Test_ParseTemplate_TemplateDifferButCurrect(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": true,
	}
	outTemplate := map[string]any{
		"bar":  "{{.foo}}",
		"fizz": "{{.baz}}",
	}
	expected := map[string]any{
		"bar":  "bar",
		"fizz": "true",
	}
	actual, err := ParseTemplate(input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
