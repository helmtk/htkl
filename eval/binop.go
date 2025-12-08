package eval

import (
	"fmt"

	"helmtk.dev/code/htkl/runtime"
)

// Arithmetic operations

func (e *evaluator) evalAdd(left, right runtime.Value) (runtime.Value, error) {
	// String concatenation
	if runtime.IsString(left) || runtime.IsString(right) {
		leftStr, err := runtime.ToString(left)
		if err != nil {
			return nil, err
		}
		rightStr, err := runtime.ToString(right)
		if err != nil {
			return nil, err
		}
		return runtime.NewString(leftStr + rightStr), nil
	}

	// Numeric addition
	leftNum, err := runtime.ToNumber(left)
	if err != nil {
		return nil, fmt.Errorf("cannot add %s and %s", left.Type(), right.Type())
	}
	rightNum, err := runtime.ToNumber(right)
	if err != nil {
		return nil, fmt.Errorf("cannot add %s and %s", left.Type(), right.Type())
	}
	return runtime.NewNumber(leftNum + rightNum), nil
}

func (e *evaluator) evalSub(left, right runtime.Value) (runtime.Value, error) {
	leftNum, err := runtime.ToNumber(left)
	if err != nil {
		return nil, fmt.Errorf("cannot subtract %s from %s", right.Type(), left.Type())
	}
	rightNum, err := runtime.ToNumber(right)
	if err != nil {
		return nil, fmt.Errorf("cannot subtract %s from %s", right.Type(), left.Type())
	}
	return runtime.NewNumber(leftNum - rightNum), nil
}

func (e *evaluator) evalMul(left, right runtime.Value) (runtime.Value, error) {
	leftNum, err := runtime.ToNumber(left)
	if err != nil {
		return nil, fmt.Errorf("cannot multiply %s and %s", left.Type(), right.Type())
	}
	rightNum, err := runtime.ToNumber(right)
	if err != nil {
		return nil, fmt.Errorf("cannot multiply %s and %s", left.Type(), right.Type())
	}
	return runtime.NewNumber(leftNum * rightNum), nil
}

func (e *evaluator) evalDiv(left, right runtime.Value) (runtime.Value, error) {
	leftNum, err := runtime.ToNumber(left)
	if err != nil {
		return nil, fmt.Errorf("cannot divide %s by %s", left.Type(), right.Type())
	}
	rightNum, err := runtime.ToNumber(right)
	if err != nil {
		return nil, fmt.Errorf("cannot divide %s by %s", left.Type(), right.Type())
	}
	if rightNum == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	return runtime.NewNumber(leftNum / rightNum), nil
}

// Comparison operations

func (e *evaluator) evalEqual(left, right runtime.Value) (runtime.Value, error) {
	return runtime.NewBool(runtime.Equal(left, right)), nil
}

func (e *evaluator) evalNotEqual(left, right runtime.Value) (runtime.Value, error) {
	return runtime.NewBool(runtime.NotEqual(left, right)), nil
}

func (e *evaluator) evalLess(left, right runtime.Value) (runtime.Value, error) {
	result, err := runtime.Less(left, right)
	if err != nil {
		return nil, err
	}
	return runtime.NewBool(result), nil
}

func (e *evaluator) evalLessEqual(left, right runtime.Value) (runtime.Value, error) {
	result, err := runtime.LessEqual(left, right)
	if err != nil {
		return nil, err
	}
	return runtime.NewBool(result), nil
}

func (e *evaluator) evalGreater(left, right runtime.Value) (runtime.Value, error) {
	result, err := runtime.Greater(left, right)
	if err != nil {
		return nil, err
	}
	return runtime.NewBool(result), nil
}

func (e *evaluator) evalGreaterEqual(left, right runtime.Value) (runtime.Value, error) {
	result, err := runtime.GreaterEqual(left, right)
	if err != nil {
		return nil, err
	}
	return runtime.NewBool(result), nil
}
