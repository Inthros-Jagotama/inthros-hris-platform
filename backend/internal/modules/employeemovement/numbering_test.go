package employeemovement

import (
	"context"
	"testing"
)

type fakeNumberingGenerator struct {
	calls  []string
	number string
}

func (f *fakeNumberingGenerator) Generate(ctx context.Context, documentType string) (string, error) {
	f.calls = append(f.calls, documentType)
	return f.number, nil
}

func TestCreateMovementGeneratesNumberWhenBlank(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	gen := &fakeNumberingGenerator{number: "SK/001/HRIS/VIII/2026"}
	svc.SetNumberingService(gen)

	toPosition := uuidStr()
	req := CreateMovementRequest{
		EmployeeID:           uuidStr(),
		MovementType:         "promotion",
		DecisionLetterNumber: "",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
		ToPositionID:         &toPosition,
	}
	resp, err := svc.CreateMovement(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateMovement: %v", err)
	}
	if resp.DecisionLetterNumber != "SK/001/HRIS/VIII/2026" {
		t.Fatalf("expected generated number, got %q", resp.DecisionLetterNumber)
	}
	if len(gen.calls) != 1 || gen.calls[0] != "employee_movement" {
		t.Fatalf("expected one Generate call for employee_movement, got %v", gen.calls)
	}
}

func TestCreateMovementKeepsManualNumber(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	gen := &fakeNumberingGenerator{number: "SK/999/HRIS/VIII/2026"}
	svc.SetNumberingService(gen)

	toPosition := uuidStr()
	req := CreateMovementRequest{
		EmployeeID:           uuidStr(),
		MovementType:         "promotion",
		DecisionLetterNumber: "SK/MANUAL/2026",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
		ToPositionID:         &toPosition,
	}
	resp, err := svc.CreateMovement(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateMovement: %v", err)
	}
	if resp.DecisionLetterNumber != "SK/MANUAL/2026" {
		t.Fatalf("expected manual number preserved, got %q", resp.DecisionLetterNumber)
	}
	if len(gen.calls) != 0 {
		t.Fatalf("expected Generate not called when number provided, got %v", gen.calls)
	}
}
