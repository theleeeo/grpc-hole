package templateparse

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/TheLeeeo/grpc-hole/fieldselector"
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
	fmt.Println(actual, err)
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
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
		"baz": "1.1",
	}
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
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
		"baz": "true",
	}
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
	if err == nil {
		t.Error(err)
	}
	if len(err) != 1 {
		t.Errorf("Expected 1 error, got %v", err)
	}
	if !strings.Contains(err[0].Error(), "unclosed action") {
		t.Errorf("Expected \"unclosed action\", got %v", err[0])
	}
	if err[0].Location() != ".baz" {
		t.Errorf("Expected \".baz\", got %v", err[0].Location())
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
	if err == nil {
		t.Error(err)
	}
	if len(err) != 1 {
		t.Errorf("Expected 1 error, got %v", err)
	}
	if err[0].Error() != "No value found" {
		t.Errorf("Expected \"No value found\", got %v", err)
	}
	if err[0].Location() != ".baz" {
		t.Errorf("Expected \".baz\", got %v", err[0].Location())
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
	if err == nil {
		t.Error(err)
	}
	if len(err) != 1 {
		t.Errorf("Expected 1 error, got %v", err)
	}
	if err[0].Error() != "No value found" {
		t.Errorf("Expected \"No value found\", got %v", err)
	}
	if err[0].Location() != ".fizz" {
		t.Errorf("Expected \".fizz\", got %v", err[0].Location())
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
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

// This test fails sometimes because the order of errors is not guaranteed.
func Test_ErrorsInNestedPlace(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"baz": true,
	}
	outTemplate := map[string]any{
		"bar": []any{"{{.bad}}", "{{.bad"},
		"fizz": map[string]any{
			"buzz": "{{.bad}}",
			"test": []any{"{{.baz}}", "{{.bad"},
		},
	}
	expected := map[string]any{
		"bar": []any{"<no value>", "<no value>"},
		"fizz": map[string]any{
			"buzz": "<no value>",
			"test": []any{"true", "<no value>"},
		},
	}
	actual, err := ParseTemplate(fieldselector.Root, input, outTemplate)
	if err == nil {
		t.Error(err)
	}
	if len(err) != 4 {
		t.Errorf("Expected 4 errors, got %v", err)
	}
	if !strings.Contains(err[0].Error(), "No value found") {
		t.Errorf("Expected \"No value found\", got %v", err)
	}
	if err[0].Location() != ".bar[0]" {
		t.Errorf("Expected \".bar[0]\", got %v", err[0].Location())
	}
	if !strings.Contains(err[1].Error(), "unclosed action") {
		t.Errorf("Expected \"unclosed action\", got %v", err[1])
	}
	if err[1].Location() != ".bar[1]" {
		t.Errorf("Expected \".bar[1]\", got %v", err[1].Location())
	}
	if !strings.Contains(err[2].Error(), "No value found") {
		t.Errorf("Expected \"No value found\", got %v", err[2])
	}
	if err[2].Location() != ".fizz.buzz" {
		t.Errorf("Expected \".fizz.buzz\", got %v", err[2].Location())
	}
	if !strings.Contains(err[3].Error(), "unclosed action") {
		t.Errorf("Expected \"unclosed action\", got %v", err[3])
	}
	if err[3].Location() != ".fizz.test[1]" {
		t.Errorf("Expected \".fizz.test[1]\", got %v", err[3].Location())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
