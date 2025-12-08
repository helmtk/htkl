package eval

import (
	"strings"
	"testing"

	"helmtk.dev/code/htkl/parser"
	"helmtk.dev/code/htkl/runtime"
)

func TestArithmetic(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "addition",
			input: "result: 5 + 3",
			want:  map[string]string{"result": "8"},
		},
		{
			name:  "subtraction",
			input: "result: 10 - 3",
			want:  map[string]string{"result": "7"},
		},
		{
			name:  "multiplication",
			input: "result: 4 * 5",
			want:  map[string]string{"result": "20"},
		},
		{
			name:  "division",
			input: "result: 15 / 3",
			want:  map[string]string{"result": "5"},
		},
		{
			name:  "negation",
			input: "result: -42",
			want:  map[string]string{"result": "-42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := evalToObject(t, tt.input)
			for key, want := range tt.want {
				got := getString(t, obj, key)
				if got != want {
					t.Errorf("%s: got %q, want %q", key, got, want)
				}
			}
		})
	}
}

func TestStringConcat(t *testing.T) {
	obj := evalToObject(t, `result: "hello" + " world"`)
	got := getString(t, obj, "result")
	want := "hello world"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestComparison(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"equal", "5 == 5", true},
		{"not equal", "5 != 3", true},
		{"less than", "3 < 5", true},
		{"less than or equal", "5 <= 5", true},
		{"greater than", "10 > 5", true},
		{"greater than or equal", "5 >= 5", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := evalToObject(t, "result: "+tt.input)
			got := getBool(t, obj, "result")
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogical(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"and true", "true && true", true},
		{"and false", "true && false", false},
		{"or true", "false || true", true},
		{"or false", "false || false", false},
		{"not true", "!false", true},
		{"not false", "!true", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := evalToObject(t, "result: "+tt.input)
			got := getBool(t, obj, "result")
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLiterals(t *testing.T) {
	obj := evalToObject(t, `
string: "hello"
number: 42
bool: true
	`)

	if got := getString(t, obj, "string"); got != "hello" {
		t.Errorf("string: got %q, want %q", got, "hello")
	}
	if got := getString(t, obj, "number"); got != "42" {
		t.Errorf("number: got %q, want %q", got, "42")
	}
	if got := getBool(t, obj, "bool"); !got {
		t.Errorf("bool: got false, want true")
	}
}

func TestArrays(t *testing.T) {
	obj := evalToObject(t, `
let items = [1, 2, 3]
first: items[0]
second: items[1]
	`)

	if got := getString(t, obj, "first"); got != "1" {
		t.Errorf("first: got %q, want %q", got, "1")
	}
	if got := getString(t, obj, "second"); got != "2" {
		t.Errorf("second: got %q, want %q", got, "2")
	}
}

func TestObjects(t *testing.T) {
	obj := evalToObject(t, `
let person = {
	name: "Alice"
	age: 30
}
personName: person.name
personAge: person.age
	`)

	if got := getString(t, obj, "personName"); got != "Alice" {
		t.Errorf("personName: got %q, want %q", got, "Alice")
	}
	if got := getString(t, obj, "personAge"); got != "30" {
		t.Errorf("personAge: got %q, want %q", got, "30")
	}
}

func TestNestedObjects(t *testing.T) {
	obj := evalToObject(t, `
let config = {
	server: {
		host: "localhost"
		port: 8080
	}
}
host: config.server.host
	`)

	if got := getString(t, obj, "host"); got != "localhost" {
		t.Errorf("host: got %q, want %q", got, "localhost")
	}
}

func TestLetStatement(t *testing.T) {
	obj := evalToObject(t, `
let x = 10
let y = 20
result: x + y
	`)

	if got := getString(t, obj, "result"); got != "30" {
		t.Errorf("result: got %q, want %q", got, "30")
	}
}

func TestAssignment(t *testing.T) {
	obj := evalToObject(t, `
let x = 10
x = 20
result: x
	`)

	if got := getString(t, obj, "result"); got != "20" {
		t.Errorf("result: got %q, want %q", got, "20")
	}
}

func TestForStatement(t *testing.T) {
	obj := evalToObject(t, `
let items = [1, 2, 3]
results: [for i, item in items do item * 2 end]
	`)

	arr := getArray(t, obj, "results")
	want := []string{"2", "4", "6"}
	if len(arr.Elements) != len(want) {
		t.Fatalf("expected %d elements, got %d", len(want), len(arr.Elements))
	}
	for i, w := range want {
		got := arr.Elements[i].String()
		if got != w {
			t.Errorf("results[%d]: got %q, want %q", i, got, w)
		}
	}
}

func TestWithStatement(t *testing.T) {
	obj := evalToObject(t, `
let config = {name: "test"}
result: with config as cfg do cfg.name end
	`)

	if got := getString(t, obj, "result"); got != "test" {
		t.Errorf("result: got %q, want %q", got, "test")
	}
}

func TestSpread(t *testing.T) {
	t.Run("array spread", func(t *testing.T) {
		obj := evalToObject(t, `
let a = [1, 2]
let b = [3, 4]
result: [spread a, spread b]
		`)

		arr := getArray(t, obj, "result")
		want := []string{"1", "2", "3", "4"}
		if len(arr.Elements) != len(want) {
			t.Fatalf("expected %d elements, got %d", len(want), len(arr.Elements))
		}
		for i, w := range want {
			got := arr.Elements[i].String()
			if got != w {
				t.Errorf("result[%d]: got %q, want %q", i, got, w)
			}
		}
	})

	t.Run("object spread", func(t *testing.T) {
		obj := evalToObject(t, `
let a = {x: 1}
let b = {y: 2}
result: {spread a, spread b}
		`)

		if got := getString(t, obj, "result.x"); got != "1" {
			t.Errorf("result.x: got %q, want %q", got, "1")
		}
		if got := getString(t, obj, "result.y"); got != "2" {
			t.Errorf("result.y: got %q, want %q", got, "2")
		}
	})
}

func TestTemplates(t *testing.T) {
	result := eval(t, `
define("makeLabel") do
	app: "myapp"
end

labels: {
	include("makeLabel")
}
	`)

	obj := getDocument(t, result, 0)
	if got := getString(t, obj, "labels.app"); got != "myapp" {
		t.Errorf("labels.app: got %q, want %q", got, "myapp")
	}
}

func TestTemplateContext(t *testing.T) {
	scope := runtime.NewScope(nil)
	values := runtime.NewObject()
	values.Set("app", runtime.NewString("foo"))
	scope.Set("Values", values)

	result := evalWithScope(t, scope, `
define("makeLabel") do
	app: Values.app
	always: "always"
end

labels: {
	include("makeLabel")
}
	`)

	obj := getDocument(t, result, 0)
	if got := getString(t, obj, "labels.app"); got != "foo" {
		t.Errorf("labels.app: got %q, want %q", got, "foo")
	}
	if got := getString(t, obj, "labels.always"); got != "always" {
		t.Errorf("labels.always: got %q, want %q", got, "always")
	}
}

func TestPipes(t *testing.T) {
	scope := runtime.NewScope(nil)
	scope.SetFunction("upper", func(args ...runtime.Value) (runtime.Value, error) {
		s, err := runtime.ToString(args[0])
		if err != nil {
			return nil, err
		}
		return runtime.NewString(strings.ToUpper(s)), nil
	})

	result := evalWithScope(t, scope, `result: "hello" | upper`)

	obj := getDocument(t, result, 0)
	if got := getString(t, obj, "result"); got != "HELLO" {
		t.Errorf("result: got %q, want %q", got, "HELLO")
	}
}

func TestMultipleDocuments(t *testing.T) {
	result := eval(t, `
{kind: "ConfigMap"}
{kind: "Deployment"}
	`)

	arr, ok := result.(*runtime.ArrayValue)
	if !ok {
		t.Fatalf("expected ArrayValue, got %T", result)
	}

	if len(arr.Elements) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(arr.Elements))
	}

	doc1 := arr.Elements[0].(*runtime.ObjectValue)
	if got := getString(t, doc1, "kind"); got != "ConfigMap" {
		t.Errorf("doc1.kind: got %q, want %q", got, "ConfigMap")
	}

	doc2 := arr.Elements[1].(*runtime.ObjectValue)
	if got := getString(t, doc2, "kind"); got != "Deployment" {
		t.Errorf("doc2.kind: got %q, want %q", got, "Deployment")
	}
}

// Error tests

func TestErrorDivisionByZero(t *testing.T) {
	expectError(t, "result: 10 / 0", "division by zero")
}

func TestErrorArrayIndexBounds(t *testing.T) {
	expectError(t, `
let items = [1, 2, 3]
result: items[10]
	`, "array index out of bounds")
}

func TestErrorUndefinedFunction(t *testing.T) {
	expectError(t, `result: unknownFunc()`, "undefined function")
}

func TestErrorUndefinedTemplate(t *testing.T) {
	expectError(t, `include("unknown")`, "undefined template")
}

// Helper functions

func eval(t *testing.T, input string) runtime.Value {
	t.Helper()
	return evalWithScope(t, runtime.NewScope(nil), input)
}

func evalWithScope(t *testing.T, scope *runtime.Scope, input string) runtime.Value {
	t.Helper()
	doc, err := parser.New(input, "test.helmtk").Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := EvalDocument(doc, scope)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}

	return result
}

func evalToObject(t *testing.T, input string) *runtime.ObjectValue {
	t.Helper()
	result := eval(t, input)
	return getDocument(t, result, 0)
}

func getDocument(t *testing.T, result runtime.Value, index int) *runtime.ObjectValue {
	t.Helper()
	arr, ok := result.(*runtime.ArrayValue)
	if !ok {
		t.Fatalf("expected ArrayValue, got %T", result)
	}
	if index >= len(arr.Elements) {
		t.Fatalf("document index %d out of bounds (len=%d)", index, len(arr.Elements))
	}
	obj, ok := arr.Elements[index].(*runtime.ObjectValue)
	if !ok {
		t.Fatalf("expected ObjectValue at index %d, got %T", index, arr.Elements[index])
	}
	return obj
}

func getString(t *testing.T, obj *runtime.ObjectValue, path string) string {
	t.Helper()
	val := getPath(t, obj, path)
	return val.String()
}

func getBool(t *testing.T, obj *runtime.ObjectValue, path string) bool {
	t.Helper()
	val := getPath(t, obj, path)
	return val.IsTruthy()
}

func getArray(t *testing.T, obj *runtime.ObjectValue, path string) *runtime.ArrayValue {
	t.Helper()
	val := getPath(t, obj, path)
	arr, ok := val.(*runtime.ArrayValue)
	if !ok {
		t.Fatalf("expected ArrayValue at %s, got %T", path, val)
	}
	return arr
}

func getPath(t *testing.T, obj *runtime.ObjectValue, path string) runtime.Value {
	t.Helper()
	parts := strings.Split(path, ".")
	var val runtime.Value = obj

	for _, part := range parts {
		obj, ok := val.(*runtime.ObjectValue)
		if !ok {
			t.Fatalf("expected ObjectValue at %s, got %T", part, val)
		}
		var found bool
		val, found = obj.Get(part)
		if !found {
			t.Fatalf("field %s not found", part)
		}
	}

	return val
}

func expectError(t *testing.T, input string, wantErr string) {
	t.Helper()
	doc, err := parser.New(input, "test.helmtk").Parse()
	if err != nil {
		if strings.Contains(err.Error(), wantErr) {
			return
		}
		t.Fatalf("unexpected parse error: %v", err)
	}

	scope := runtime.NewScope(nil)
	_, err = EvalDocument(doc, scope)
	if err == nil {
		t.Fatal("expected error but got none")
	}

	if !strings.Contains(err.Error(), wantErr) {
		t.Errorf("error mismatch\ngot: %v\nwant substring: %s", err, wantErr)
	}
}
