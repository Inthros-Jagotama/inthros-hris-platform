package calculator

import (
	"fmt"
	"sort"
	"strings"
)

// Engine adalah entry point Formula Engine payroll. Engine bersifat stateless
// dan aman dipakai bersamaan: parse/evaluate murni berdasarkan input.
type Engine struct {
	registry *Registry
}

// NewEngine membuat Engine dengan registry built-in default.
func NewEngine() *Engine {
	return &Engine{registry: DefaultRegistry()}
}

// NewEngineWithRegistry membuat Engine dengan registry kustom.
func NewEngineWithRegistry(registry *Registry) *Engine {
	if registry == nil {
		registry = DefaultRegistry()
	}
	return &Engine{registry: registry}
}

// Registry mengembalikan registry variabel engine.
func (e *Engine) Registry() *Registry {
	return e.registry
}

// Parse meng-parse formula menjadi AST.
func (e *Engine) Parse(expr string) (Node, error) {
	if strings.TrimSpace(expr) == "" {
		return nil, fmt.Errorf("formula tidak boleh kosong")
	}
	return Parse(expr)
}

// Validate memastikan formula valid secara sintaks. Tidak mengecek apakah
// variabel tersedia (itu urusan konteks kalkulasi / validasi dependency).
func (e *Engine) Validate(expr string) error {
	_, err := e.Parse(expr)
	return err
}

// Evaluate meng-parse lalu mengevaluasi formula dengan resolver variabel.
func (e *Engine) Evaluate(expr string, resolver VariableResolver) (float64, error) {
	node, err := e.Parse(expr)
	if err != nil {
		return 0, err
	}
	return Evaluate(node, resolver)
}

// ReferencedVariables mengembalikan daftar variabel (built-in + referensi
// komponen lain) yang dipakai formula, diurutkan dan tanpa duplikat.
func (e *Engine) ReferencedVariables(expr string) ([]string, error) {
	node, err := e.Parse(expr)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	collectVariables(node, seen)
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

// ValidateReferences memastikan seluruh variabel formula ter-resolve: variabel
// built-in dikenal registry, atau tercantum sebagai kode komponen yang ada.
// dependencies adalah map kode komponen yang valid (biasanya dari DB).
func (e *Engine) ValidateReferences(expr string, availableComponents map[string]bool) ([]string, error) {
	vars, err := e.ReferencedVariables(expr)
	if err != nil {
		return nil, err
	}
	var unresolved []string
	for _, name := range vars {
		if e.registry.IsBuiltIn(name) {
			continue
		}
		if availableComponents[normalizeVarName(name)] {
			continue
		}
		unresolved = append(unresolved, name)
	}
	if len(unresolved) > 0 {
		return vars, fmt.Errorf("variabel tak dikenal pada formula: %s", strings.Join(unresolved, ", "))
	}
	return vars, nil
}

func collectVariables(node Node, seen map[string]bool) {
	switch n := node.(type) {
	case *VariableNode:
		seen[normalizeVarName(n.Name)] = true
	case *PercentNode:
		collectVariables(n.Operand, seen)
	case *UnaryNode:
		collectVariables(n.Operand, seen)
	case *BinaryOpNode:
		collectVariables(n.Left, seen)
		collectVariables(n.Right, seen)
	}
}

// Cycle adalah siklus dependency yang terdeteksi, mis. A→B→A.
type Cycle struct {
	Path []string `json:"path"` // urutan komponen yang membentuk siklus
}

func (c Cycle) String() string {
	return strings.Join(c.Path, " → ")
}

// DetectCycles mendeteksi circular dependency pada graph komponen.
// deps memetakan kode komponen ke daftar kode yang direferensikan.
// Mengembalikan satu siklus per komponen "akar" siklus (deduplikasi otomatis).
func DetectCycles(deps map[string][]string) []Cycle {
	// State DFS: 0 = belum dikunjungi, 1 = sedang diproses (di stack), 2 = selesai.
	state := map[string]int{}
	var cycles []Cycle
	seenCycles := map[string]bool{}

	var dfs func(node string, stack []string)
	dfs = func(node string, stack []string) {
		state[node] = 1
		stack = append(stack, node)
		for _, dep := range deps[node] {
			dep = normalizeVarName(dep)
			if dep == node {
				// Self-reference: A → A.
				key := "self:" + node
				if !seenCycles[key] {
					seenCycles[key] = true
					cycles = append(cycles, Cycle{Path: []string{node, node}})
				}
				continue
			}
			switch state[dep] {
			case 0:
				dfs(dep, stack)
			case 1:
				// Siklus ditemukan: potong stack dari posisi dep.
				start := -1
				for i, n := range stack {
					if n == dep {
						start = i
						break
					}
				}
				if start >= 0 {
					path := append([]string{}, stack[start:]...)
					path = append(path, dep)
					if key := cycleKey(path); !seenCycles[key] {
						seenCycles[key] = true
						cycles = append(cycles, Cycle{Path: path})
					}
				}
			}
		}
		state[node] = 2
	}

	// Urutkan deterministik agar hasil test stabil.
	keys := make([]string, 0, len(deps))
	for k := range deps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if state[k] == 0 {
			dfs(k, nil)
		}
	}
	return cycles
}

// cycleKey membuat kunci kanonik siklus agar siklus yang sama tidak dihitung
// dua kali walau terdeteksi dari titik masuk berbeda.
func cycleKey(path []string) string {
	if len(path) < 2 {
		return strings.Join(path, "→")
	}
	// Rotasi agar kunci dimulai dari elemen terkecil (leksikografis).
	minIdx := 0
	for i := 1; i < len(path)-1; i++ {
		if path[i] < path[minIdx] {
			minIdx = i
		}
	}
	rotated := append(append([]string{}, path[minIdx:len(path)-1]...), path[:minIdx]...)
	return strings.Join(rotated, "→")
}
