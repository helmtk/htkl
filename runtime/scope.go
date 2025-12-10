package runtime

import (
	"fmt"

	"helmtk.dev/code/htkl/parser"
)

// Scope manages variable bindings and template definitions
type Scope struct {
	parent    *Scope
	vars      map[string]Value
	globals   map[string]Value
	funcs     map[string]Func
	templates map[string]*Template
}

// NewScope creates a new scope with an optional parent
func NewScope(parent *Scope) *Scope {
	s := &Scope{
		parent:    parent,
		vars:      make(map[string]Value),
		globals:   make(map[string]Value),
		funcs:     make(map[string]Func),
		templates: make(map[string]*Template),
	}
	if parent != nil {
		s.Link(parent)
	}
	return s
}

func (s *Scope) GetFunction(name string) (Func, bool) {
	f, ok := s.funcs[name]
	return f, ok
}

func (s *Scope) SetFunction(name string, f Func) {
	s.funcs[name] = f
}

// Get retrieves a variable value from this scope or parent scopes
func (s *Scope) Get(name string) (Value, error) {
	// Check this scope
	if val, ok := s.vars[name]; ok {
		return val, nil
	}

	// Check parent scope
	if s.parent != nil {
		return s.parent.Get(name)
	}

	// Check this scope
	if val, ok := s.globals[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("undefined variable: %s", name)
}

// Set binds a variable to a value in the current scope
func (s *Scope) Set(name string, val Value) {
	s.vars[name] = val
}

func (s *Scope) SetGlobal(name string, val Value) {
	s.globals[name] = val
}

// DefineTemplate registers a template in the current scope
func (s *Scope) DefineTemplate(name string, tmpl *Template) {
	s.templates[name] = tmpl
}

func (s *Scope) Link(other *Scope) {
	s.templates = other.templates
	s.globals = other.globals
	s.funcs = other.funcs
}

// GetTemplate retrieves a template from this scope or parent scopes
func (s *Scope) GetTemplate(name string) (*Template, error) {
	// Check this scope
	if tmpl, ok := s.templates[name]; ok {
		return tmpl, nil
	}

	// Check parent scope
	if s.parent != nil {
		return s.parent.GetTemplate(name)
	}

	return nil, fmt.Errorf("undefined template: %s", name)
}

// Template represents a user-defined template
type Template struct {
	Name     string
	Body     []parser.Node // The AST nodes to evaluate
	Filename string        // Source file where template was defined
}

// NewTemplate creates a new template with source file information
func NewTemplate(name string, body []parser.Node, filename string) *Template {
	return &Template{
		Name:     name,
		Body:     body,
		Filename: filename,
	}
}

type Func func(args ...Value) (Value, error)
