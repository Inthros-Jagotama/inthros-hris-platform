package calculator

import "fmt"

// Node adalah representasi AST (Abstract Syntax Tree) dari sebuah formula.
// Tidak ada evaluasi via eval()/reflect — semuanya di-parse dan dievaluasi
// secara eksplisit (lihat parser.go / evaluator.go).
type Node interface {
	// String mengembalikan representasi infix dari node (berguna untuk
	// menampilkan kembali formula dan debugging).
	String() string
}

// NumberNode menyimpan literal angka (mis. 500000, 2).
type NumberNode struct {
	Value float64
}

func (n *NumberNode) String() string {
	if n.Value == float64(int64(n.Value)) {
		return fmt.Sprintf("%d", int64(n.Value))
	}
	return fmt.Sprintf("%g", n.Value)
}

// VariableNode mereferensikan sebuah variabel: variabel built-in dari registry
// (mis. GROSS, BPJS_WAGE) atau kode salary component lain (mis. BASIC).
type VariableNode struct {
	Name string
}

func (n *VariableNode) String() string { return n.Name }

// BinaryOpNode adalah operasi biner: + - * /
type BinaryOpNode struct {
	Op    TokenType
	Left  Node
	Right Node
}

func (n *BinaryOpNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Op.String(), n.Right.String())
}

// UnaryNode adalah operasi unary: -x atau +x
type UnaryNode struct {
	Op      TokenType
	Operand Node
}

func (n *UnaryNode) String() string {
	return fmt.Sprintf("(%s%s)", n.Op.String(), n.Operand.String())
}

// PercentNode menerapkan postfix '%': nilai operand dibagi 100 (mis. 2% = 0.02).
type PercentNode struct {
	Operand Node
}

func (n *PercentNode) String() string {
	return fmt.Sprintf("%s%%", n.Operand.String())
}
