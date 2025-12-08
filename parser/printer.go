package parser

import (
	"fmt"
	"io"
	"strings"
)

// Printer formats an AST for display
type Printer struct {
	indent int
	w      io.Writer
}

// NewPrinter creates a new AST printer
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

func (p *Printer) println(format string, args ...interface{}) {
	fmt.Fprintf(p.w, strings.Repeat("  ", p.indent)+format+"\n", args...)
}

// PrintDocument prints a Document node
func (p *Printer) PrintDocument(doc *Document) {
	p.println("Document")
	p.indent++

	// Print definitions first, then body statements
	stmtIdx := 0
	for _, def := range doc.Definitions {
		p.println("Statement[%d]:", stmtIdx)
		p.indent++
		p.PrintDefineStatement(def)
		p.indent--
		stmtIdx++
	}
	for _, stmt := range doc.Body {
		p.println("Statement[%d]:", stmtIdx)
		p.indent++
		p.PrintStatement(stmt)
		p.indent--
		stmtIdx++
	}
	p.indent--
}

// PrintStatement prints a Statement node
func (p *Printer) PrintStatement(stmt Statement) {
	switch s := stmt.(type) {
	case *KeyValueStatement:
		p.PrintKeyValue(s)
	case *LetStatement:
		p.PrintLetStatement(s)
	default:
		p.println("Unknown statement: %T", s)
	}
}

// PrintKeyValue prints a KeyValue node
func (p *Printer) PrintKeyValue(kv *KeyValueStatement) {
	p.println("KeyValue")
	p.indent++
	p.println("Key: %q", kv.Key)
	p.println("Value:")
	p.indent++
	p.PrintValueStatement(kv.Value)
	p.indent--
	p.indent--
}

// PrintLetStatement prints a LetStatement node
func (p *Printer) PrintLetStatement(let *LetStatement) {
	p.println("LetStatement")
	p.indent++
	p.println("Name: %q", let.Name)
	p.println("Value:")
	p.indent++
	p.PrintValueStatement(let.Value)
	p.indent--
	p.indent--
}

// PrintDefineStatement prints a DefineStatement node
func (p *Printer) PrintDefineStatement(def *Definition) {
	p.println("DefineStatement")
	p.indent++
	p.println("Name: %q", def.Name)
	p.println("Body:")
	p.indent++
	for i, val := range def.Body {
		p.println("Value[%d]:", i)
		p.indent++
		p.PrintNode(val)
		p.indent--
	}
	p.indent--
	p.indent--
}

// PrintComment prints a Comment node
func (p *Printer) PrintComment(c *Comment) {
	p.println("Comment: %q", c.Text)
}

// PrintBinaryOp prints a BinaryOp node
func (p *Printer) PrintBinaryOp(b *BinaryOp) {
	p.println("BinaryOp")
	p.indent++
	p.println("Operator: %q", b.Operator)
	p.println("Left:")
	p.indent++
	p.PrintValue(b.Left)
	p.indent--
	p.println("Right:")
	p.indent++
	p.PrintValue(b.Right)
	p.indent--
	p.indent--
}

// PrintUnaryOp prints a UnaryOp node
func (p *Printer) PrintUnaryOp(u *UnaryOp) {
	p.println("UnaryOp")
	p.indent++
	p.println("Operator: %q", u.Operator)
	p.println("Operand:")
	p.indent++
	p.PrintValue(u.Operand)
	p.indent--
	p.indent--
}

// PrintMemberExpression prints a MemberExpression node
func (p *Printer) PrintMemberExpression(m *MemberExpression) {
	p.println("MemberExpression")
	p.indent++
	p.println("Object:")
	p.indent++
	p.PrintValue(m.Object)
	p.indent--
	p.println("Member: %q", m.Member)
	p.indent--
}

// PrintIndexExpression prints an IndexExpression node
func (p *Printer) PrintIndexExpression(idx *IndexExpression) {
	p.println("IndexExpression")
	p.indent++
	p.println("Object:")
	p.indent++
	p.PrintValue(idx.Object)
	p.indent--
	p.println("Index:")
	p.indent++
	p.PrintValue(idx.Index)
	p.indent--
	p.indent--
}

// PrintInterpolatedString prints an InterpolatedString node
func (p *Printer) PrintInterpolatedString(s *InterpolatedString) {
	p.println("InterpolatedString")
	p.indent++
	for i, part := range s.Parts {
		p.println("Part[%d]:", i)
		p.indent++
		p.PrintValue(part)
		p.indent--
	}
	p.indent--
}

// PrintIncludeExpression prints an IncludeExpression node
func (p *Printer) PrintIncludeExpression(inc *IncludeExpression) {
	p.println("IncludeExpression")
	p.indent++
	p.println("Name: %q", inc.Name)
	if inc.Context != nil {
		p.println("Content:")
		p.indent++
		p.PrintValue(inc.Context)
		p.indent--
	} else {
		p.println("Args: []")
	}
	p.indent--
}

// PrintCallExpression prints a CallExpression node
func (p *Printer) PrintCallExpression(call *CallExpression) {
	p.println("CallExpression")
	p.indent++
	p.println("Function:")
	p.indent++
	p.PrintValue(call.Function)
	p.indent--
	if len(call.Args) > 0 {
		p.println("Args:")
		p.indent++
		for i, arg := range call.Args {
			p.println("Arg[%d]:", i)
			p.indent++
			p.PrintValue(arg)
			p.indent--
		}
		p.indent--
	} else {
		p.println("Args: []")
	}
	p.indent--
}

func (p *Printer) PrintValueStatement(vs ValueStatement) {
	switch v := vs.(type) {
	case *IfStatement:
		p.PrintIfStatement(v)
	case *ForStatement:
		p.PrintForStatement(v)
	case *WithStatement:
		p.PrintWithStatement(v)
	case *IncludeExpression:
		p.PrintIncludeExpression(v)
	case Expression:
		p.PrintValue(v)
	}
}

// PrintValue prints an Expression node
func (p *Printer) PrintValue(val Expression) {
	switch v := val.(type) {
	case *StringLiteral:
		p.println("StringLiteral: %q", v.Value)
	case *InterpolatedString:
		p.PrintInterpolatedString(v)
	case *NumberLiteral:
		p.println("NumberLiteral: %v", v.Value)
	case *BooleanLiteral:
		p.println("BooleanLiteral: %v", v.Value)
	case *NullLiteral:
		p.println("NullLiteral")
	case *CurrentContext:
		p.println("CurrentContext")
	case *Identifier:
		p.println("Identifier: %s", v.Name)
	case *MemberExpression:
		p.PrintMemberExpression(v)
	case *IndexExpression:
		p.PrintIndexExpression(v)
	case *BinaryOp:
		p.PrintBinaryOp(v)
	case *UnaryOp:
		p.PrintUnaryOp(v)
	case *CallExpression:
		p.PrintCallExpression(v)
	case *Object:
		p.PrintObject(v)
	case *Array:
		p.PrintArray(v)
	default:
		p.println("Unknown expression: %T", v)
	}
}

// PrintSpreadElement prints a SpreadElement node
func (p *Printer) PrintSpreadElement(s *SpreadStatement) {
	p.println("SpreadElement")
	p.indent++
	p.println("Operand:")
	p.indent++
	p.PrintValueStatement(s.Operand)
	p.indent--
	p.indent--
}

// PrintNode prints any Node
func (p *Printer) PrintNode(node Node) {
	switch n := node.(type) {
	case *KeyValueStatement:
		p.PrintKeyValue(n)
	case *LetStatement:
		p.PrintLetStatement(n)
	case *SpreadStatement:
		p.PrintSpreadElement(n)
	case *IfStatement:
		p.PrintIfStatement(n)
	case *ForStatement:
		p.PrintForStatement(n)
	case *WithStatement:
		p.PrintWithStatement(n)
	case *BreakStatement:
		p.println("BreakStatement")
	case *ContinueStatement:
		p.println("ContinueStatement")
	case Expression:
		p.PrintValue(n)
	default:
		p.println("Unknown node: %T", n)
	}
}

// PrintObject prints an Object node
func (p *Printer) PrintObject(obj *Object) {
	p.println("Object")
	p.indent++
	for i, field := range obj.Body {
		p.println("Field[%d]:", i)
		p.indent++
		p.PrintNode(field)
		p.indent--
	}
	p.indent--
}

// PrintArray prints an Array node
func (p *Printer) PrintArray(arr *Array) {
	p.println("Array")
	p.indent++
	for i, elem := range arr.Body {
		p.println("Element[%d]:", i)
		p.indent++
		p.PrintNode(elem)
		p.indent--
	}
	p.indent--
}

// PrintIfStatement prints an IfStatement node
func (p *Printer) PrintIfStatement(ifStmt *IfStatement) {
	p.println("IfStatement")
	p.indent++
	p.println("Condition:")
	p.indent++
	p.PrintValue(ifStmt.Condition)
	p.indent--
	p.println("Body:")
	p.indent++
	for i, val := range ifStmt.Body {
		p.println("Value[%d]:", i)
		p.indent++
		p.PrintNode(val)
		p.indent--
	}
	p.indent--
	if len(ifStmt.Else) > 0 {
		p.println("Else:")
		p.indent++
		for i, val := range ifStmt.Else {
			p.println("Value[%d]:", i)
			p.indent++
			p.PrintNode(val)
			p.indent--
		}
		p.indent--
	}
	p.indent--
}

// PrintWithStatement prints a WithStatement node
func (p *Printer) PrintWithStatement(withStmt *WithStatement) {
	p.println("WithStatement")
	p.indent++
	p.println("Context:")
	p.indent++
	p.PrintValue(withStmt.Context)
	p.indent--
	p.println("Body:")
	p.indent++
	for i, val := range withStmt.Body {
		p.println("Value[%d]:", i)
		p.indent++
		p.PrintNode(val)
		p.indent--
	}
	p.indent--
	p.indent--
}

// PrintForStatement prints a ForStatement node
func (p *Printer) PrintForStatement(forStmt *ForStatement) {
	p.println("ForStatement")
	p.indent++
	p.println("KeyVar: %q", forStmt.KeyVar)
	p.println("ValueVar: %q", forStmt.ValueVar)
	p.println("Iterable:")
	p.indent++
	p.PrintValue(forStmt.Iterable)
	p.indent--
	p.println("Body:")
	p.indent++
	for i, val := range forStmt.Body {
		p.println("Value[%d]:", i)
		p.indent++
		p.PrintNode(val)
		p.indent--
	}
	p.indent--
	p.indent--
}
