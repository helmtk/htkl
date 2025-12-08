package funcs

import (
	"testing"

	"github.com/helmtk/htkl/runtime"
)

func TestUpperFunc(t *testing.T) {
	result, err := upperFunc(runtime.NewString("hello"))
	if err != nil {
		t.Fatalf("upper() error = %v", err)
	}
	if result.String() != "HELLO" {
		t.Errorf("upper(hello) = %v, want HELLO", result.String())
	}
}

func TestLowerFunc(t *testing.T) {
	result, err := lowerFunc(runtime.NewString("WORLD"))
	if err != nil {
		t.Fatalf("lower() error = %v", err)
	}
	if result.String() != "world" {
		t.Errorf("lower(WORLD) = %v, want world", result.String())
	}
}

func TestTrimFunc(t *testing.T) {
	result, err := trimFunc(runtime.NewString("  hello  "))
	if err != nil {
		t.Fatalf("trim() error = %v", err)
	}
	if result.String() != "hello" {
		t.Errorf("trim('  hello  ') = %v, want 'hello'", result.String())
	}
}

func TestQuoteFunc(t *testing.T) {
	result, err := quoteFunc(runtime.NewString("hello"))
	if err != nil {
		t.Fatalf("quote() error = %v", err)
	}
	if result.String() != `"hello"` {
		t.Errorf("quote(hello) = %v, want \"hello\"", result.String())
	}
}

func TestNindentFunc(t *testing.T) {
	result, err := nindentFunc(runtime.NewString("line1\nline2"), runtime.NewNumber(2))
	if err != nil {
		t.Fatalf("nindent() error = %v", err)
	}
	want := "  line1\n  line2"
	if result.String() != want {
		t.Errorf("nindent() = %q, want %q", result.String(), want)
	}
}

func TestDefaultFunc(t *testing.T) {
	tests := []struct {
		name string
		args []runtime.Value
		want string
	}{
		{
			name: "null returns default",
			args: []runtime.Value{runtime.NewString("default"), runtime.NewNull()},
			want: "default",
		},
		{
			name: "empty string returns default",
			args: []runtime.Value{runtime.NewString("default"), runtime.NewString("")},
			want: "default",
		},
		{
			name: "value returns value",
			args: []runtime.Value{runtime.NewString("default"), runtime.NewString("value")},
			want: "value",
		},
		{
			name: "zero returns default",
			args: []runtime.Value{runtime.NewString("default"), runtime.NewNumber(0)},
			want: "default",
		},
		{
			name: "non-zero returns value",
			args: []runtime.Value{runtime.NewString("default"), runtime.NewNumber(42)},
			want: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := defaultFunc(tt.args...)
			if err != nil {
				t.Fatalf("default() error = %v", err)
			}
			if result.String() != tt.want {
				t.Errorf("default() = %v, want %v", result.String(), tt.want)
			}
		})
	}
}

func TestLenFunc(t *testing.T) {
	tests := []struct {
		name string
		arg  runtime.Value
		want float64
	}{
		{
			name: "string length",
			arg:  runtime.NewString("hello"),
			want: 5,
		},
		{
			name: "array length",
			arg:  runtime.NewArray(runtime.NewNumber(1), runtime.NewNumber(2), runtime.NewNumber(3)),
			want: 3,
		},
		{
			name: "object length",
			arg: func() runtime.Value {
				obj := runtime.NewObject()
				obj.Set("a", runtime.NewNumber(1))
				obj.Set("b", runtime.NewNumber(2))
				return obj
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lenFunc(tt.arg)
			if err != nil {
				t.Fatalf("len() error = %v", err)
			}
			num, ok := result.(*runtime.NumberValue)
			if !ok {
				t.Fatalf("expected NumberValue, got %T", result)
			}
			if num.Value != tt.want {
				t.Errorf("len() = %v, want %v", num.Value, tt.want)
			}
		})
	}
}

func TestMathFuncs(t *testing.T) {
	tests := []struct {
		name string
		fn   runtime.Func
		arg  float64
		want float64
	}{
		{"round up", roundFunc, 3.6, 4},
		{"round down", roundFunc, 3.4, 3},
		{"floor", floorFunc, 3.9, 3},
		{"ceil", ceilFunc, 3.1, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.fn(runtime.NewNumber(tt.arg))
			if err != nil {
				t.Fatalf("%s() error = %v", tt.name, err)
			}
			num, ok := result.(*runtime.NumberValue)
			if !ok {
				t.Fatalf("expected NumberValue, got %T", result)
			}
			if num.Value != tt.want {
				t.Errorf("%s() = %v, want %v", tt.name, num.Value, tt.want)
			}
		})
	}
}
