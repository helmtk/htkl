package parser

import (
	"strings"
	"testing"
)

func TestParseErrorFormatting(t *testing.T) {
	input := `config: {
  name: "myapp"
  version: "1.0"
  replicas 3
  ports: [80, 443]
}`

	_, err := New(input, "").Parse()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}

	// Check that formatted error includes context
	formatted := parseErr.FormatWithContext()
	if !strings.Contains(formatted, "Parse error at line") {
		t.Error("formatted error should include 'Parse error at line'")
	}
	if !strings.Contains(formatted, "replicas") {
		t.Error("formatted error should include the problematic line")
	}
	if !strings.Contains(formatted, "^") {
		t.Error("formatted error should include a pointer (^) to the error location")
	}

	t.Logf("Formatted error:\n%s", formatted)
}

func TestWithStatementRequiresAs(t *testing.T) {
	input := `with Values.routes do
  name: "test"
end`

	_, err := New(input, "test.helmtk").Parse()
	if err == nil {
		t.Fatal("expected parse error for 'with' without 'as', got nil")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}

	// Check that error mentions 'as'
	errMsg := parseErr.Error()
	if !strings.Contains(errMsg, "'as'") {
		t.Errorf("error should mention 'as', got: %s", errMsg)
	}

	t.Logf("Parse error (as expected):\n%s", parseErr.FormatWithContext())
}
