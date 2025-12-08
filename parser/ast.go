package parser

// Node represents a node in the abstract syntax tree
type Node interface {
	GetPos() Pos
	node()
}

type Pos struct {
	Filename string
	Line     int
	Col      int
}

// Document represents the root of a helmtk template
type Document struct {
	Body        []Statement
	Definitions []*Definition
}

func (d *Document) node()       {}
func (d *Document) GetPos() Pos { return Pos{} }

// Statement represents a top-level statement
type Statement interface {
	Node
	statement()
}

type ValueStatement interface {
	Node
	valueStatement()
}

type Expression interface {
	Node
	Statement
	ValueStatement
	expression()
}

// KeyValueStatement represents a key-value pair (e.g., apiVersion: "apps/v1")
type KeyValueStatement struct {
	Key   string
	Value ValueStatement
	Pos   Pos
}

func (kv *KeyValueStatement) node()       {}
func (kv *KeyValueStatement) statement()  {}
func (kv *KeyValueStatement) GetPos() Pos { return kv.Pos }

// StringLiteral represents a quoted string
type StringLiteral struct {
	Value string
	Pos   Pos
}

func (s *StringLiteral) node()           {}
func (s *StringLiteral) expression()     {}
func (s *StringLiteral) statement()      {}
func (s *StringLiteral) valueStatement() {}
func (s *StringLiteral) GetPos() Pos     { return s.Pos }

// InterpolatedString represents a string with embedded expressions (e.g., "Hello ${name}!")
type InterpolatedString struct {
	Parts []Expression // Alternating StringLiteral and expression values
	Pos   Pos
}

func (i *InterpolatedString) node()           {}
func (i *InterpolatedString) expression()     {}
func (i *InterpolatedString) statement()      {}
func (i *InterpolatedString) valueStatement() {}
func (i *InterpolatedString) GetPos() Pos     { return i.Pos }

// NumberLiteral represents a numeric value
type NumberLiteral struct {
	Value float64
	Pos   Pos
}

func (n *NumberLiteral) node()           {}
func (n *NumberLiteral) expression()     {}
func (n *NumberLiteral) statement()      {}
func (n *NumberLiteral) valueStatement() {}
func (n *NumberLiteral) GetPos() Pos     { return n.Pos }

// BooleanLiteral represents a boolean value (true or false)
type BooleanLiteral struct {
	Value bool
	Pos   Pos
}

func (b *BooleanLiteral) node()           {}
func (b *BooleanLiteral) expression()     {}
func (b *BooleanLiteral) statement()      {}
func (b *BooleanLiteral) valueStatement() {}
func (b *BooleanLiteral) GetPos() Pos     { return b.Pos }

// NullLiteral represents a null value
type NullLiteral struct {
	Pos Pos
}

func (n *NullLiteral) node()           {}
func (n *NullLiteral) expression()     {}
func (n *NullLiteral) statement()      {}
func (n *NullLiteral) valueStatement() {}
func (n *NullLiteral) GetPos() Pos     { return n.Pos }

// CurrentContext represents a reference to the current context (.)
type CurrentContext struct {
	Pos Pos
}

func (c *CurrentContext) node()           {}
func (c *CurrentContext) expression()     {}
func (c *CurrentContext) statement()      {}
func (c *CurrentContext) valueStatement() {}
func (c *CurrentContext) GetPos() Pos     { return c.Pos }

// Identifier represents a variable reference (e.g., Values, name)
type Identifier struct {
	Name string
	Pos  Pos
}

func (i *Identifier) node()           {}
func (i *Identifier) expression()     {}
func (i *Identifier) statement()      {}
func (i *Identifier) valueStatement() {}
func (i *Identifier) GetPos() Pos     { return i.Pos }

// MemberExpression represents member access (e.g., obj.key)
type MemberExpression struct {
	Object Expression
	Member string
	Pos    Pos
}

func (m *MemberExpression) node()           {}
func (m *MemberExpression) expression()     {}
func (m *MemberExpression) statement()      {}
func (m *MemberExpression) valueStatement() {}
func (m *MemberExpression) GetPos() Pos     { return m.Pos }

// IndexExpression represents array/object indexing (e.g., array[0], obj[key])
type IndexExpression struct {
	Object Expression
	Index  Expression
	Pos    Pos
}

func (idx *IndexExpression) node()           {}
func (idx *IndexExpression) expression()     {}
func (idx *IndexExpression) statement()      {}
func (idx *IndexExpression) valueStatement() {}
func (idx *IndexExpression) GetPos() Pos     { return idx.Pos }

// BinaryOp represents a binary operation (e.g., Values.debug && Values.verbose)
type BinaryOp struct {
	Left     Expression
	Operator string // "&&", "||", "==", "!=", "<", "<=", ">", ">=", "+", "-", "*", "/"
	Right    Expression
	Pos      Pos
}

func (b *BinaryOp) node()           {}
func (b *BinaryOp) expression()     {}
func (b *BinaryOp) statement()      {}
func (b *BinaryOp) valueStatement() {}
func (b *BinaryOp) GetPos() Pos     { return b.Pos }

// UnaryOp represents a unary operation (e.g., !Values.debug)
type UnaryOp struct {
	Operator string // "!"
	Operand  Expression
	Pos      Pos
}

func (u *UnaryOp) node()           {}
func (u *UnaryOp) expression()     {}
func (u *UnaryOp) statement()      {}
func (u *UnaryOp) valueStatement() {}
func (u *UnaryOp) GetPos() Pos     { return u.Pos }

// Object represents an object (e.g., {key: value})
type Object struct {
	Body []Node
	Pos  Pos
}

func (o *Object) node()           {}
func (o *Object) expression()     {}
func (o *Object) statement()      {}
func (o *Object) valueStatement() {}
func (o *Object) GetPos() Pos     { return o.Pos }

// Array represents an array (e.g., [value1, value2])
type Array struct {
	Body []Node
	Pos  Pos
}

func (a *Array) node()           {}
func (a *Array) expression()     {}
func (a *Array) statement()      {}
func (a *Array) valueStatement() {}
func (a *Array) GetPos() Pos     { return a.Pos }

// SpreadStatement represents a spread operator (e.g., spread obj, spread arr)
type SpreadStatement struct {
	Operand ValueStatement
	Pos     Pos
}

func (s *SpreadStatement) node()       {}
func (s *SpreadStatement) statement()  {}
func (s *SpreadStatement) GetPos() Pos { return s.Pos }

// IfStatement represents a conditional (e.g., if Values.debug { ... })
type IfStatement struct {
	Condition Expression
	Body      []Node
	Else      []Node // Optional else clause
	Pos       Pos
}

func (i *IfStatement) node()           {}
func (i *IfStatement) statement()      {}
func (i *IfStatement) valueStatement() {}
func (i *IfStatement) GetPos() Pos     { return i.Pos }

// ForStatement represents a loop (e.g., for k, v in Values.extraEnvs { ... })
type ForStatement struct {
	KeyVar   string
	ValueVar string
	Iterable Expression
	Body     []Node
	Pos      Pos
}

func (f *ForStatement) node()           {}
func (f *ForStatement) statement()      {}
func (f *ForStatement) valueStatement() {}
func (f *ForStatement) GetPos() Pos     { return f.Pos }

// WithStatement represents a context change (e.g., with Values.ingress as ing do ... end)
type WithStatement struct {
	Context Expression
	VarName string // Variable name for the context (optional, empty string means use ".")
	Body    []Node
	Pos     Pos
}

func (w *WithStatement) node()           {}
func (w *WithStatement) statement()      {}
func (w *WithStatement) valueStatement() {}
func (w *WithStatement) GetPos() Pos     { return w.Pos }

// BreakStatement represents a break statement in a loop
type BreakStatement struct {
	Pos Pos
}

func (b *BreakStatement) node()       {}
func (b *BreakStatement) statement()  {}
func (b *BreakStatement) GetPos() Pos { return b.Pos }

// ContinueStatement represents a continue statement in a loop
type ContinueStatement struct {
	Pos Pos
}

func (c *ContinueStatement) node()       {}
func (c *ContinueStatement) statement()  {}
func (c *ContinueStatement) GetPos() Pos { return c.Pos }

// LetStatement represents a variable definition (e.g., let name = "helmtk")
type LetStatement struct {
	Name  string
	Value ValueStatement
	Pos   Pos
}

func (l *LetStatement) node()       {}
func (l *LetStatement) statement()  {}
func (l *LetStatement) GetPos() Pos { return l.Pos }

// AssignmentStatement represents variable reassignment (e.g., name = "new value")
type AssignmentStatement struct {
	Name  string
	Value ValueStatement
	Pos   Pos
}

func (a *AssignmentStatement) node()       {}
func (a *AssignmentStatement) statement()  {}
func (a *AssignmentStatement) GetPos() Pos { return a.Pos }

// Comment represents a comment line
type Comment struct {
	Text string
	Pos  Pos
}

func (c *Comment) node()       {}
func (c *Comment) GetPos() Pos { return c.Pos }

// Definition represents a template definition (e.g., define(name, arg1, arg2) body)
type Definition struct {
	Name   string
	Body   []Node // Single value for expression form, multiple for do block
	Pos    Pos
}

func (d *Definition) node()       {}
func (d *Definition) GetPos() Pos { return d.Pos }

// IncludeExpression represents a template inclusion (e.g., include(name, arg1, arg2))
type IncludeExpression struct {
	Name string
	Context Expression
	Pos  Pos
}

func (i *IncludeExpression) node()       {}
func (i *IncludeExpression) expression() {}
func (i *IncludeExpression) statement()  {}
func (i *IncludeExpression) valueStatement()  {}
func (i *IncludeExpression) GetPos() Pos { return i.Pos }

// CallExpression represents a function call (e.g., upper(name), quote(str))
type CallExpression struct {
	Function Expression
	Args     []Expression
	Pos      Pos
}

func (c *CallExpression) node()           {}
func (c *CallExpression) expression()     {}
func (c *CallExpression) statement()      {}
func (c *CallExpression) valueStatement() {}
func (c *CallExpression) GetPos() Pos     { return c.Pos }
