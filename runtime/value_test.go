package runtime

import (
	"testing"
)

func TestValueTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		wantType ValueType
		wantStr  string
		truthy   bool
	}{
		{"string", NewString("hello"), StringType, "hello", true},
		{"empty string", NewString(""), StringType, "", false},
		{"number", NewNumber(42), NumberType, "42", true},
		{"zero", NewNumber(0), NumberType, "0", false},
		{"float", NewNumber(3.14), NumberType, "3.14", true},
		{"bool true", NewBool(true), BoolType, "true", true},
		{"bool false", NewBool(false), BoolType, "false", false},
		{"null", NewNull(), NullType, "null", false},
		{"array", NewArray(NewNumber(1), NewNumber(2)), ArrayType, "[1, 2]", true},
		{"empty array", NewArray(), ArrayType, "[]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}
			if got := tt.value.String(); got != tt.wantStr {
				t.Errorf("String() = %q, want %q", got, tt.wantStr)
			}
			if got := tt.value.IsTruthy(); got != tt.truthy {
				t.Errorf("IsTruthy() = %v, want %v", got, tt.truthy)
			}
		})
	}
}

func TestObjectValue(t *testing.T) {
	obj := NewObject()
	obj.Set("name", NewString("test"))
	obj.Set("count", NewNumber(42))

	// Test Get
	val, ok := obj.Get("name")
	if !ok {
		t.Fatal("expected to find 'name' field")
	}
	if str, ok := val.(*StringValue); !ok || str.Value != "test" {
		t.Errorf("Get(name) = %v, want 'test'", val)
	}

	// Test missing field
	_, ok = obj.Get("missing")
	if ok {
		t.Error("expected 'missing' field to not exist")
	}

	// Test truthy
	if !obj.IsTruthy() {
		t.Error("non-empty object should be truthy")
	}

	emptyObj := NewObject()
	if emptyObj.IsTruthy() {
		t.Error("empty object should be falsy")
	}
}

