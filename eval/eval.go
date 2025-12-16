package eval

import (
	"errors"
	"fmt"

	"helmtk.dev/code/htkl/parser"
	"helmtk.dev/code/htkl/runtime"
)

// EvalDocument evaluates a complete helmtk document
// Returns an ArrayValue containing all root-level documents
func EvalDocument(doc *parser.Document, root *runtime.Scope) (runtime.Value, error) {

	docColl := &documentCollector{}
	e := evaluator{
		scope: root,
		coll:  docColl,
	}

	// process all "define" blocks to register templates
	for _, def := range doc.Definitions {
		// Get filename from the body nodes
		filename := ""
		if len(def.Body) > 0 {
			filename = def.Body[0].GetPos().Filename
		}

		// Create template with filename for better error messages
		tmpl := runtime.NewTemplate(def.Name, def.Body, filename)

		// Register it in the scope
		e.scope.DefineTemplate(def.Name, tmpl)
	}

	// evaluate all statements in the document context
	for _, stmt := range doc.Body {
		if err := e.evalStatement(stmt); err != nil {
			return nil, err
		}
	}

	// Return array of documents
	arr := &runtime.ArrayValue{
		Elements: docColl.documents,
	}
	return arr, nil
}

// evaluator evaluates AST nodes into runtime values
type evaluator struct {
	scope *runtime.Scope
	coll  any
}

// Eval evaluates an AST value node and returns a runtime value
func (e *evaluator) evalExpression(node parser.Expression) (runtime.Value, error) {
	switch n := node.(type) {

	case *parser.StringLiteral:
		return evalStringLiteral(n)
	case *parser.NumberLiteral:
		return evalNumberLiteral(n)
	case *parser.BooleanLiteral:
		return evalBooleanLiteral(n)
	case *parser.NullLiteral:
		return evalNullLiteral(n)
	case *parser.InterpolatedString:
		return e.evalInterpolatedString(n)
	case *parser.Identifier:
		return e.evalIdentifier(n)
	case *parser.BinaryOp:
		return e.evalBinaryOp(n)
	case *parser.UnaryOp:
		return e.evalUnaryOp(n)
	case *parser.CallExpression:
		return e.evalCallExpression(n)
	case *parser.MemberExpression:
		return e.evalMemberExpression(n)
	case *parser.IndexExpression:
		return e.evalIndexExpression(n)
	case *parser.Array:
		return e.evalArray(n)
	case *parser.Object:
		return e.evalObject(n)
	case *parser.IncludeExpression:
		return e.collectSingleValue(node, func(sub *evaluator) error {
			return sub.evalIncludeStatement(n)
		})
	case *parser.CurrentContext:
		return e.evalCurrentContext(n)
	default:
		return nil, errorf(n.GetPos(), "unsupported node type: %T", node)
	}
}

func (e *evaluator) evalStatement(node parser.Statement) error {
	switch n := node.(type) {
	case *parser.LetStatement:
		return e.evalLetStatement(n)
	case *parser.AssignmentStatement:
		return e.evalAssignmentStatement(n)
	case *parser.WithStatement:
		return e.evalWithStatement(n)
	case *parser.ForStatement:
		return e.evalForStatement(n)
	case *parser.KeyValueStatement:
		return e.evalKeyValue(n)
	case *parser.SpreadStatement:
		return e.evalSpreadStatement(n)
	case *parser.IfStatement:
		return e.evalIfStatement(n)
	case *parser.IncludeExpression:
		return e.evalIncludeStatement(n)
	case parser.Expression:
		// Evaluate the expression
		val, err := e.evalExpression(n)
		if err != nil {
			return err
		}

		// If we're in a document collector context, add it as a document
		if docColl, ok := e.coll.(*documentCollector); ok {
			docColl.addDocument(val)
		}

		return nil
	default:
		return errorf(n.GetPos(), "unsupported statement: %T", node)
	}
}

func (e *evaluator) evalKeyValue(n *parser.KeyValueStatement) error {
	// Check if we're in a document collector - if so, we need an implicit root object
	if docColl, ok := e.coll.(*documentCollector); ok {
		// Create an implicit root object if we encounter key:value at document level
		// Check if the last document is an object we can add to
		var obj *runtime.ObjectValue
		if len(docColl.documents) > 0 {
			if lastObj, ok := docColl.documents[len(docColl.documents)-1].(*runtime.ObjectValue); ok {
				obj = lastObj
			}
		}

		// If we don't have an object yet, create one
		if obj == nil {
			obj = &runtime.ObjectValue{}
			docColl.addDocument(obj)
		}

		// Evaluate the value
		val, err := e.evalValueStatement(n.Value)
		if err != nil {
			return err
		}

		obj.Set(n.Key, val)
		return nil
	}

	// Normal object context
	obj, ok := e.coll.(*runtime.ObjectValue)
	if !ok {
		return errorf(n.Pos, "key:value pair in non-object context")
	}

	val, err := e.evalValueStatement(n.Value)
	if err != nil {
		return err
	}

	obj.Set(n.Key, val)
	return nil
}

func (e *evaluator) evalValueStatement(node parser.ValueStatement) (runtime.Value, error) {
	switch it := node.(type) {
	case *parser.IfStatement:
		return e.collectSingleValue(node, func(sub *evaluator) error {
			return sub.evalIfStatement(it)
		})
	case *parser.WithStatement:
		return e.collectSingleValue(node, func(sub *evaluator) error {
			return sub.evalWithStatement(it)
		})
	case parser.Expression:
		return e.evalExpression(it)
	default:
		return nil, errorf(node.GetPos(), "unexpected node %T", node)
	}
}

// evalArray evaluates an array literal
func (e *evaluator) evalArray(node *parser.Array) (runtime.Value, error) {
	arr := &runtime.ArrayValue{}
	sub := evaluator{scope: e.scope, coll: arr}

	for _, item := range node.Body {
		if err := sub.collectNode(item); err != nil {
			return nil, err
		}
	}

	return arr, nil
}

// evalObject evaluates an object literal
func (e *evaluator) evalObject(node *parser.Object) (runtime.Value, error) {
	obj := &runtime.ObjectValue{}
	sub := evaluator{scope: e.scope, coll: obj}

	for _, item := range node.Body {

		it, ok := item.(parser.Statement)
		if !ok {
			return nil, errorf(it.GetPos(), "unsupported node: %T", it)
		}

		if err := sub.evalStatement(it); err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func (e *evaluator) collectNode(node parser.Node) error {
	switch it := node.(type) {

	case parser.Expression:
		val, err := e.evalExpression(it)
		if err != nil {
			return err
		}

		switch coll := e.coll.(type) {
		case *runtime.ArrayValue:
			coll.Elements = append(coll.Elements, val)
		case *singleValueCollector:
			return coll.setVal(val)
		case *documentCollector:
			// Root-level expressions (typically object literals) become documents
			coll.addDocument(val)
		default:
			return errorf(it.GetPos(), "unexpected value")
		}
		return nil

	case parser.Statement:
		return e.evalStatement(it)

	default:
		return errorf(node.GetPos(), "unsupported node: %T", it)
	}
}

func (e *evaluator) evalIfStatement(n *parser.IfStatement) error {
	// Evaluate the condition
	cond, err := e.evalExpression(n.Condition)
	if err != nil {
		return err
	}

	// Determine which branch to execute
	var branch []parser.Node
	if cond.IsTruthy() {
		branch = n.Body
	} else {
		branch = n.Else
	}

	// Emit all items from the branch
	for _, item := range branch {
		if err := e.collectNode(item); err != nil {
			return err
		}
	}

	return nil
}

// documentCollector collects root-level objects as separate documents
type documentCollector struct {
	documents []runtime.Value
}

func (d *documentCollector) addDocument(val runtime.Value) {
	d.documents = append(d.documents, val)
}

type singleValueCollector struct {
	val runtime.Value
}

func (s *singleValueCollector) setVal(v runtime.Value) error {
	if s.val != nil {
		return fmt.Errorf("unexpected value, expected only a single value")
	}
	s.val = v
	return nil
}

func (e *evaluator) collectSingleValue(n parser.Node, cb func(*evaluator) error) (runtime.Value, error) {

	coll := &singleValueCollector{}
	sub := &evaluator{scope: e.scope, coll: coll}

	if err := cb(sub); err != nil {
		return nil, err
	}

	if coll.val == nil {
		return nil, errorf(n.GetPos(), "expected value")
	}

	return coll.val, nil
}

func (e *evaluator) evalWithStatement(n *parser.WithStatement) error {
	// Evaluate the context
	context, err := e.evalExpression(n.Context)
	if err != nil {
		return err
	}

	// Create new scope for with body and bind the context to the variable
	newScope := runtime.NewScope(e.scope)
	newScope.Set(n.VarName, context)

	sub := evaluator{
		scope: newScope,
		coll:  e.coll,
	}

	// Emit all items from the body
	for _, item := range n.Body {
		if err := sub.collectNode(item); err != nil {
			return err
		}
	}

	return nil
}

func (e *evaluator) evalSpreadStatement(n *parser.SpreadStatement) error {
	// Evaluate the operand
	val, err := e.evalValueStatement(n.Operand)
	if err != nil {
		return err
	}

	// Spread into the current collection
	switch coll := e.coll.(type) {
	case *runtime.ArrayValue:
		// Spread array into array
		arr, ok := val.(*runtime.ArrayValue)
		if !ok {
			return errorf(n.Pos, "cannot spread %s into array", val.Type())
		}
		coll.Elements = append(coll.Elements, arr.Elements...)

	case *runtime.ObjectValue:
		// Spread object into object
		obj, ok := val.(*runtime.ObjectValue)
		if !ok {
			return errorf(n.Pos, "cannot spread %s into object", val.Type())
		}
		for k, v := range obj.Fields {
			coll.Set(k, v)
		}

	default:
		return errorf(n.Pos, "cannot spread into this context")
	}

	return nil
}

func (e *evaluator) evalForStatement(n *parser.ForStatement) error {
	// Evaluate the iterable
	iterable, err := e.evalExpression(n.Iterable)
	if err != nil {
		return err
	}

	switch iter := iterable.(type) {
	case *runtime.ArrayValue:
		for i, elem := range iter.Elements {
			key := runtime.NewNumber(float64(i))
			err := e.evalForIteration(n, key, elem)
			if err == breakSignal {
				break
			}
			if err != nil {
				return err
			}
		}

	case *runtime.ObjectValue:
		for key, val := range iter.Fields {
			key := runtime.NewString(key)
			err := e.evalForIteration(n, key, val)
			if err == breakSignal {
				break
			}
			if err != nil {
				return err
			}
		}

	default:
		return errorf(n.Iterable.GetPos(), "cannot iterate over %s", iterable.Type())
	}

	return nil
}

var breakSignal = errors.New("break")

// evalForIteration evaluates a single iteration of a for loop
func (e *evaluator) evalForIteration(n *parser.ForStatement, key, value runtime.Value) error {
	// Create new scope for loop variables
	loopScope := runtime.NewScope(e.scope)
	sub := &evaluator{scope: loopScope, coll: e.coll}

	// Bind loop variables
	if n.KeyVar != "" {
		loopScope.Set(n.KeyVar, key)
	}
	loopScope.Set(n.ValueVar, value)

	// Emit all items from the body
	for _, item := range n.Body {
		switch item.(type) {
		case *parser.BreakStatement:
			return breakSignal
		case *parser.ContinueStatement:
			break
		}

		if err := sub.collectNode(item); err != nil {
			return err
		}
	}

	return nil
}

// evalMemberExpression evaluates member access (e.g., obj.field)
func (e *evaluator) evalMemberExpression(n *parser.MemberExpression) (runtime.Value, error) {
	// Evaluate the object
	objVal, err := e.evalExpression(n.Object)
	if err != nil {
		return nil, err
	}

	// If the object is null, return null (allows chaining through null values)
	// This matches Helm's behavior where undefined.field returns empty/null
	if objVal.Type() == runtime.NullType {
		return runtime.NewNull(), nil
	}

	// It must be an object
	obj, ok := objVal.(*runtime.ObjectValue)
	if !ok {
		return nil, errorf(n.Pos, "cannot access member of %s", objVal.Type())
	}

	// Get the field
	val, ok := obj.Get(n.Member)
	if !ok {
		// Return null for undefined fields instead of erroring
		// This matches Helm's behavior where undefined values are treated as empty/null
		return runtime.NewNull(), nil
	}

	return val, nil
}

// evalIndexExpression evaluates array/object indexing (e.g., arr[0], obj["key"])
func (e *evaluator) evalIndexExpression(n *parser.IndexExpression) (runtime.Value, error) {
	// Evaluate the object/array
	objVal, err := e.evalExpression(n.Object)
	if err != nil {
		return nil, err
	}

	// Evaluate the index
	indexVal, err := e.evalExpression(n.Index)
	if err != nil {
		return nil, err
	}

	switch obj := objVal.(type) {
	case *runtime.ArrayValue:
		// Index must be a number
		num, ok := indexVal.(*runtime.NumberValue)
		if !ok {
			return nil, errorf(n.Pos, "array index must be a number, got %s", indexVal.Type())
		}

		idx := int(num.Value)
		if idx < 0 || idx >= len(obj.Elements) {
			return nil, errorf(n.Pos, "array index out of bounds: %d", idx)
		}

		return obj.Elements[idx], nil

	case *runtime.ObjectValue:
		// Index must be a string
		key, err := runtime.ToString(indexVal)
		if err != nil {
			return nil, errorf(n.Pos, "object index must be a string")
		}

		val, ok := obj.Get(key)
		if !ok {
			return nil, errorf(n.Pos, "undefined field: %s", key)
		}

		return val, nil

	default:
		return nil, errorf(n.Pos, "cannot index %s", objVal.Type())
	}
}

// evalBinaryOp evaluates a binary operation
func (e *evaluator) evalBinaryOp(n *parser.BinaryOp) (runtime.Value, error) {
	// Handle pipe operator specially
	if n.Operator == "|" {
		return e.evalPipe(n)
	}

	// Evaluate left and right operands
	left, err := e.evalExpression(n.Left)
	if err != nil {
		return nil, err
	}

	right, err := e.evalExpression(n.Right)
	if err != nil {
		return nil, err
	}

	// Dispatch based on operator
	switch n.Operator {
	// Arithmetic operators
	case "+":
		return e.evalAdd(left, right)
	case "-":
		return e.evalSub(left, right)
	case "*":
		return e.evalMul(left, right)
	case "/":
		return e.evalDiv(left, right)

	// Comparison operators
	case "==":
		return e.evalEqual(left, right)
	case "!=":
		return e.evalNotEqual(left, right)
	case "<":
		return e.evalLess(left, right)
	case "<=":
		return e.evalLessEqual(left, right)
	case ">":
		return e.evalGreater(left, right)
	case ">=":
		return e.evalGreaterEqual(left, right)

	// Logical operators
	case "&&":
		return runtime.NewBool(left.IsTruthy() && right.IsTruthy()), nil
	case "||":
		return runtime.NewBool(left.IsTruthy() || right.IsTruthy()), nil

	default:
		return nil, errorf(n.Pos, "unknown operator: %s", n.Operator)
	}
}

// evalPipe evaluates the pipe operator
func (e *evaluator) evalPipe(n *parser.BinaryOp) (runtime.Value, error) {
	// Evaluate the left side (the value being piped)
	val, err := e.evalExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// The right side should be either an identifier or a function call
	switch right := n.Right.(type) {
	case *parser.Identifier:
		// Simple pipe: val | funcName
		// Call the function with val as the last argument (Go template behavior)
		return e.callFunction(right.Pos, right.Name, []runtime.Value{val})

	case *parser.CallExpression:
		// Pipe with function call: val | funcName(arg1, arg2)
		// Append val to the arguments (Go template behavior)
		funcName, ok := right.Function.(*parser.Identifier)
		if !ok {
			return nil, errorf(n.Pos, "pipe right side must be a function name")
		}

		// Evaluate the arguments, then append the piped value as the last argument
		var args []runtime.Value
		for _, arg := range right.Args {
			argVal, err := e.evalExpression(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, argVal)
		}
		args = append(args, val)

		return e.callFunction(right.Function.GetPos(), funcName.Name, args)

	default:
		return nil, errorf(n.Pos, "invalid pipe right side: %T", n.Right)
	}
}

// callFunction is a helper for calling functions
func (e *evaluator) callFunction(pos parser.Pos, name string, args []runtime.Value) (runtime.Value, error) {
	// Look up the function in the registry
	fn, ok := e.scope.GetFunction(name)
	if !ok {
		return nil, errorf(pos, "undefined function: %s", name)
	}

	// Call the function
	res, err := fn(args...)
	if err != nil {
		return nil, errorf(pos, "%s", err)
	}
	return res, nil
}

// evalUnaryOp evaluates a unary operation
func (e *evaluator) evalUnaryOp(n *parser.UnaryOp) (runtime.Value, error) {
	operand, err := e.evalExpression(n.Operand)
	if err != nil {
		return nil, err
	}

	switch n.Operator {
	case "!":
		// Logical not
		return runtime.NewBool(!operand.IsTruthy()), nil

	case "-":
		// Negation
		num, err := runtime.ToNumber(operand)
		if err != nil {
			return nil, errorf(n.Pos, "cannot negate %s", operand.Type())
		}
		return runtime.NewNumber(-num), nil

	default:
		return nil, errorf(n.Pos, "unknown unary operator: %s", n.Operator)
	}
}

func (e *evaluator) evalCallExpression(n *parser.CallExpression) (runtime.Value, error) {
	// Get the function name
	funcName, ok := n.Function.(*parser.Identifier)
	if !ok {
		return nil, errorf(n.Pos, "function must be an identifier")
	}

	// Evaluate arguments
	args := make([]runtime.Value, len(n.Args))
	for i, arg := range n.Args {
		val, err := e.evalExpression(arg)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	// Call the function
	return e.callFunction(n.Pos, funcName.Name, args)
}

func (e *evaluator) evalIncludeStatement(n *parser.IncludeExpression) error {
	// Get the template
	tmpl, err := e.scope.GetTemplate(n.Name)
	if err != nil {
		return errorf(n.Pos, "%s", err.Error())
	}

	// Create new scope for template evaluation
	tmplScope := runtime.NewScope(nil)
	tmplScope.Link(e.scope)

	if n.Context != nil {
		val, err := e.evalExpression(n.Context)
		if err != nil {
			return err
		}
		obj, ok := val.(*runtime.ObjectValue)
		if !ok {
			return errorf(n.Context.GetPos(), "template context must be an object")
		}
		for k, v := range obj.Fields {
			tmplScope.Set(k, v)
		}
	}

	tmplEval := &evaluator{
		scope: tmplScope,
		coll:  e.coll,
	}

	for _, node := range tmpl.Body {
		if err := tmplEval.collectNode(node); err != nil {
			return errorf(n.Pos, "include %q: %s", n.Name, err)
		}
	}

	return nil
}

func (e *evaluator) evalAssignmentStatement(n *parser.AssignmentStatement) error {
	// Evaluate the value
	val, err := e.evalValueStatement(n.Value)
	if err != nil {
		return err
	}

	// Update the variable in the current scope
	// Unlike let, assignment should update an existing variable
	e.scope.Set(n.Name, val)

	// Assignment statements don't produce a value
	return nil
}

func (e *evaluator) evalLetStatement(n *parser.LetStatement) error {
	// Evaluate the value
	val, err := e.evalValueStatement(n.Value)
	if err != nil {
		return err
	}

	// Bind it in the current scope
	e.scope.Set(n.Name, val)

	// Let statements don't produce a value
	return nil
}

// evalInterpolatedString evaluates an interpolated string with ${} expressions
func (e *evaluator) evalInterpolatedString(n *parser.InterpolatedString) (runtime.Value, error) {
	var result string
	for _, part := range n.Parts {
		val, err := e.evalExpression(part)
		if err != nil {
			return nil, err
		}
		str, err := runtime.ToString(val)
		if err != nil {
			return nil, wraperr(n.Pos, err)
		}
		result += str
	}
	return runtime.NewString(result), nil
}

// evalIdentifier looks up an identifier in the current scope
func (e *evaluator) evalIdentifier(n *parser.Identifier) (runtime.Value, error) {
	val, err := e.scope.Get(n.Name)
	if err != nil {
		return nil, errorf(n.Pos, "%s", err.Error())
	}
	return val, nil
}

func (e *evaluator) evalCurrentContext(_ *parser.CurrentContext) (runtime.Value, error) {
	// Try to get the implicit context from the scope
	// In templates, this will be the parameter passed to the template
	// For now, we'll create a default context with Release, Chart, and Values
	ctx := runtime.NewObject()

	// Copy Release, Chart, Values from current scope if they exist
	if release, err := e.scope.Get("Release"); err == nil {
		ctx.Set("Release", release)
	}
	if chart, err := e.scope.Get("Chart"); err == nil {
		ctx.Set("Chart", chart)
	}
	if values, err := e.scope.Get("Values"); err == nil {
		ctx.Set("Values", values)
	}

	return ctx, nil
}
