package calculator

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

// =============================================================================
// Parsing
// =============================================================================

func TestParseValidExpressions(t *testing.T) {
	valid := []string{
		"BASIC + POSITION_ALLOWANCE",
		"BPJS_WAGE * 2%",
		"OVERTIME_HOURS * OVERTIME_RATE",
		"GROSS - TOTAL_EMPLOYEE_DEDUCTION",
		"(BASIC + ALLOWANCE) * 12",
		"500000",
		"2.5",
		"BASIC + 10%",
		"TOTAL_EARNINGS - TOTAL_DEDUCTIONS",
		"  BASIC   +   ALLOWANCE  ",
	}
	for _, expr := range valid {
		if _, err := Parse(expr); err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", expr, err)
		}
	}
}

func TestParseInvalidExpressions(t *testing.T) {
	invalid := []string{
		"",
		"   ",
		"BASIC +",
		"* BASIC",
		"(BASIC + ALLOWANCE",
		"BASIC + )",
		"BASIC & ALLOWANCE",
		"BASIC + 2*",
	}
	for _, expr := range invalid {
		if _, err := Parse(expr); err == nil {
			t.Errorf("Parse(%q) expected error, got nil", expr)
		}
	}
}

func TestParseNumberLiteral(t *testing.T) {
	node, err := Parse("2%")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	v, err := Evaluate(node, nil)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if math.Abs(v-0.02) > 1e-9 {
		t.Errorf("expected 0.02, got %v", v)
	}
}

// =============================================================================
// Evaluation
// =============================================================================

func TestEvaluateBasicArithmetic(t *testing.T) {
	values := map[string]float64{
		"BASIC":                    10000000,
		"POSITION_ALLOWANCE":       500000,
		"OVERTIME_HOURS":           10,
		"OVERTIME_RATE":            50000,
		"GROSS":                    10500000,
		"TOTAL_EMPLOYEE_DEDUCTION": 3000000,
	}
	resolver := func(name string) (float64, bool) {
		v, ok := values[name]
		return v, ok
	}

	cases := []struct {
		expr string
		want float64
	}{
		{"BASIC + POSITION_ALLOWANCE", 10500000},
		{"OVERTIME_HOURS * OVERTIME_RATE", 500000},
		{"GROSS - TOTAL_EMPLOYEE_DEDUCTION", 7500000},
		{"(BASIC + POSITION_ALLOWANCE) * 12", 126000000},
		{"BASIC / 1000000", 10},
	}
	for _, c := range cases {
		got, err := EvaluateFormula(c.expr, resolver)
		if err != nil {
			t.Errorf("Evaluate(%q): %v", c.expr, err)
			continue
		}
		if got != c.want {
			t.Errorf("Evaluate(%q) = %v, want %v", c.expr, got, c.want)
		}
	}
}

func TestEvaluatePercentage(t *testing.T) {
	values := map[string]float64{"BPJS_WAGE": 10000000}
	resolver := func(name string) (float64, bool) {
		v, ok := values[name]
		return v, ok
	}
	got, err := EvaluateFormula("BPJS_WAGE * 2%", resolver)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if got != 200000 {
		t.Errorf("expected 200000, got %v", got)
	}
}

func TestEvaluateDivisionByZero(t *testing.T) {
	_, err := EvaluateFormula("BASIC / 0", func(string) (float64, bool) { return 1, true })
	if err == nil {
		t.Fatal("expected division by zero error")
	}
	if !strings.Contains(err.Error(), "nol") {
		t.Errorf("expected error mentioning division by zero, got %q", err.Error())
	}
}

func TestEvaluateUnknownVariable(t *testing.T) {
	_, err := EvaluateFormula("UNKNOWN_VAR + 1", func(string) (float64, bool) { return 0, false })
	if err == nil {
		t.Fatal("expected unknown variable error")
	}
	if !strings.Contains(err.Error(), "UNKNOWN_VAR") {
		t.Errorf("expected error mentioning UNKNOWN_VAR, got %q", err.Error())
	}
}

// =============================================================================
// Engine API
// =============================================================================

func TestEngineValidateAndReferencedVariables(t *testing.T) {
	e := NewEngine()
	if err := e.Validate("BASIC + POSITION_ALLOWANCE"); err != nil {
		t.Fatalf("Validate valid formula: %v", err)
	}
	if err := e.Validate("BASIC +"); err == nil {
		t.Fatal("expected error for invalid formula")
	}

	vars, err := e.ReferencedVariables("BASIC + basic + ALLOWANCE")
	if err != nil {
		t.Fatalf("ReferencedVariables: %v", err)
	}
	// BASIC di-normalize uppercase, deduplikasi, urut.
	if !reflect.DeepEqual(vars, []string{"ALLOWANCE", "BASIC"}) {
		t.Errorf("expected [ALLOWANCE BASIC], got %v", vars)
	}
}

func TestEngineValidateReferences(t *testing.T) {
	e := NewEngine()
	available := map[string]bool{"BASIC": true, "POSITION_ALLOWANCE": true}

	// Semua variabel ter-resolve (built-in atau komponen tersedia).
	if _, err := e.ValidateReferences("BASIC + GROSS", available); err != nil {
		t.Errorf("expected valid references, got %v", err)
	}

	// Komponen yang tidak ada → error.
	if _, err := e.ValidateReferences("BASIC + MISSING_COMP", available); err == nil {
		t.Error("expected unresolved variable error")
	}
}

func TestEngineRegistryBuiltin(t *testing.T) {
	e := NewEngine()
	reg := e.Registry()
	if !reg.IsBuiltIn("gross") {
		t.Error("expected GROSS to be a built-in variable (case-insensitive)")
	}
	if !reg.IsBuiltIn("BASIC") {
		t.Error("expected BASIC as a built-in variable per docs variable registry")
	}
	if reg.IsBuiltIn("SOME_COMPONENT_CODE") {
		t.Error("unregistered component code should not be a built-in variable")
	}
	if !reg.IsBuiltIn("NET_SALARY") {
		t.Error("expected NET_SALARY built-in")
	}
	all := reg.All()
	if len(all) == 0 {
		t.Error("expected non-empty default registry")
	}
	// Pastikan terurut.
	for i := 1; i < len(all); i++ {
		if all[i-1].Name > all[i].Name {
			t.Errorf("registry not sorted at index %d: %s > %s", i, all[i-1].Name, all[i].Name)
		}
	}
}

// =============================================================================
// Circular Dependency Detection
// =============================================================================

func TestDetectCyclesNone(t *testing.T) {
	deps := map[string][]string{
		"BASIC":       {},
		"ALLOWANCE":   {"BASIC"},
		"BONUS":       {"BASIC", "ALLOWANCE"},
		"GROSS_TOTAL": {"BASIC", "ALLOWANCE", "BONUS"},
	}
	cycles := DetectCycles(deps)
	if len(cycles) != 0 {
		t.Errorf("expected no cycles, got %v", cycles)
	}
}

func TestDetectCyclesDirect(t *testing.T) {
	deps := map[string][]string{
		"A": {"B"},
		"B": {"A"},
	}
	cycles := DetectCycles(deps)
	if len(cycles) != 1 {
		t.Fatalf("expected 1 cycle, got %v", cycles)
	}
	if !reflect.DeepEqual(cycles[0].Path, []string{"A", "B", "A"}) {
		t.Errorf("unexpected cycle path: %v", cycles[0].Path)
	}
}

func TestDetectCyclesSelfReference(t *testing.T) {
	deps := map[string][]string{
		"A": {"A"},
	}
	cycles := DetectCycles(deps)
	if len(cycles) != 1 {
		t.Fatalf("expected 1 self-cycle, got %v", cycles)
	}
	if !reflect.DeepEqual(cycles[0].Path, []string{"A", "A"}) {
		t.Errorf("unexpected self-cycle path: %v", cycles[0].Path)
	}
}

func TestDetectCyclesIndirect(t *testing.T) {
	deps := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {"A"},
		"D": {"A"},
	}
	cycles := DetectCycles(deps)
	if len(cycles) != 1 {
		t.Fatalf("expected 1 cycle, got %v", cycles)
	}
	if !reflect.DeepEqual(cycles[0].Path, []string{"A", "B", "C", "A"}) {
		t.Errorf("unexpected cycle path: %v", cycles[0].Path)
	}
}

func TestDetectCyclesMultiple(t *testing.T) {
	deps := map[string][]string{
		"A": {"B"},
		"B": {"A"},
		"C": {"D"},
		"D": {"C"},
	}
	cycles := DetectCycles(deps)
	if len(cycles) != 2 {
		t.Fatalf("expected 2 cycles, got %v", cycles)
	}
}

// =============================================================================
// Rounding
// =============================================================================

func TestRoundingModes(t *testing.T) {
	if got := Round(2.5, RoundingRound); got != 3 {
		t.Errorf("Round(2.5)=%v, want 3", got)
	}
	if got := Round(2.4, RoundingRound); got != 2 {
		t.Errorf("Round(2.4)=%v, want 2", got)
	}
	if got := Round(2.1, RoundingCeil); got != 3 {
		t.Errorf("Ceil(2.1)=%v, want 3", got)
	}
	if got := Round(2.9, RoundingFloor); got != 2 {
		t.Errorf("Floor(2.9)=%v, want 2", got)
	}
	if got := Round(2.5, RoundingNone); got != 2.5 {
		t.Errorf("None(2.5)=%v, want 2.5", got)
	}
}

func TestRoundToUnit(t *testing.T) {
	if got := RoundToUnit(12345, 1000, RoundingRound); got != 12000 {
		t.Errorf("RoundToUnit(12345, 1000, ROUND)=%v, want 12000", got)
	}
	if got := RoundToUnit(12345, 1000, RoundingCeil); got != 13000 {
		t.Errorf("RoundToUnit(12345, 1000, CEIL)=%v, want 13000", got)
	}
	if got := RoundToUnit(12345, 1000, RoundingFloor); got != 12000 {
		t.Errorf("RoundToUnit(12345, 1000, FLOOR)=%v, want 12000", got)
	}
	if got := RoundToUnit(12345, 1000, RoundingNone); got != 12345 {
		t.Errorf("RoundToUnit(12345, 1000, NONE)=%v, want 12345", got)
	}
	if got := RoundToUnit(12345, 0, RoundingRound); got != 12345 {
		t.Errorf("RoundToUnit with unit=0 should be identity, got %v", got)
	}
}

func TestNormalizeRoundingMode(t *testing.T) {
	if NormalizeRoundingMode("") != RoundingRound {
		t.Error("empty mode should default to ROUND")
	}
	if NormalizeRoundingMode("WEIRD") != RoundingRound {
		t.Error("unknown mode should default to ROUND")
	}
	if NormalizeRoundingMode("CEIL") != RoundingCeil {
		t.Error("CEIL should be preserved")
	}
}
