package parser

import (
	"testing"
)

func TestConditionalFieldsInObjects(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "Simple conditional field",
			input: `config: {
  name: "myapp"

  if Values.debug do
    logLevel: "debug"
  end
}`,
		},
		{
			name: "Conditional field with else",
			input: `config: {
  name: "myapp"

  if Values.debug do
    logLevel: "debug"
  else
    logLevel: "info"
  end
}`,
		},
		{
			name: "Kubernetes metadata example",
			input: `metadata: {
  name: include("cert-manager.fullname")
  namespace: include("cert-manager.namespace")
  labels: {
    app: include("cert-manager.name")
    "app.kubernetes.io/name": include("cert-manager.name")
    "app.kubernetes.io/instance": Release.Name
    "app.kubernetes.io/component": "controller"

    spread include("labels")
  }

  if Values.deploymentAnnotations do
    annotations: Values.deploymentAnnotations
  end
}`,
		},
		{
			name: "Multiple conditionals",
			input: `config: {
  name: "myapp"

  if Values.debug do
    logLevel: "debug"
  end

  if Values.metrics do
    metricsPort: 9090
  end
}`,
		},
		{
			name: "For loop in object",
			input: `env: {
  for k, v in Values.extraEnv do
    k: v
  end
}`,
		},
		{
			name: "Mixed fields, conditionals, and spread",
			input: `metadata: {
  name: "myapp"

  spread defaults

  if Values.labels do
    labels: Values.labels
  end

  version: "1.0"
}`,
		},
		{
			name: "Spread inside conditional",
			input: `template: {
  metadata: {
    labels: {
      app: include("cert-manager.name")
      "app.kubernetes.io/name": include("cert-manager.name")
      "app.kubernetes.io/instance": Release.Name
      "app.kubernetes.io/component": "controller"

      spread include("labels")

      if Values.podLabels do
        spread Values.podLabels
      end
    }
  }
}`,
		},
		{
			name: "With statement in object",
			input: `metadata: {
  name: include("cert-manager.fullname")

  with Values.deploymentAnnotations as a do
    annotations: a
  end
}`,
		},
		{
			name: "Let statement in object",
			input: `spec: {
  let finalNodeSelector = Values.nodeSelector | default(Values.global.nodeSelector)

  if finalNodeSelector do
    nodeSelector: finalNodeSelector
  end
}`,
		},
		{
			name: "Multiple with statements",
			input: `spec: {
  with Values.resources as r do
    resources: r
  end

  with Values.livenessProbe as p do
    livenessProbe: p
  end

  with Values.podDnsPolicy as d do
    dnsPolicy: d
  end
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
