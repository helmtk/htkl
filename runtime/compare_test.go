package runtime

import "testing"

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		left     Value
		right    Value
		expected bool
	}{
		{
			name:     "equal strings",
			left:     NewString("hello"),
			right:    NewString("hello"),
			expected: true,
		},
		{
			name:     "unequal strings",
			left:     NewString("hello"),
			right:    NewString("world"),
			expected: false,
		},
		{
			name:     "equal numbers",
			left:     NewNumber(42),
			right:    NewNumber(42),
			expected: true,
		},
		{
			name:     "unequal numbers",
			left:     NewNumber(42),
			right:    NewNumber(43),
			expected: false,
		},
		{
			name:     "equal bools",
			left:     NewBool(true),
			right:    NewBool(true),
			expected: true,
		},
		{
			name:     "unequal bools",
			left:     NewBool(true),
			right:    NewBool(false),
			expected: false,
		},
		{
			name:     "equal nulls",
			left:     NewNull(),
			right:    NewNull(),
			expected: true,
		},
		{
			name:     "different types",
			left:     NewString("42"),
			right:    NewNumber(42),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equal(tt.left, tt.right)
			if result != tt.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", tt.left, tt.right, result, tt.expected)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		name     string
		left     Value
		right    Value
		expected bool
	}{
		{
			name:     "equal strings",
			left:     NewString("hello"),
			right:    NewString("hello"),
			expected: false,
		},
		{
			name:     "unequal strings",
			left:     NewString("hello"),
			right:    NewString("world"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NotEqual(tt.left, tt.right)
			if result != tt.expected {
				t.Errorf("NotEqual(%v, %v) = %v, want %v", tt.left, tt.right, result, tt.expected)
			}
		})
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		name      string
		left      Value
		right     Value
		expected  bool
		shouldErr bool
	}{
		{
			name:     "less than",
			left:     NewNumber(1),
			right:    NewNumber(2),
			expected: true,
		},
		{
			name:     "equal",
			left:     NewNumber(2),
			right:    NewNumber(2),
			expected: false,
		},
		{
			name:     "greater than",
			left:     NewNumber(3),
			right:    NewNumber(2),
			expected: false,
		},
		{
			name:      "non-numeric left",
			left:      NewString("hello"),
			right:     NewNumber(2),
			shouldErr: true,
		},
		{
			name:      "non-numeric right",
			left:      NewNumber(1),
			right:     NewString("world"),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Less(tt.left, tt.right)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Less(%v, %v) should return error", tt.left, tt.right)
				}
				return
			}
			if err != nil {
				t.Errorf("Less(%v, %v) unexpected error: %v", tt.left, tt.right, err)
				return
			}
			if result != tt.expected {
				t.Errorf("Less(%v, %v) = %v, want %v", tt.left, tt.right, result, tt.expected)
			}
		})
	}
}

func TestLessEqual(t *testing.T) {
	tests := []struct {
		name     string
		left     Value
		right    Value
		expected bool
	}{
		{
			name:     "less than",
			left:     NewNumber(1),
			right:    NewNumber(2),
			expected: true,
		},
		{
			name:     "equal",
			left:     NewNumber(2),
			right:    NewNumber(2),
			expected: true,
		},
		{
			name:     "greater than",
			left:     NewNumber(3),
			right:    NewNumber(2),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LessEqual(tt.left, tt.right)
			if err != nil {
				t.Errorf("LessEqual(%v, %v) unexpected error: %v", tt.left, tt.right, err)
				return
			}
			if result != tt.expected {
				t.Errorf("LessEqual(%v, %v) = %v, want %v", tt.left, tt.right, result, tt.expected)
			}
		})
	}
}

func TestGreater(t *testing.T) {
	tests := []struct {
		name     string
		left     Value
		right    Value
		expected bool
	}{
		{
			name:     "less than",
			left:     NewNumber(1),
			right:    NewNumber(2),
			expected: false,
		},
		{
			name:     "equal",
			left:     NewNumber(2),
			right:    NewNumber(2),
			expected: false,
		},
		{
			name:     "greater than",
			left:     NewNumber(3),
			right:    NewNumber(2),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Greater(tt.left, tt.right)
			if err != nil {
				t.Errorf("Greater(%v, %v) unexpected error: %v", tt.left, tt.right, err)
				return
			}
			if result != tt.expected {
				t.Errorf("Greater(%v, %v) = %v, want %v", tt.left, tt.right, result, tt.expected)
			}
		})
	}
}

func TestGreaterEqual(t *testing.T) {
	tests := []struct {
		name     string
		left     Value
		right    Value
		expected bool
	}{
		{
			name:     "less than",
			left:     NewNumber(1),
			right:    NewNumber(2),
			expected: false,
		},
		{
			name:     "equal",
			left:     NewNumber(2),
			right:    NewNumber(2),
			expected: true,
		},
		{
			name:     "greater than",
			left:     NewNumber(3),
			right:    NewNumber(2),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GreaterEqual(tt.left, tt.right)
			if err != nil {
				t.Errorf("GreaterEqual(%v, %v) unexpected error: %v", tt.left, tt.right, err)
				return
			}
			if result != tt.expected {
				t.Errorf("GreaterEqual(%v, %v) = %v, want %v", tt.left, tt.right, result, tt.expected)
			}
		})
	}
}
