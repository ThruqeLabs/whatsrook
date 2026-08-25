package fun

import (
	"testing"
)

func TestEvalMathExpr(t *testing.T) {
	cases := []struct {
		expr     string
		expected float64
		hasErr   bool
	}{
		{"2 + 2", 4, false},
		{"(3 + 5) * 2", 16, false},
		{"10 / 2 + 3 * 4", 17, false},
		{"10 % 3", 1, false},
		{"10 / 0", 0, true},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, c := range cases {
		got, err := EvalMathExpr(c.expr)
		if (err != nil) != c.hasErr {
			t.Errorf("EvalMathExpr(%q) error = %v, wantErr = %v", c.expr, err, c.hasErr)
		}
		if err == nil && got != c.expected {
			t.Errorf("EvalMathExpr(%q) = %v; want %v", c.expr, got, c.expected)
		}
	}
}
