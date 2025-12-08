package eval

import (
	"github.com/helmtk/htkl/parser"
	"github.com/helmtk/htkl/runtime"
)

// evalStringLiteral evaluates a string literal
func evalStringLiteral(n *parser.StringLiteral) (runtime.Value, error) {
	return runtime.NewString(n.Value), nil
}

// evalNumberLiteral evaluates a number literal
func evalNumberLiteral(n *parser.NumberLiteral) (runtime.Value, error) {
	return runtime.NewNumber(n.Value), nil
}

// evalBooleanLiteral evaluates a boolean literal
func evalBooleanLiteral(n *parser.BooleanLiteral) (runtime.Value, error) {
	return runtime.NewBool(n.Value), nil
}

// evalNullLiteral evaluates a null literal
func evalNullLiteral(n *parser.NullLiteral) (runtime.Value, error) {
	return runtime.NewNull(), nil
}
