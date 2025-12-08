package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseError represents a parsing error with position information
type ParseError struct {
	Message string
	Line    int
	Col     int
	Offset  int
	Source  string // The full source code for context
}

func (e *ParseError) Error() string {
	return e.FormatWithContext()
}

// FormatWithContext returns a formatted error message with source context
func (e *ParseError) FormatWithContext() string {
	var sb strings.Builder

	// Write the basic error message
	sb.WriteString(fmt.Sprintf("Parse error at line %d, column %d: %s\n", e.Line, e.Col, e.Message))

	// Add source context (3 lines before, error line, 3 lines after)
	if e.Source != "" {
		lines := strings.Split(e.Source, "\n")
		if e.Line > 0 && e.Line <= len(lines) {
			sb.WriteString("\n")

			// Show 3 lines before (if they exist)
			contextBefore := 3
			for i := contextBefore; i >= 1; i-- {
				lineNum := e.Line - i
				if lineNum > 0 {
					sb.WriteString(fmt.Sprintf("%4d | %s\n", lineNum, lines[lineNum-1]))
				}
			}

			// Error line with pointer
			errorLine := lines[e.Line-1]
			sb.WriteString(fmt.Sprintf("%4d | %s\n", e.Line, errorLine))

			// Pointer to error column
			pointer := strings.Repeat(" ", 7+e.Col-1) + "^"
			sb.WriteString(pointer + "\n")

			// Show 3 lines after (if they exist)
			contextAfter := 3
			for i := 1; i <= contextAfter; i++ {
				lineNum := e.Line + i
				if lineNum <= len(lines) {
					sb.WriteString(fmt.Sprintf("%4d | %s\n", lineNum, lines[lineNum-1]))
				}
			}
		}
	}

	return sb.String()
}

// Parser represents a helmtk template parser
type Parser struct {
	lexer    *Lexer
	current  Token
	peek     Token
	source   string // Store source for error reporting
	filename string // Source filename for position tracking
}

func New(source, filename string) *Parser {
	p := &Parser{
		lexer:    NewLexer(source),
		source:   source,
		filename: filename,
	}
	// Initialize current and peek tokens
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) skipNewlines() {
	for p.currentIs(TokenNewline) {
		p.nextToken()
	}
}

func (p *Parser) currentIs(t TokenType) bool {
	return p.current.Type == t
}

func (p *Parser) peekIs(t TokenType) bool {
	return p.peek.Type == t
}

func (p *Parser) pos() Pos {
	return Pos{
		Filename: p.filename,
		Line:     p.current.Line,
		Col:      p.current.Col,
	}
}

func (p *Parser) expectCurrent(t TokenType) error {
	if !p.currentIs(t) {
		return p.error(fmt.Sprintf("expected %v, got %v", t, p.current.Type))
	}
	return nil
}

// error creates a ParseError at the current token position
func (p *Parser) error(message string) *ParseError {
	return &ParseError{
		Message: message,
		Line:    p.current.Line,
		Col:     p.current.Col,
		Offset:  p.lexer.pos,
		Source:  p.source,
	}
}

// Operator precedence levels (higher = tighter binding)
const (
	PREC_LOWEST     = iota
	PREC_PIPE       // |
	PREC_OR         // ||
	PREC_AND        // &&
	PREC_EQUALS     // ==, !=
	PREC_COMPARISON // <, <=, >, >=
	PREC_SUM        // +, -
	PREC_PRODUCT    // *, /
)

func (p *Parser) tokenPrecedence(t TokenType) int {
	switch t {
	case TokenPipe:
		return PREC_PIPE
	case TokenOr:
		return PREC_OR
	case TokenAnd:
		return PREC_AND
	case TokenEq, TokenNeq:
		return PREC_EQUALS
	case TokenLt, TokenLte, TokenGt, TokenGte:
		return PREC_COMPARISON
	case TokenPlus, TokenMinus:
		return PREC_SUM
	case TokenMul, TokenDiv:
		return PREC_PRODUCT
	default:
		return PREC_LOWEST
	}
}

func (p *Parser) peekPrecedence() int {
	return p.tokenPrecedence(p.peek.Type)
}

// Parse parses the input and returns a Document AST node
func (p *Parser) Parse() (*Document, error) {
	doc := &Document{}

	for !p.currentIs(TokenEOF) {
		// Skip newlines and comments
		if p.currentIs(TokenNewline) || p.currentIs(TokenComment) {
			p.nextToken()
			continue
		}

		if p.currentIs(TokenDefine) {
			d, err := p.parseDefinition()
			if err != nil {
				return nil, err
			}
			doc.Definitions = append(doc.Definitions, d)
			p.skipNewlines()
			continue
		}

		node, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if node != nil {
			doc.Body = append(doc.Body, node)
		}

		p.nextToken()
		p.skipNewlines()
	}

	return doc, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	switch p.current.Type {
	case TokenFor:
		return p.parseForStatement()
	case TokenWith:
		return p.parseWithStatement()
	case TokenBreak:
		return &BreakStatement{Pos: p.pos()}, nil
	case TokenContinue:
		return &ContinueStatement{Pos: p.pos()}, nil
	case TokenLet:
		return p.parseLetStatement()
	case TokenSpread:
		return p.parseSpread()
	case TokenIf:
		return p.parseIfStatement()
	case TokenComment:
		// TODO
		// return &Comment{Text: p.current.Value}, nil
		return nil, nil
	case TokenString:
		// Check if this is a key-value pair (for use in objects/conditionals)
		if p.peekIs(TokenColon) {
			return p.parseKeyValue()
		}
	case TokenIdent:
		// Check if this is an assignment statement
		if p.peekIs(TokenAssign) {
			return p.parseAssignmentStatement()
		}
		// Check if this is a key-value pair (for use in objects/conditionals)
		if p.peekIs(TokenColon) {
			return p.parseKeyValue()
		}
	case TokenEOF, TokenEnd:
		return nil, nil
	}

	return p.parseExpression()
}

func (p *Parser) parseSpread() (*SpreadStatement, error) {
	pos := p.pos()
	p.nextToken() // skip 'spread'
	operand, err := p.parseValueStatement()
	if err != nil {
		return nil, err
	}
	return &SpreadStatement{Operand: operand, Pos: pos}, nil
}

func (p *Parser) parseAssignmentStatement() (*AssignmentStatement, error) {
	pos := p.pos()

	if err := p.expectCurrent(TokenIdent); err != nil {
		return nil, err
	}

	name := p.current.Value
	p.nextToken()

	// Expect '='
	if err := p.expectCurrent(TokenAssign); err != nil {
		return nil, err
	}
	p.nextToken()

	value, err := p.parseValueStatement()
	if err != nil {
		return nil, err
	}

	return &AssignmentStatement{Name: name, Value: value, Pos: pos}, nil
}

func (p *Parser) parseLetStatement() (*LetStatement, error) {
	pos := p.pos()
	p.nextToken() // skip 'let'

	if err := p.expectCurrent(TokenIdent); err != nil {
		return nil, err
	}

	name := p.current.Value
	p.nextToken()

	// Expect '='
	if err := p.expectCurrent(TokenAssign); err != nil {
		return nil, err
	}
	p.nextToken()

	value, err := p.parseValueStatement()
	if err != nil {
		return nil, err
	}

	return &LetStatement{
		Name:  name,
		Value: value,
		Pos:   pos,
	}, nil
}

func (p *Parser) parseDefinition() (*Definition, error) {
	pos := p.pos()
	p.nextToken() // skip 'define'

	// Expect '('
	if err := p.expectCurrent(TokenLParen); err != nil {
		return nil, err
	}
	p.nextToken()

	// Parse name (first parameter - must be a string)
	if err := p.expectCurrent(TokenString); err != nil {
		return nil, err
	}
	name := p.current.Value
	p.nextToken()

	// Expect ')'
	if err := p.expectCurrent(TokenRParen); err != nil {
		return nil, err
	}
	p.nextToken()

	// Parse body
	var body []Node

	// Check for 'do' keyword (block form)
	if p.currentIs(TokenDo) {
		p.nextToken() // skip 'do'
		p.skipNewlines()

		// Parse block body (multiple statements)
		for !p.currentIs(TokenEnd) && !p.currentIs(TokenEOF) {
			// Skip comments and newlines
			if p.currentIs(TokenComment) || p.currentIs(TokenNewline) {
				p.nextToken()
				continue
			}

			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}

			body = append(body, stmt)
			p.nextToken()
			p.skipNewlines()

			// Optional comma
			if p.currentIs(TokenComma) {
				p.nextToken()
				p.skipNewlines()
			}
		}

		if err := p.expectCurrent(TokenEnd); err != nil {
			return nil, err
		}
		p.nextToken() // skip 'end'
	} else {
		// Expression form (single value)
		value, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		body = []Node{value}
		p.nextToken() // advance past expression
	}

	return &Definition{
		Name:   name,
		Body:   body,
		Pos:    pos,
	}, nil
}

func (p *Parser) parseValueStatement() (ValueStatement, error) {
	switch p.current.Type {
	case TokenFor:
		return p.parseForStatement()
	case TokenWith:
		return p.parseWithStatement()
	case TokenIf:
		return p.parseIfStatement()
	case TokenComment:
		// TODO
		// return &Comment{Text: p.current.Value}, nil
		return nil, nil
	case TokenEOF, TokenEnd:
		return nil, nil
	}

	return p.parseExpression()
}

func (p *Parser) parseKeyValue() (*KeyValueStatement, error) {
	pos := p.pos()
	key := p.current.Value
	p.nextToken()

	if err := p.expectCurrent(TokenColon); err != nil {
		return nil, err
	}
	p.nextToken()

	value, err := p.parseValueStatement()
	if err != nil {
		return nil, err
	}

	// After a key-value expression, we should see a statement terminator
	if err := p.expectStatementEnd(); err != nil {
		return nil, err
	}

	return &KeyValueStatement{
		Key:   key,
		Value: value,
		Pos:   pos,
	}, nil
}

// expectStatementEnd checks that the current position is a valid statement terminator
func (p *Parser) expectStatementEnd() error {
	// Valid terminators: newline, comma, closing brace/bracket, end, else, EOF
	switch p.peek.Type {
	case TokenNewline, TokenComma, TokenRBrace, TokenRBracket, TokenEnd, TokenElse, TokenEOF:
		return nil
	default:
		return p.error(fmt.Sprintf("unexpected token %v after expression", p.peek.Type))
	}
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseValueWithPrecedence(PREC_LOWEST)
}

func (p *Parser) parseValueWithPrecedence(minPrecedence int) (Expression, error) {
	left, err := p.parsePostfixValue()
	if err != nil {
		return nil, err
	}

	// Precedence climbing: handle binary operators
	for p.peekPrecedence() > minPrecedence {
		p.nextToken() // move to operator
		pos := p.pos()
		operator := p.current.Value
		precedence := p.tokenPrecedence(p.current.Type)

		p.nextToken() // move to right operand

		// Parse right operand with higher precedence for left-associativity
		right, err := p.parseValueWithPrecedence(precedence)
		if err != nil {
			return nil, err
		}

		left = &BinaryOp{
			Left:     left,
			Operator: operator,
			Right:    right,
			Pos:      pos,
		}
	}

	return left, nil
}

func (p *Parser) parsePostfixValue() (Expression, error) {
	value, err := p.parsePrimaryValue()
	if err != nil {
		return nil, err
	}

	// Handle postfix operators: ., [, and (
	for {
		if p.peekIs(TokenDot) {
			p.nextToken() // move to .
			pos := p.pos()
			p.nextToken() // move to member name
			if err := p.expectCurrent(TokenIdent); err != nil {
				return nil, err
			}
			value = &MemberExpression{
				Object: value,
				Member: p.current.Value,
				Pos:    pos,
			}
		} else if p.peekIs(TokenLBracket) {
			p.nextToken() // move to [
			pos := p.pos()
			p.nextToken() // move to index expression
			index, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			p.nextToken() // move to ]
			if err := p.expectCurrent(TokenRBracket); err != nil {
				return nil, err
			}
			value = &IndexExpression{
				Object: value,
				Index:  index,
				Pos:    pos,
			}
		} else if p.peekIs(TokenLParen) {
			p.nextToken() // move to (
			pos := p.pos()
			p.nextToken() // move to first arg or )

			args := []Expression{}
			for !p.currentIs(TokenRParen) && !p.currentIs(TokenEOF) {
				arg, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				p.nextToken()

				// Optional comma
				if p.currentIs(TokenComma) {
					p.nextToken()
				}
			}

			if err := p.expectCurrent(TokenRParen); err != nil {
				return nil, err
			}

			value = &CallExpression{
				Function: value,
				Args:     args,
				Pos:      pos,
			}
		} else {
			break
		}
	}

	return value, nil
}

func (p *Parser) parsePrimaryValue() (Expression, error) {
	pos := p.pos()

	switch p.current.Type {

	case TokenString:
		return p.parseStringLiteral()
	case TokenTrue:
		return &BooleanLiteral{Value: true, Pos: pos}, nil
	case TokenFalse:
		return &BooleanLiteral{Value: false, Pos: pos}, nil
	case TokenNull:
		return &NullLiteral{Pos: pos}, nil
	case TokenDot:
		// Check if this is .identifier (member access on current context)
		// or just . (current context reference)
		if p.peekIs(TokenIdent) {
			// This is .identifier - create a MemberExpression
			p.nextToken() // move to identifier
			return &MemberExpression{
				Object: &CurrentContext{Pos: pos},
				Member: p.current.Value,
				Pos:    pos,
			}, nil
		}
		// Just . by itself
		return &CurrentContext{Pos: pos}, nil
	case TokenLBrace:
		return p.parseObject()
	case TokenLBracket:
		return p.parseArray()
	case TokenInclude:
		return p.parseIncludeExpression()

	case TokenNumber:
		num, err := strconv.ParseFloat(p.current.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", p.current.Value, err)
		}
		return &NumberLiteral{Value: num, Pos: pos}, nil

	case TokenIdent:
		return p.parseIdentifier()

	case TokenNot:
		p.nextToken() // move past !
		operand, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return &UnaryOp{
			Operator: "!",
			Operand:  operand,
			Pos:      pos,
		}, nil

	case TokenLParen:
		p.nextToken() // move past (
		value, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.nextToken() // move to )
		if err := p.expectCurrent(TokenRParen); err != nil {
			return nil, err
		}
		return value, nil

	default:
		return nil, p.error(fmt.Sprintf("unexpected token %v", p.current.Type))
	}
}

func (p *Parser) parseIdentifier() (*Identifier, error) {
	return &Identifier{Name: p.current.Value, Pos: p.pos()}, nil
}

func (p *Parser) parseIncludeExpression() (*IncludeExpression, error) {
	pos := p.pos()
	p.nextToken() // skip 'include'

	// Expect '('
	if err := p.expectCurrent(TokenLParen); err != nil {
		return nil, err
	}
	p.nextToken()

	// Parse name (first argument - must be a string)
	if err := p.expectCurrent(TokenString); err != nil {
		return nil, err
	}
	name := p.current.Value
	p.nextToken()

	// Parse arguments
	var context Expression
	if p.currentIs(TokenComma) {
		p.nextToken() // skip comma

		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		context = arg
		p.nextToken()
	}

	// Expect ')'
	if err := p.expectCurrent(TokenRParen); err != nil {
		return nil, err
	}

	return &IncludeExpression{
		Name: name,
		Context: context,
		Pos:  pos,
	}, nil
}

func (p *Parser) parseStringLiteral() (Expression, error) {
	stringPos := p.pos()
	str := p.current.Value

	// Check if string contains interpolation (but not escaped \x00${)
	hasInterpolation := false
	for i := 0; i < len(str)-1; i++ {
		if str[i] == '$' && str[i+1] == '{' {
			// Check if this ${ is escaped (preceded by \x00)
			if i > 0 && str[i-1] == '\x00' {
				continue
			}
			hasInterpolation = true
			break
		}
	}

	if !hasInterpolation {
		// No interpolation, just clean up escaped $ markers
		cleanStr := strings.ReplaceAll(str, "\x00$", "$")
		return &StringLiteral{Value: cleanStr, Pos: stringPos}, nil
	}

	// Parse interpolated string
	parts := []Expression{}
	pos := 0

	for {
		// Find next interpolation (skip \x00$ which is escaped)
		start := -1
		searchPos := pos
		for searchPos < len(str) {
			idx := strings.Index(str[searchPos:], "${")
			if idx == -1 {
				break
			}
			// Check if this ${ is escaped (preceded by \x00)
			absoluteIdx := searchPos + idx
			if absoluteIdx > 0 && str[absoluteIdx-1] == '\x00' {
				// This is an escaped ${, skip it
				searchPos = absoluteIdx + 2
				continue
			}
			start = idx
			pos = searchPos
			break
		}

		if start == -1 {
			// No more interpolations, add remaining string if non-empty
			if pos < len(str) {
				cleanStr := strings.ReplaceAll(str[pos:], "\x00$", "$")
				parts = append(parts, &StringLiteral{Value: cleanStr, Pos: stringPos})
			}
			break
		}

		// Add string before interpolation if non-empty
		if start > 0 {
			cleanStr := strings.ReplaceAll(str[pos:pos+start], "\x00$", "$")
			parts = append(parts, &StringLiteral{Value: cleanStr, Pos: stringPos})
		}

		// Find end of interpolation
		pos += start + 2 // skip "${"
		end := strings.Index(str[pos:], "}")
		if end == -1 {
			return nil, p.error("unclosed interpolation in string")
		}

		// Parse the expression inside ${}
		exprStr := str[pos : pos+end]
		exprParser := &Parser{
			lexer:    NewLexer(exprStr),
			filename: p.filename,
		}
		exprParser.nextToken()
		exprParser.nextToken()

		expr, err := exprParser.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse interpolation expression %q: %w", exprStr, err)
		}

		parts = append(parts, expr)
		pos += end + 1 // skip past "}"
	}

	return &InterpolatedString{Parts: parts, Pos: stringPos}, nil
}

func (p *Parser) parseObject() (*Object, error) {
	obj := &Object{Pos: p.pos()}

	p.nextToken() // skip '{'
	p.skipNewlines()

	for !p.currentIs(TokenRBrace) && !p.currentIs(TokenEOF) {
		// Skip comments and newlines
		if p.currentIs(TokenComment) || p.currentIs(TokenNewline) {
			p.nextToken()
			continue
		}

		// In object context, identifiers must be followed by colon (key-value pairs)
		if p.currentIs(TokenIdent) && !p.peekIs(TokenColon) {
			return nil, p.error(fmt.Sprintf("expected ':', got %v", p.peek.Type))
		}
		if p.currentIs(TokenString) && !p.peekIs(TokenColon) {
			return nil, p.error(fmt.Sprintf("expected ':', got %v", p.peek.Type))
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		p.nextToken()
		p.skipNewlines()

		obj.Body = append(obj.Body, stmt)

		// Optional comma
		if p.currentIs(TokenComma) {
			p.nextToken()
			p.skipNewlines()
		}
	}

	if !p.currentIs(TokenRBrace) {
		return nil, p.error(fmt.Sprintf("expected '}', got %v", p.current.Type))
	}

	return obj, nil
}

func (p *Parser) parseArray() (*Array, error) {
	arr := &Array{Pos: p.pos()}

	p.nextToken() // skip '['
	p.skipNewlines()

	for !p.currentIs(TokenRBracket) && !p.currentIs(TokenEOF) {
		// Skip comments and newlines
		if p.currentIs(TokenComment) || p.currentIs(TokenNewline) {
			p.nextToken()
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		arr.Body = append(arr.Body, stmt)
		p.nextToken()
		p.skipNewlines()

		// Optional comma
		if p.currentIs(TokenComma) {
			p.nextToken()
			p.skipNewlines()
		}
	}

	if !p.currentIs(TokenRBracket) {
		return nil, p.error(fmt.Sprintf("expected ']', got %v", p.current.Type))
	}

	return arr, nil
}

func (p *Parser) parseIfStatement() (*IfStatement, error) {
	pos := p.pos()
	p.nextToken() // skip 'if'

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.nextToken()

	if err := p.expectCurrent(TokenDo); err != nil {
		return nil, err
	}

	p.nextToken() // skip 'do'
	p.skipNewlines()

	body := []Node{}
	for !p.currentIs(TokenElse) && !p.currentIs(TokenEnd) && !p.currentIs(TokenEOF) {
		// Skip comments and newlines
		if p.currentIs(TokenComment) || p.currentIs(TokenNewline) {
			p.nextToken()
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		body = append(body, stmt)
		p.nextToken()
		p.skipNewlines()

		// Optional comma
		if p.currentIs(TokenComma) {
			p.nextToken()
			p.skipNewlines()
		}
	}

	// Check for optional else clause
	elseBody := []Node{}
	if p.currentIs(TokenElse) {
		p.nextToken() // skip 'else'
		p.skipNewlines()

		// Check for else if
		if p.currentIs(TokenIf) {
			// Parse the if statement directly (which will consume its own 'end')
			ifStmt, err := p.parseIfStatement()
			if err != nil {
				return nil, err
			}
			elseBody = append(elseBody, ifStmt)
			// The nested if already consumed the 'end', so we're done
			return &IfStatement{
				Condition: condition,
				Body:      body,
				Else:      elseBody,
				Pos:       pos,
			}, nil
		} else {
			// Parse else block (no 'do' needed)
			for !p.currentIs(TokenEnd) && !p.currentIs(TokenEOF) {
				// Skip comments and newlines
				if p.currentIs(TokenComment) || p.currentIs(TokenNewline) {
					p.nextToken()
					continue
				}

				stmt, err := p.parseStatement()
				if err != nil {
					return nil, err
				}

				elseBody = append(elseBody, stmt)
				p.nextToken()
				p.skipNewlines()

				// Optional comma
				if p.currentIs(TokenComma) {
					p.nextToken()
					p.skipNewlines()
				}
			}
		}
	}

	if !p.currentIs(TokenEnd) {
		return nil, fmt.Errorf("expected 'end', got %v", p.current.Type)
	}

	return &IfStatement{
		Condition: condition,
		Body:      body,
		Else:      elseBody,
		Pos:       pos,
	}, nil
}

func (p *Parser) parseWithStatement() (*WithStatement, error) {
	pos := p.pos()
	p.nextToken() // skip 'with'

	context, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.nextToken()

	// Parse required "as <name>"
	if err := p.expectCurrent(TokenAs); err != nil {
		return nil, err
	}
	p.nextToken() // skip 'as'

	if err := p.expectCurrent(TokenIdent); err != nil {
		return nil, err
	}

	varName := p.current.Value
	p.nextToken()

	if err := p.expectCurrent(TokenDo); err != nil {
		return nil, err
	}

	p.nextToken() // skip 'do'
	p.skipNewlines()

	body := []Node{}
	for !p.currentIs(TokenEnd) && !p.currentIs(TokenEOF) {
		// Skip comments and newlines
		if p.currentIs(TokenComment) || p.currentIs(TokenNewline) {
			p.nextToken()
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		body = append(body, stmt)
		p.nextToken()
		p.skipNewlines()

		// Optional comma
		if p.currentIs(TokenComma) {
			p.nextToken()
			p.skipNewlines()
		}
	}

	if !p.currentIs(TokenEnd) {
		return nil, fmt.Errorf("expected 'end', got %v", p.current.Type)
	}

	return &WithStatement{
		Context: context,
		VarName: varName,
		Body:    body,
		Pos:     pos,
	}, nil
}

func (p *Parser) parseForStatement() (*ForStatement, error) {
	pos := p.pos()
	p.nextToken() // skip 'for'

	if err := p.expectCurrent(TokenIdent); err != nil {
		return nil, err
	}

	keyVar := p.current.Value
	p.nextToken()

	if err := p.expectCurrent(TokenComma); err != nil {
		return nil, err
	}
	p.nextToken()

	if err := p.expectCurrent(TokenIdent); err != nil {
		return nil, err
	}

	valueVar := p.current.Value
	p.nextToken()

	if err := p.expectCurrent(TokenIn); err != nil {
		return nil, err
	}
	p.nextToken()

	iterable, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.nextToken()

	if err := p.expectCurrent(TokenDo); err != nil {
		return nil, err
	}

	p.nextToken() // skip 'do'
	p.skipNewlines()

	body := []Node{}
	for !p.currentIs(TokenEnd) && !p.currentIs(TokenEOF) {
		// Skip comments and newlines
		if p.currentIs(TokenComment) || p.currentIs(TokenNewline) {
			p.nextToken()
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		body = append(body, stmt)
		p.nextToken()
		p.skipNewlines()

		// Optional comma
		if p.currentIs(TokenComma) {
			p.nextToken()
			p.skipNewlines()
		}
	}

	if !p.currentIs(TokenEnd) {
		return nil, fmt.Errorf("expected 'end', got %v", p.current.Type)
	}

	return &ForStatement{
		KeyVar:   keyVar,
		ValueVar: valueVar,
		Iterable: iterable,
		Body:     body,
		Pos:      pos,
	}, nil
}
