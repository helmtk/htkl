package runtime

import (
	"testing"
)

func TestScopeVariables(t *testing.T) {
	scope := NewScope(nil)

	// Set a variable
	scope.Set("x", NewNumber(42))

	// Get the variable
	val, err := scope.Get("x")
	if err != nil {
		t.Fatalf("Get(x) error = %v", err)
	}
	num, ok := val.(*NumberValue)
	if !ok || num.Value != 42 {
		t.Errorf("Get(x) = %v, want 42", val)
	}

	// Get undefined variable
	_, err = scope.Get("undefined")
	if err == nil {
		t.Error("expected error for undefined variable")
	}
}

func TestScopeParent(t *testing.T) {
	parent := NewScope(nil)
	parent.Set("x", NewNumber(10))
	parent.Set("y", NewString("parent"))

	child := NewScope(parent)
	child.Set("y", NewString("child"))
	child.Set("z", NewBool(true))

	// Child can access parent variable
	val, err := child.Get("x")
	if err != nil {
		t.Fatalf("child.Get(x) error = %v", err)
	}
	if num, ok := val.(*NumberValue); !ok || num.Value != 10 {
		t.Errorf("child.Get(x) = %v, want 10", val)
	}

	// Child shadows parent variable
	val, err = child.Get("y")
	if err != nil {
		t.Fatalf("child.Get(y) error = %v", err)
	}
	if str, ok := val.(*StringValue); !ok || str.Value != "child" {
		t.Errorf("child.Get(y) = %v, want 'child'", val)
	}

	// Child has its own variable
	val, err = child.Get("z")
	if err != nil {
		t.Fatalf("child.Get(z) error = %v", err)
	}
	if b, ok := val.(*BoolValue); !ok || !b.Value {
		t.Errorf("child.Get(z) = %v, want true", val)
	}

	// Parent can't access child variable
	_, err = parent.Get("z")
	if err == nil {
		t.Error("expected error when parent tries to access child variable")
	}
}

func TestScopeTemplates(t *testing.T) {
	scope := NewScope(nil)

	// Define a template
	tmpl := NewTemplate("myTemplate", nil, "test.helmtk")
	scope.DefineTemplate("myTemplate", tmpl)

	// Get the template
	got, err := scope.GetTemplate("myTemplate")
	if err != nil {
		t.Fatalf("GetTemplate error = %v", err)
	}
	if got.Name != "myTemplate" {
		t.Errorf("template name = %s, want myTemplate", got.Name)
	}

	// Get undefined template
	_, err = scope.GetTemplate("undefined")
	if err == nil {
		t.Error("expected error for undefined template")
	}
}

func TestScopeTemplateInheritance(t *testing.T) {
	parent := NewScope(nil)
	tmpl1 := NewTemplate("parent", nil, "parent.helmtk")
	parent.DefineTemplate("parent", tmpl1)

	child := NewScope(parent)
	tmpl2 := NewTemplate("child", nil, "child.helmtk")
	child.DefineTemplate("child", tmpl2)

	// Child can access parent template
	got, err := child.GetTemplate("parent")
	if err != nil {
		t.Fatalf("child.GetTemplate(parent) error = %v", err)
	}
	if got.Name != "parent" {
		t.Errorf("template name = %s, want parent", got.Name)
	}

	// Child can access its own template
	got, err = child.GetTemplate("child")
	if err != nil {
		t.Fatalf("child.GetTemplate(child) error = %v", err)
	}
	if got.Name != "child" {
		t.Errorf("template name = %s, want child", got.Name)
	}

	// Parent can't access child template
	_, err = parent.GetTemplate("child")
	if err == nil {
		t.Error("expected error when parent tries to access child template")
	}
}
