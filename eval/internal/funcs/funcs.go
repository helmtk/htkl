package funcs

import (
	"fmt"
	"math"
	"strings"

	"helmtk.dev/code/htkl/runtime"
)

// Registry holds all built-in functions
var Registry = map[string]runtime.Func{
	// String functions
	"upper":      upperFunc,
	"lower":      lowerFunc,
	"trim":       trimFunc,
	"quote":      quoteFunc,
	"nindent":    nindentFunc,
	"contains":   containsFunc,
	"trunc":      truncFunc,
	"trimSuffix": trimSuffixFunc,
	"replace":    replaceFunc,
	"printf":     printfFunc,

	// Conversion functions
	"toJson":   toJsonFunc,
	"toString": toStringFunc,

	// Utility functions
	"default":  defaultFunc,
	"len":      lenFunc,
	"has":      hasFunc,
	"coalesce": coalesceFunc,
	"empty":    emptyFunc,

	// Math functions
	"round": roundFunc,
	"floor": floorFunc,
	"ceil":  ceilFunc,

	// List functions
	"first":   firstFunc,
	"last":    lastFunc,
	"initial": initialFunc,
	"rest":    restFunc,
	"append":  appendFunc,
	"prepend": prependFunc,
	"concat":  concatFunc,
	"reverse": reverseFunc,
	"uniq":    uniqFunc,

	// String functions (additional)
	"split":      splitFunc,
	"join":       joinFunc,
	"trimPrefix": trimPrefixFunc,
	"hasPrefix":  hasPrefixFunc,
	"hasSuffix":  hasSuffixFunc,
	"repeat":     repeatFunc,

	// Dict/Object functions
	"keys":   keysFunc,
	"values": valuesFunc,
	"pick":   pickFunc,
	"omit":   omitFunc,
	"merge":  mergeFunc,
	"get":    getFunc,
	"set":    setFunc,

	// Encoding functions
	"b64enc": b64encFunc,
	"b64dec": b64decFunc,
}

// String functions

func upperFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upper expects 1 argument, got %d", len(args))
	}
	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}
	return runtime.NewString(strings.ToUpper(str)), nil
}

func lowerFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("lower expects 1 argument, got %d", len(args))
	}
	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}
	return runtime.NewString(strings.ToLower(str)), nil
}

func trimFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trim expects 1 argument, got %d", len(args))
	}
	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}
	return runtime.NewString(strings.TrimSpace(str)), nil
}

func quoteFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("quote expects 1 argument, got %d", len(args))
	}
	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}
	return runtime.NewString(fmt.Sprintf("%q", str)), nil
}

func nindentFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("nindent expects 2 arguments, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	n, err := runtime.ToNumber(args[1])
	if err != nil {
		return nil, err
	}

	indent := strings.Repeat(" ", int(n))
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}

	return runtime.NewString(strings.Join(lines, "\n")), nil
}

func containsFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contains expects 2 arguments, got %d", len(args))
	}

	str1, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	str2, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.Contains(str2, str1)), nil
}

func truncFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("trunc expects 2 arguments, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	maxLen, err := runtime.ToNumber(args[1])
	if err != nil {
		return nil, err
	}

	n := int(maxLen)
	if n < 0 {
		n = 0
	}
	if len(str) <= n {
		return runtime.NewString(str), nil
	}
	return runtime.NewString(str[:n]), nil
}

func trimSuffixFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("trimSuffix expects 2 arguments, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	suffix, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.TrimSuffix(str, suffix)), nil
}

func replaceFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace expects 3 arguments, got %d", len(args))
	}

	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	old, err := runtime.ToString(args[1])
	if err != nil {
		return nil, err
	}

	new, err := runtime.ToString(args[2])
	if err != nil {
		return nil, err
	}

	// Replace all occurrences (like Helm's replace function)
	return runtime.NewString(strings.ReplaceAll(str, old, new)), nil
}

func printfFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("printf expects at least 1 argument, got %d", len(args))
	}

	format, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}

	// Convert runtime values to interface{} for fmt.Sprintf
	fmtArgs := make([]interface{}, len(args)-1)
	for i, arg := range args[1:] {
		fmtArgs[i] = runtimeToNative(arg)
	}

	result := fmt.Sprintf(format, fmtArgs...)
	return runtime.NewString(result), nil
}

// Conversion functions

func toJsonFunc(args ...runtime.Value) (runtime.Value, error) {
	// TODO: implement toJson
	return nil, fmt.Errorf("toJson not yet implemented")
}

func toStringFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("toString expects 1 argument, got %d", len(args))
	}
	str, err := runtime.ToString(args[0])
	if err != nil {
		return nil, err
	}
	return runtime.NewString(str), nil
}

// Utility functions

func defaultFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("default expects 2 arguments, got %d", len(args))
	}

	// If first arg is null or falsy, return default
	if args[1].Type() == runtime.NullType || !args[1].IsTruthy() {
		return args[0], nil
	}

	return args[1], nil
}

func lenFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len expects 1 argument, got %d", len(args))
	}

	switch val := args[0].(type) {
	case *runtime.StringValue:
		return runtime.NewNumber(float64(len(val.Value))), nil
	case *runtime.ArrayValue:
		return runtime.NewNumber(float64(len(val.Elements))), nil
	case *runtime.ObjectValue:
		return runtime.NewNumber(float64(len(val.Fields))), nil
	default:
		return nil, fmt.Errorf("len does not support %s", val.Type())
	}
}

func hasFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("has expects 2 arguments, got %d", len(args))
	}

	query := args[0]
	arrVal := args[1]

	if arrVal == nil {
		return runtime.NewBool(false), nil
	}

	arr, ok := arrVal.(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("has expects first argument to be an array, got %s", arrVal.Type())
	}

	exists := false
	for _, el := range arr.Elements {
		if runtime.Equal(query, el) {
			exists = true
			break
		}
	}
	return runtime.NewBool(exists), nil
}

// Math functions

func roundFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("round expects 1 argument, got %d", len(args))
	}
	num, err := runtime.ToNumber(args[0])
	if err != nil {
		return nil, err
	}
	return runtime.NewNumber(math.Round(num)), nil
}

func floorFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("floor expects 1 argument, got %d", len(args))
	}
	num, err := runtime.ToNumber(args[0])
	if err != nil {
		return nil, err
	}
	return runtime.NewNumber(math.Floor(num)), nil
}

func ceilFunc(args ...runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ceil expects 1 argument, got %d", len(args))
	}
	num, err := runtime.ToNumber(args[0])
	if err != nil {
		return nil, err
	}
	return runtime.NewNumber(math.Ceil(num)), nil
}

// Helper function to convert runtime values to native Go types for YAML marshaling
func runtimeToNative(val runtime.Value) interface{} {
	switch v := val.(type) {
	case *runtime.StringValue:
		return v.Value
	case *runtime.NumberValue:
		return v.Value
	case *runtime.BoolValue:
		return v.Value
	case *runtime.NullValue:
		return nil
	case *runtime.ArrayValue:
		result := make([]interface{}, len(v.Elements))
		for i, elem := range v.Elements {
			result[i] = runtimeToNative(elem)
		}
		return result
	case *runtime.ObjectValue:
		result := make(map[string]interface{})
		for k, val := range v.Fields {
			result[k] = runtimeToNative(val)
		}
		return result
	default:
		return nil
	}
}
