package performance

import (
	"testing"

	"github.com/google/uuid"
)

// TestOKRRepository_LatestOKRScoreByEmployee verifies plan §12.11: the latest
// COMPLETED/APPROVED OKR evaluation's final score is returned for a given
// employee (input untuk promotion eligibility di module employeemovement).
func TestOKRRepository_LatestOKRScoreByEmployee(t *testing.T) {
	_, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	repo := NewOKRRepository()

	employeeID := uuid.New()
	orgID := uuid.New()
	periodID := uuid.New()

	// Evaluation pertama — DRAFT, tidak boleh terpilih.
	if err := repo.CreateOKREvaluation(db, &OKREvaluation{
		EmployeeID:     employeeID,
		OrganizationID: orgID,
		PeriodID:       periodID,
		Status:         OKRStatusDraft,
		FinalScore:     10,
	}); err != nil {
		t.Fatalf("failed to create draft evaluation: %v", err)
	}

	// Evaluasi COMPLETED dengan skor 95 — harus terpilih.
	completed := &OKREvaluation{
		EmployeeID:     employeeID,
		OrganizationID: orgID,
		PeriodID:       periodID,
		Status:         OKRStatusCompleted,
		FinalScore:     95,
	}
	if err := repo.CreateOKREvaluation(db, completed); err != nil {
		t.Fatalf("failed to create completed evaluation: %v", err)
	}

	score, found, err := repo.LatestOKRScoreByEmployee(db, employeeID)
	if err != nil {
		t.Fatalf("LatestOKRScoreByEmployee failed: %v", err)
	}
	if !found {
		t.Error("expected found=true for employee with completed evaluation")
	}
	if score != 95 {
		t.Errorf("expected score 95, got %v", score)
	}

	// Employee tanpa evaluasi → found=false.
	other, found, err := repo.LatestOKRScoreByEmployee(db, uuid.New())
	if err != nil {
		t.Fatalf("LatestOKRScoreByEmployee (missing) failed: %v", err)
	}
	if found {
		t.Error("expected found=false for employee without evaluations")
	}
	if other != 0 {
		t.Errorf("expected score 0 for missing employee, got %v", other)
	}
}

// TestOKRRepository_LatestOKRScoreByEmployee_Approved verifies APPROVED
// evaluations are also eligible as final scores (selain COMPLETED).
func TestOKRRepository_LatestOKRScoreByEmployee_Approved(t *testing.T) {
	_, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	repo := NewOKRRepository()

	employeeID := uuid.New()
	orgID := uuid.New()
	periodID := uuid.New()

	if err := repo.CreateOKREvaluation(db, &OKREvaluation{
		EmployeeID:     employeeID,
		OrganizationID: orgID,
		PeriodID:       periodID,
		Status:         OKRStatusApproved,
		FinalScore:     87,
	}); err != nil {
		t.Fatalf("failed to create approved evaluation: %v", err)
	}

	score, found, err := repo.LatestOKRScoreByEmployee(db, employeeID)
	if err != nil {
		t.Fatalf("LatestOKRScoreByEmployee failed: %v", err)
	}
	if !found {
		t.Error("expected found=true for APPROVED evaluation")
	}
	if score != 87 {
		t.Errorf("expected score 87, got %v", score)
	}
}
