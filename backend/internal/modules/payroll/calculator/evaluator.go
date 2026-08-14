package calculator

import "fmt"

// VariableResolver menyediakan nilai sebuah variabel saat evaluasi.
// Return (0, false) jika variabel tidak diketahui pada konteks kalkulasi.
type VariableResolver func(name string) (float64, bool)

// Evaluate mengevaluasi AST dengan resolver variabel yang diberikan.
// Mengembalikan error untuk: variabel tak dikenal, pembagian dengan nol.
func Evaluate(node Node, resolver VariableResolver) (float64, error) {
	if resolver == nil {
		resolver = func(string) (float64, bool) { return 0, false }
	}
	return eval(node, resolver)
}

func eval(node Node, resolver VariableResolver) (float64, error) {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value, nil

	case *VariableNode:
		value, ok := resolver(n.Name)
		if !ok {
			return 0, fmt.Errorf("variabel %q tidak tersedia pada konteks kalkulasi", n.Name)
		}
		return value, nil

	case *PercentNode:
		v, err := eval(n.Operand, resolver)
		if err != nil {
			return 0, err
		}
		return v / 100, nil

	case *UnaryNode:
		v, err := eval(n.Operand, resolver)
		if err != nil {
			return 0, err
		}
		if n.Op == TokenMinus {
			return -v, nil
		}
		return v, nil

	case *BinaryOpNode:
		left, err := eval(n.Left, resolver)
		if err != nil {
			return 0, err
		}
		right, err := eval(n.Right, resolver)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case TokenPlus:
			return left + right, nil
		case TokenMinus:
			return left - right, nil
		case TokenStar:
			return left * right, nil
		case TokenSlash:
			if right == 0 {
				return 0, fmt.Errorf("pembagian dengan nol pada formula")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("operator tidak didukung: %s", n.Op.String())
		}

	default:
		return 0, fmt.Errorf("node formula tidak dikenal: %T", node)
	}
}

// EvaluateFormula meng-parse lalu mengevaluasi formula string dalam satu langkah.
func EvaluateFormula(expr string, resolver VariableResolver) (float64, error) {
	node, err := Parse(expr)
	if err != nil {
		return 0, err
	}
	return Evaluate(node, resolver)
}
