package parser

import (
	"testing"
)

func TestQuotedKeys(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "Kubernetes labels",
			input: `labels: {
  "app.kubernetes.io/name": "myapp"
  "app.kubernetes.io/instance": "prod"
  "app.kubernetes.io/version": "1.0.0"
}`,
		},
		{
			name: "Mixed quoted and unquoted keys",
			input: `config: {
  name: "myapp"
  "special-key": "value"
  "key.with.dots": "value2"
}`,
		},
		{
			name: "Nested objects with quoted keys",
			input: `metadata: {
  labels: {
    "app.kubernetes.io/name": "myapp"
  }
  annotations: {
    "prometheus.io/scrape": "true"
    "prometheus.io/port": "9090"
  }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := New(tt.input, "").Parse()
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			if doc == nil {
				t.Fatal("expected document, got nil")
			}

			if len(doc.Body) == 0 {
				t.Fatal("expected at least one statement")
			}

			t.Logf("Successfully parsed: %s", tt.name)
		})
	}
}
