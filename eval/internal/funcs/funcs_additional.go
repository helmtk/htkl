package funcs

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/helmtk/htkl/runtime"
)

// Utility functions

func coalesceFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("coalesce expects at least 1 argument")
	}

	// Return the first non-null, non-empty value
	for _, arg := range args {
		if arg.Type() != runtime.NullType && arg.IsTruthy() {
			return arg, nil
		}
	}

	// If all are null/empty, return the last one
	return args[len(args)-1], nil
}

func emptyFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("empty expects 1 argument, got %d", len(args))
	}

	val := args[0]
	switch v := val.(type) {
	case *runtime.NullValue:
		return runtime.NewBool(true), nil
	case *runtime.StringValue:
		return runtime.NewBool(v.Value == ""), nil
	case *runtime.ArrayValue:
		return runtime.NewBool(len(v.Elements) == 0), nil
	case *runtime.ObjectValue:
		return runtime.NewBool(len(v.Fields) == 0), nil
	default:
		return runtime.NewBool(false), nil
	}
}

// List functions

func firstFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("first expects 1 argument, got %d", len(args))
	}

	arr, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("first expects an array, got %s", args[0].Type())
	}

	if len(arr.Elements) == 0 {
		return nil, fmt.Errorf("first: array is empty")
	}

	return arr.Elements[0], nil
}

func lastFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("last expects 1 argument, got %d", len(args))
	}

	arr, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("last expects an array, got %s", args[0].Type())
	}

	if len(arr.Elements) == 0 {
		return nil, fmt.Errorf("last: array is empty")
	}

	return arr.Elements[len(arr.Elements)-1], nil
}

func initialFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("initial expects 1 argument, got %d", len(args))
	}

	arr, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("initial expects an array, got %s", args[0].Type())
	}

	if len(arr.Elements) == 0 {
		return &runtime.ArrayValue{Elements: []runtime.Value{}}, nil
	}

	return &runtime.ArrayValue{Elements: arr.Elements[:len(arr.Elements)-1]}, nil
}

func restFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("rest expects 1 argument, got %d", len(args))
	}

	arr, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("rest expects an array, got %s", args[0].Type())
	}

	if len(arr.Elements) == 0 {
		return &runtime.ArrayValue{Elements: []runtime.Value{}}, nil
	}

	return &runtime.ArrayValue{Elements: arr.Elements[1:]}, nil
}

func appendFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("append expects at least 2 arguments, got %d", len(args))
	}

	arr, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("append expects first argument to be an array, got %s", args[0].Type())
	}

	result := make([]runtime.Value, len(arr.Elements))
	copy(result, arr.Elements)
	result = append(result, args[1:]...)

	return &runtime.ArrayValue{Elements: result}, nil
}

func prependFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("prepend expects at least 2 arguments, got %d", len(args))
	}

	arr, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("prepend expects first argument to be an array, got %s", args[0].Type())
	}

	result := make([]runtime.Value, 0, len(arr.Elements)+len(args)-1)
	result = append(result, args[1:]...)
	result = append(result, arr.Elements...)

	return &runtime.ArrayValue{Elements: result}, nil
}

func concatFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) == 0 {
		return &runtime.ArrayValue{Elements: []runtime.Value{}}, nil
	}

	var result []runtime.Value
	for _, arg := range args {
		arr, ok := arg.(*runtime.ArrayValue)
		if !ok {
			return nil, fmt.Errorf("concat expects all arguments to be arrays, got %s", arg.Type())
		}
		result = append(result, arr.Elements...)
	}

	return &runtime.ArrayValue{Elements: result}, nil
}

func reverseFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("reverse expects 1 argument, got %d", len(args))
	}

	arr, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("reverse expects an array, got %s", args[0].Type())
	}

	result := make([]runtime.Value, len(arr.Elements))
	for i, v := range arr.Elements {
		result[len(arr.Elements)-1-i] = v
	}

	return &runtime.ArrayValue{Elements: result}, nil
}

func uniqFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("uniq expects 1 argument, got %d", len(args))
	}

	arr, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("uniq expects an array, got %s", args[0].Type())
	}

	seen := make(map[string]bool)
	var result []runtime.Value

	for _, v := range arr.Elements {
		// Use string representation as key for deduplication
		str, _ := runtime.ToString(v)
		if !seen[str] {
			seen[str] = true
			result = append(result, v)
		}
	}

	return &runtime.ArrayValue{Elements: result}, nil
}

// String functions (additional)

func splitFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("split expects 2 arguments, got %d", len(args))
	}

	sep, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	str, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	parts := strings.Split(str, sep)
	elements := make([]runtime.Value, len(parts))
	for i, part := range parts {
		elements[i] = runtime.NewString(part)
	}

	return &runtime.ArrayValue{Elements: elements}, nil
}

func joinFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("join expects 2 arguments, got %d", len(args))
	}

	sep, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	arr, ok := args[1].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("join expects second argument to be an array, got %s", args[1].Type())
	}

	parts := make([]string, len(arr.Elements))
	for i, elem := range arr.Elements {
		parts[i], _ = runtime.ToString(elem)
	}

	return runtime.NewString(strings.Join(parts, sep)), nil
}

func trimPrefixFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("trimPrefix expects 2 arguments, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	prefix, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.TrimPrefix(str, prefix)), nil
}

func hasPrefixFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("hasPrefix expects 2 arguments, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	prefix, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.HasPrefix(str, prefix)), nil
}

func hasSuffixFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("hasSuffix expects 2 arguments, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	suffix, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.HasSuffix(str, suffix)), nil
}

func repeatFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("repeat expects 2 arguments, got %d", len(args))
	}

	count, err := runtime.ToNumber(args[0])
	if err != nil {
		return nil, err
	}

	str, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.Repeat(str, int(count))), nil
}

// Dict/Object functions

func keysFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("keys expects 1 argument, got %d", len(args))
	}

	obj, ok := args[0].(*runtime.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("keys expects an object, got %s", args[0].Type())
	}

	var result []runtime.Value
	for k := range obj.Fields {
		result = append(result, runtime.NewString(k))
	}

	return &runtime.ArrayValue{Elements: result}, nil
}

func valuesFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("values expects 1 argument, got %d", len(args))
	}

	obj, ok := args[0].(*runtime.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("values expects an object, got %s", args[0].Type())
	}

	var result []runtime.Value
	for _, v := range obj.Fields {
		result = append(result, v)
	}

	return &runtime.ArrayValue{Elements: result}, nil
}

func pickFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("pick expects at least 2 arguments, got %d", len(args))
	}

	obj, ok := args[0].(*runtime.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("pick expects first argument to be an object, got %s", args[0].Type())
	}

	result := runtime.NewObject()
	for i := 1; i < len(args); i++ {
		key, err := runtime.ToString(args[i])
		if err != nil {
			return nil, err
		}
		if val, ok := obj.Fields[key]; ok {
			result.Set(key, val)
		}
	}

	return result, nil
}

func omitFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("omit expects at least 2 arguments, got %d", len(args))
	}

	obj, ok := args[0].(*runtime.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("omit expects first argument to be an object, got %s", args[0].Type())
	}

	omitKeys := make(map[string]bool)
	for i := 1; i < len(args); i++ {
		key, err := runtime.ToString(args[i])
		if err != nil {
			return nil, err
		}
		omitKeys[key] = true
	}

	result := runtime.NewObject()
	for k, v := range obj.Fields {
		if !omitKeys[k] {
			result.Set(k, v)
		}
	}

	return result, nil
}

func mergeFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) == 0 {
		return runtime.NewObject(), nil
	}

	result := runtime.NewObject()
	for _, arg := range args {
		obj, ok := arg.(*runtime.ObjectValue)
		if !ok {
			return nil, fmt.Errorf("merge expects all arguments to be objects, got %s", arg.Type())
		}
		for k, v := range obj.Fields {
			result.Set(k, v)
		}
	}

	return result, nil
}

func getFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("get expects 2 arguments, got %d", len(args))
	}

	obj, ok := args[0].(*runtime.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("get expects first argument to be an object, got %s", args[0].Type())
	}

	key, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	if val, ok := obj.Fields[key]; ok {
		return val, nil
	}

	return runtime.NewNull(), nil
}

func setFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("set expects 3 arguments, got %d", len(args))
	}

	obj, ok := args[0].(*runtime.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("set expects first argument to be an object, got %s", args[0].Type())
	}

	key, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	// Create a copy to avoid mutating the original
	result := runtime.NewObject()
	for k, v := range obj.Fields {
		result.Set(k, v)
	}
	result.Set(key, args[2])

	return result, nil
}

// Encoding functions

func b64encFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("b64enc expects 1 argument, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return runtime.NewString(encoded), nil
}

func b64decFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("b64dec expects 1 argument, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("b64dec: %w", err)
	}

	return runtime.NewString(string(decoded)), nil
}
