package runtime

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ValueType represents the type of a runtime value
type ValueType int

const (
	StringType ValueType = iota
	NumberType
	BoolType
	NullType
	ArrayType
	ObjectType
)

func (vt ValueType) String() string {
	switch vt {
	case StringType:
		return "string"
	case NumberType:
		return "number"
	case BoolType:
		return "bool"
	case NullType:
		return "null"
	case ArrayType:
		return "array"
	case ObjectType:
		return "object"
	default:
		return "unknown"
	}
}

// Value is the interface for all runtime values
type Value interface {
	Type() ValueType
	String() string
	IsTruthy() bool
}

// StringValue represents a string value
type StringValue struct {
	Value string
}

func (s *StringValue) Type() ValueType { return StringType }
func (s *StringValue) String() string  { return s.Value }
func (s *StringValue) IsTruthy() bool  { return s.Value != "" }

// NumberValue represents a numeric value
type NumberValue struct {
	Value float64
}

func (n *NumberValue) Type() ValueType { return NumberType }
func (n *NumberValue) String() string  { return strconv.FormatFloat(n.Value, 'f', -1, 64) }
func (n *NumberValue) IsTruthy() bool  { return n.Value != 0 }

// BoolValue represents a boolean value
type BoolValue struct {
	Value bool
}

func (b *BoolValue) Type() ValueType { return BoolType }
func (b *BoolValue) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b *BoolValue) IsTruthy() bool { return b.Value }

// NullValue represents a null value
type NullValue struct{}

func (n *NullValue) Type() ValueType { return NullType }
func (n *NullValue) String() string  { return "null" }
func (n *NullValue) IsTruthy() bool  { return false }

// ArrayValue represents an array of values
type ArrayValue struct {
	Elements []Value
}

func (a *ArrayValue) Type() ValueType { return ArrayType }
func (a *ArrayValue) String() string {
	var parts []string
	for _, elem := range a.Elements {
		parts = append(parts, elem.String())
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
func (a *ArrayValue) IsTruthy() bool { return len(a.Elements) > 0 }

// ObjectValue represents an object (map of string keys to values)
type ObjectValue struct {
	Fields map[string]Value
}

func (o *ObjectValue) Type() ValueType { return ObjectType }
func (o *ObjectValue) String() string {
	var parts []string
	for k, v := range o.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
func (o *ObjectValue) IsTruthy() bool { return len(o.Fields) > 0 }

// Get retrieves a field from the object
func (o *ObjectValue) Get(key string) (Value, bool) {
	val, ok := o.Fields[key]
	return val, ok
}

// Set sets a field in the object
func (o *ObjectValue) Set(key string, val Value) {
	if o.Fields == nil {
		o.Fields = make(map[string]Value)
	}
	o.Fields[key] = val
}

// Helper functions for type checking

func IsString(v Value) bool {
	_, ok := v.(*StringValue)
	return ok
}

func IsNumber(v Value) bool {
	_, ok := v.(*NumberValue)
	return ok
}

func IsBool(v Value) bool {
	_, ok := v.(*BoolValue)
	return ok
}

func IsNull(v Value) bool {
	_, ok := v.(*NullValue)
	return ok
}

func IsArray(v Value) bool {
	_, ok := v.(*ArrayValue)
	return ok
}

func IsObject(v Value) bool {
	_, ok := v.(*ObjectValue)
	return ok
}

// Helper functions for type conversion

func ToString(v Value) (string, error) {
	switch val := v.(type) {
	case *StringValue:
		return val.Value, nil
	case *NumberValue:
		return val.String(), nil
	case *BoolValue:
		return val.String(), nil
	case *NullValue:
		return "null", nil
	default:
		return "", fmt.Errorf("cannot convert %s to string", v.Type())
	}
}

func ToNumber(v Value) (float64, error) {
	switch val := v.(type) {
	case *NumberValue:
		return val.Value, nil
	case *StringValue:
		return strconv.ParseFloat(val.Value, 64)
	case *BoolValue:
		if val.Value {
			return 1, nil
		}
		return 0, nil
	case *NullValue:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %s to number", v.Type())
	}
}

func ToBool(v Value) bool {
	return v.IsTruthy()
}

// Constructor helpers

func NewString(s string) *StringValue {
	return &StringValue{Value: s}
}

func NewNumber(n float64) *NumberValue {
	return &NumberValue{Value: n}
}

func NewBool(b bool) *BoolValue {
	return &BoolValue{Value: b}
}

func NewNull() *NullValue {
	return &NullValue{}
}

func NewArray(elements ...Value) *ArrayValue {
	return &ArrayValue{Elements: elements}
}

func NewObject() *ObjectValue {
	return &ObjectValue{Fields: make(map[string]Value)}
}

func NewValue(val any) Value {
	if val == nil {
		return NewNull()
	}

	switch v := val.(type) {
	case string:
		return NewString(v)
	case int:
		return NewNumber(float64(v))
	case int64:
		return NewNumber(float64(v))
	case float64:
		return NewNumber(v)
	case bool:
		return NewBool(v)
	case []any:
		arr := NewArray()
		for _, item := range v {
			arr.Elements = append(arr.Elements, NewValue(item))
		}
		return arr
	case map[string]any:
		obj := NewObject()
		for key, value := range v {
			obj.Set(key, NewValue(value))
		}
		return obj
	default:
		return newValueReflect(reflect.ValueOf(val))
	}
}

// newValueReflect converts a reflect.Value to a runtime Value
func newValueReflect(rv reflect.Value) Value {
	// Handle invalid or nil values
	if !rv.IsValid() {
		return NewNull()
	}

	// Dereference pointers
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return NewNull()
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return NewString(rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewNumber(float64(rv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewNumber(float64(rv.Uint()))
	case reflect.Float32, reflect.Float64:
		return NewNumber(rv.Float())
	case reflect.Bool:
		return NewBool(rv.Bool())
	case reflect.Slice, reflect.Array:
		arr := NewArray()
		for i := 0; i < rv.Len(); i++ {
			arr.Elements = append(arr.Elements, newValueReflect(rv.Index(i)))
		}
		return arr
	case reflect.Map:
		obj := NewObject()
		iter := rv.MapRange()
		for iter.Next() {
			key := iter.Key()
			// Convert key to string
			var keyStr string
			if key.Kind() == reflect.String {
				keyStr = key.String()
			} else {
				keyStr = fmt.Sprintf("%v", key.Interface())
			}
			obj.Set(keyStr, newValueReflect(iter.Value()))
		}
		return obj
	case reflect.Struct:
		obj := NewObject()
		t := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			field := t.Field(i)
			// Skip unexported fields
			if !field.IsExported() {
				continue
			}
			name := field.Name
			obj.Set(name, newValueReflect(rv.Field(i)))
		}
		return obj
	default:
		return NewNull()
	}
}
