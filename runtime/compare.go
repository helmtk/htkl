package runtime

import "fmt"

// Equal returns true if two values are equal
func Equal(left, right Value) bool {
	// Type must match
	if left.Type() != right.Type() {
		return false
	}

	// Compare values
	switch l := left.(type) {
	case *StringValue:
		r := right.(*StringValue)
		return l.Value == r.Value
	case *NumberValue:
		r := right.(*NumberValue)
		return l.Value == r.Value
	case *BoolValue:
		r := right.(*BoolValue)
		return l.Value == r.Value
	case *NullValue:
		return true
	default:
		// Arrays and objects are compared by reference
		return left == right
	}
}

// NotEqual returns true if two values are not equal
func NotEqual(left, right Value) bool {
	return !Equal(left, right)
}

// Less returns true if left < right (numeric comparison)
func Less(left, right Value) (bool, error) {
	leftNum, err := ToNumber(left)
	if err != nil {
		return false, fmt.Errorf("cannot compare %s and %s", left.Type(), right.Type())
	}
	rightNum, err := ToNumber(right)
	if err != nil {
		return false, fmt.Errorf("cannot compare %s and %s", left.Type(), right.Type())
	}
	return leftNum < rightNum, nil
}

// LessEqual returns true if left <= right (numeric comparison)
func LessEqual(left, right Value) (bool, error) {
	leftNum, err := ToNumber(left)
	if err != nil {
		return false, fmt.Errorf("cannot compare %s and %s", left.Type(), right.Type())
	}
	rightNum, err := ToNumber(right)
	if err != nil {
		return false, fmt.Errorf("cannot compare %s and %s", left.Type(), right.Type())
	}
	return leftNum <= rightNum, nil
}

// Greater returns true if left > right (numeric comparison)
func Greater(left, right Value) (bool, error) {
	leftNum, err := ToNumber(left)
	if err != nil {
		return false, fmt.Errorf("cannot compare %s and %s", left.Type(), right.Type())
	}
	rightNum, err := ToNumber(right)
	if err != nil {
		return false, fmt.Errorf("cannot compare %s and %s", left.Type(), right.Type())
	}
	return leftNum > rightNum, nil
}

// GreaterEqual returns true if left >= right (numeric comparison)
func GreaterEqual(left, right Value) (bool, error) {
	leftNum, err := ToNumber(left)
	if err != nil {
		return false, fmt.Errorf("cannot compare %s and %s", left.Type(), right.Type())
	}
	rightNum, err := ToNumber(right)
	if err != nil {
		return false, fmt.Errorf("cannot compare %s and %s", left.Type(), right.Type())
	}
	return leftNum >= rightNum, nil
}
