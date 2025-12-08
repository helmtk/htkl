package parser

import (
	"testing"
)

func TestParseEmpty(t *testing.T) {
	input := ""
	doc, err := New(input, "").Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if doc == nil {
		t.Fatal("expected document, got nil")
	}

	if len(doc.Body) != 0 {
		t.Errorf("expected 0 Body, got %d", len(doc.Body))
	}
}

func TestParseSimpleKeyValue(t *testing.T) {
	input := `apiVersion: "apps/v1"`

	doc, err := New(input, "").Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if doc == nil {
		t.Fatal("expected document, got nil")
	}

	if len(doc.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(doc.Body))
	}

	kv, ok := doc.Body[0].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue, got %T", doc.Body[0])
	}

	if kv.Key != "apiVersion" {
		t.Errorf("expected key 'apiVersion', got '%s'", kv.Key)
	}

	str, ok := kv.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral value, got %T", kv.Value)
	}

	if str.Value != "apps/v1" {
		t.Errorf("expected value 'apps/v1', got '%s'", str.Value)
	}
}

func TestParseObject(t *testing.T) {
	input := `
metadata: {
	name: "example"
	labels: {
		app: "test"
	}
}
`

	doc, err := New(input, "").Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if doc == nil {
		t.Fatal("expected document, got nil")
	}

	if len(doc.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(doc.Body))
	}

	kv, ok := doc.Body[0].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue, got %T", doc.Body[0])
	}

	if kv.Key != "metadata" {
		t.Errorf("expected key 'metadata', got '%s'", kv.Key)
	}

	obj, ok := kv.Value.(*Object)
	if !ok {
		t.Fatalf("expected Object value, got %T", kv.Value)
	}

	if len(obj.Body) != 2 {
		t.Fatalf("expected 2 Body, got %d", len(obj.Body))
	}

	// Check first field: name: "example"
	field0, ok := obj.Body[0].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue for field 0, got %T", obj.Body[0])
	}
	if field0.Key != "name" {
		t.Errorf("expected first field key 'name', got '%s'", field0.Key)
	}

	nameStr, ok := field0.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral for name, got %T", field0.Value)
	}

	if nameStr.Value != "example" {
		t.Errorf("expected name value 'example', got '%s'", nameStr.Value)
	}

	// Check second field: labels object
	field1, ok := obj.Body[1].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue for field 1, got %T", obj.Body[1])
	}
	if field1.Key != "labels" {
		t.Errorf("expected second field key 'labels', got '%s'", field1.Key)
	}

	labelsObj, ok := field1.Value.(*Object)
	if !ok {
		t.Fatalf("expected Object for labels, got %T", field1.Value)
	}

	if len(labelsObj.Body) != 1 {
		t.Fatalf("expected 1 field in labels, got %d", len(labelsObj.Body))
	}

	labelsField0, ok := labelsObj.Body[0].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue for labels field 0, got %T", labelsObj.Body[0])
	}
	if labelsField0.Key != "app" {
		t.Errorf("expected labels field key 'app', got '%s'", labelsField0.Key)
	}

	appStr, ok := labelsField0.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral for app, got %T", labelsField0.Value)
	}

	if appStr.Value != "test" {
		t.Errorf("expected app value 'test', got '%s'", appStr.Value)
	}
}

func TestParseIdentifier(t *testing.T) {
	input := `name: Values.name`

	doc, err := New(input, "").Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if doc == nil {
		t.Fatal("expected document, got nil")
	}

	if len(doc.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(doc.Body))
	}

	kv, ok := doc.Body[0].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue, got %T", doc.Body[0])
	}

	if kv.Key != "name" {
		t.Errorf("expected key 'name', got '%s'", kv.Key)
	}

	// Values.name should be parsed as MemberExpression
	member, ok := kv.Value.(*MemberExpression)
	if !ok {
		t.Fatalf("expected MemberExpression value, got %T", kv.Value)
	}

	// Check object is "Values"
	obj, ok := member.Object.(*Identifier)
	if !ok {
		t.Fatalf("expected Identifier object, got %T", member.Object)
	}

	if obj.Name != "Values" {
		t.Errorf("expected object 'Values', got '%s'", obj.Name)
	}

	// Check member is "name"
	if member.Member != "name" {
		t.Errorf("expected member 'name', got '%s'", member.Member)
	}
}

func TestParseArray(t *testing.T) {
	input := `
ports: [
	{name: "http", containerPort: 80}
	{name: "debug", containerPort: 5005}
]
`

	doc, err := New(input, "").Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if doc == nil {
		t.Fatal("expected document, got nil")
	}

	if len(doc.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(doc.Body))
	}

	kv, ok := doc.Body[0].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue, got %T", doc.Body[0])
	}

	if kv.Key != "ports" {
		t.Errorf("expected key 'ports', got '%s'", kv.Key)
	}

	arr, ok := kv.Value.(*Array)
	if !ok {
		t.Fatalf("expected Array value, got %T", kv.Value)
	}

	if len(arr.Body) != 2 {
		t.Fatalf("expected 2 Body, got %d", len(arr.Body))
	}

	// Check first element
	obj1, ok := arr.Body[0].(*Object)
	if !ok {
		t.Fatalf("expected Object for first element, got %T", arr.Body[0])
	}

	if len(obj1.Body) != 2 {
		t.Fatalf("expected 2 Body in first object, got %d", len(obj1.Body))
	}

	obj1field0, ok := obj1.Body[0].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue for obj1 field 0, got %T", obj1.Body[0])
	}
	if obj1field0.Key != "name" {
		t.Errorf("expected first field key 'name', got '%s'", obj1field0.Key)
	}

	nameStr, ok := obj1field0.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral for name, got %T", obj1field0.Value)
	}

	if nameStr.Value != "http" {
		t.Errorf("expected name 'http', got '%s'", nameStr.Value)
	}

	obj1field1, ok := obj1.Body[1].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue for obj1 field 1, got %T", obj1.Body[1])
	}
	if obj1field1.Key != "containerPort" {
		t.Errorf("expected second field key 'containerPort', got '%s'", obj1field1.Key)
	}

	port, ok := obj1field1.Value.(*NumberLiteral)
	if !ok {
		t.Fatalf("expected NumberLiteral for containerPort, got %T", obj1field1.Value)
	}

	if port.Value != 80 {
		t.Errorf("expected port 80, got %f", port.Value)
	}

	// Check second element
	obj2, ok := arr.Body[1].(*Object)
	if !ok {
		t.Fatalf("expected Object for second element, got %T", arr.Body[1])
	}

	if len(obj2.Body) != 2 {
		t.Fatalf("expected 2 Body in second object, got %d", len(obj2.Body))
	}

	obj2field0, ok := obj2.Body[0].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue for obj2 field 0, got %T", obj2.Body[0])
	}
	nameStr2, ok := obj2field0.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral for name, got %T", obj2field0.Value)
	}

	if nameStr2.Value != "debug" {
		t.Errorf("expected name 'debug', got '%s'", nameStr2.Value)
	}

	obj2field1, ok := obj2.Body[1].(*KeyValueStatement)
	if !ok {
		t.Fatalf("expected KeyValue for obj2 field 1, got %T", obj2.Body[1])
	}
	port2, ok := obj2field1.Value.(*NumberLiteral)
	if !ok {
		t.Fatalf("expected NumberLiteral for containerPort, got %T", obj2field1.Value)
	}

	if port2.Value != 5005 {
		t.Errorf("expected port 5005, got %f", port2.Value)
	}
}

func TestParseStringEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escaped double quote",
			input:    `message: "Hello \"World\""`,
			expected: `Hello "World"`,
		},
		{
			name:     "escaped backslash",
			input:    `path: "C:\\Users\\test"`,
			expected: `C:\Users\test`,
		},
		{
			name:     "escaped newline",
			input:    `text: "line1\nline2"`,
			expected: "line1\nline2",
		},
		{
			name:     "escaped tab",
			input:    `text: "col1\tcol2"`,
			expected: "col1\tcol2",
		},
		{
			name:     "escaped carriage return",
			input:    `text: "line1\rline2"`,
			expected: "line1\rline2",
		},
		{
			name:     "multiple escapes",
			input:    `text: "He said \"hello\\world\"\nNext line"`,
			expected: "He said \"hello\\world\"\nNext line",
		},
		{
			name:     "no escapes",
			input:    `text: "simple string"`,
			expected: "simple string",
		},
		{
			name:     "escaped dollar prevents interpolation",
			input:    `text: "Price is \${100}"`,
			expected: `Price is ${100}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := New(tt.input, "").Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(doc.Body) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(doc.Body))
			}

			kv, ok := doc.Body[0].(*KeyValueStatement)
			if !ok {
				t.Fatalf("expected KeyValue, got %T", doc.Body[0])
			}

			str, ok := kv.Value.(*StringLiteral)
			if !ok {
				t.Fatalf("expected StringLiteral value, got %T", kv.Value)
			}

			if str.Value != tt.expected {
				t.Errorf("expected value %q, got %q", tt.expected, str.Value)
			}
		})
	}
}

func TestParseInterpolatedStringEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string // expected string parts (non-expressions, empty strings omitted)
	}{
		{
			name:     "escaped quote before interpolation",
			input:    `text: "He said \"hello\" ${name}"`,
			expected: []string{`He said "hello" `},
		},
		{
			name:     "escaped quote after interpolation",
			input:    `text: "${name} said \"hello\""`,
			expected: []string{` said "hello"`},
		},
		{
			name:     "escaped backslash in interpolated string",
			input:    `path: "${base}\\subdir"`,
			expected: []string{`\subdir`},
		},
		{
			name:     "escaped newline with interpolation",
			input:    `text: "line1\n${value}"`,
			expected: []string{"line1\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := New(tt.input, "").Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(doc.Body) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(doc.Body))
			}

			kv, ok := doc.Body[0].(*KeyValueStatement)
			if !ok {
				t.Fatalf("expected KeyValue, got %T", doc.Body[0])
			}

			interp, ok := kv.Value.(*InterpolatedString)
			if !ok {
				t.Fatalf("expected InterpolatedString value, got %T", kv.Value)
			}

			// Collect string literal parts
			var parts []string
			for _, part := range interp.Parts {
				if strLit, ok := part.(*StringLiteral); ok {
					parts = append(parts, strLit.Value)
				}
			}

			if len(parts) != len(tt.expected) {
				t.Fatalf("expected %d string parts, got %d", len(tt.expected), len(parts))
			}

			for i, expected := range tt.expected {
				if parts[i] != expected {
					t.Errorf("part %d: expected %q, got %q", i, expected, parts[i])
				}
			}
		})
	}
}
