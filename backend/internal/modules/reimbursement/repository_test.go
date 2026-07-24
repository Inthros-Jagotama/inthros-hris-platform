package reimbursement

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Reimbursement Type Repository Tests
// =========================================================================

func TestRepo_CreateReimbursementType_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	tp := &ReimbursementType{
		Name:        "Medical Reimbursement",
		Description: "Biaya berobat",
		IsActive:    true,
	}

	if err := repo.CreateReimbursementType(context.Background(), tp); err != nil {
		t.Fatalf("CreateReimbursementType failed: %v", err)
	}

	if tp.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
	if tp.Code != "" {
		t.Errorf("expected empty code, got '%s'", tp.Code)
	}
}

func TestRepo_FindReimbursementTypeByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestReimbursementType(repo)

	found, err := repo.FindReimbursementTypeByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindReimbursementTypeByID failed: %v", err)
	}

	if found.Name != created.Name {
		t.Errorf("expected name '%s', got '%s'", created.Name, found.Name)
	}
}

func TestRepo_FindReimbursementTypeByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.FindReimbursementTypeByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent reimbursement type")
	}
}

func TestRepo_ListReimbursementTypes_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestReimbursementType(repo)
	}

	types, total, err := repo.ListReimbursementTypes(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("ListReimbursementTypes failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(types) != 3 {
		t.Errorf("expected 3 types, got %d", len(types))
	}
}

func TestRepo_UpdateReimbursementType_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestReimbursementType(repo)
	created.Name = "Updated Type"
	created.Code = "UPD"

	if err := repo.UpdateReimbursementType(context.Background(), created); err != nil {
		t.Fatalf("UpdateReimbursementType failed: %v", err)
	}

	found, _ := repo.FindReimbursementTypeByID(context.Background(), created.ID)
	if found.Name != "Updated Type" {
		t.Errorf("expected name 'Updated Type', got '%s'", found.Name)
	}
	if found.Code != "UPD" {
		t.Errorf("expected code 'UPD', got '%s'", found.Code)
	}
}

func TestRepo_DeleteReimbursementType_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestReimbursementType(repo)

	if err := repo.DeleteReimbursementType(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteReimbursementType failed: %v", err)
	}

	_, err := repo.FindReimbursementTypeByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting reimbursement type")
	}
}

// =========================================================================
// Reimbursement Request Repository Tests
// =========================================================================

func TestRepo_CreateReimbursementRequest_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()

	rr := &ReimbursementRequest{
		EmployeeID:    empID,
		RequestTypeID: rType.ID,
		Title:         "Test Request",
		Description:   "Test description",
		Currency:      "IDR",
		Status:        ReimbStatusDraft,
	}

	if err := repo.CreateReimbursementRequest(context.Background(), rr); err != nil {
		t.Fatalf("CreateReimbursementRequest failed: %v", err)
	}

	if rr.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindReimbursementRequestByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	found, err := repo.FindReimbursementRequestByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindReimbursementRequestByID failed: %v", err)
	}

	if found.Status != ReimbStatusDraft {
		t.Errorf("expected status DRAFT, got '%s'", found.Status)
	}
}

func TestRepo_FindReimbursementRequestByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.FindReimbursementRequestByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent request")
	}
}

func TestRepo_ListReimbursementRequests_ByEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()

	createTestReimbursementRequest(repo, empID, rType.ID)
	createTestReimbursementRequest(repo, empID, rType.ID)

	requests, total, err := repo.ListReimbursementRequests(context.Background(), &empID, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListReimbursementRequests failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 requests, got %d", total)
	}
	if len(requests) != 2 {
		t.Errorf("expected 2 requests, got %d", len(requests))
	}
}

func TestRepo_ListReimbursementRequests_ByStatus(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	status := "DRAFT"

	for i := 0; i < 2; i++ {
		createTestReimbursementRequest(repo, empID, rType.ID)
	}

	requests, total, err := repo.ListReimbursementRequests(context.Background(), nil, &status, 1, 10)
	if err != nil {
		t.Fatalf("ListReimbursementRequests failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 requests, got %d", total)
	}
	if len(requests) != 2 {
		t.Errorf("expected 2 requests, got %d", len(requests))
	}
}

func TestRepo_UpdateReimbursementRequest_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	created.Title = "Updated Title"
	created.Status = ReimbStatusSubmitted

	if err := repo.UpdateReimbursementRequest(context.Background(), created); err != nil {
		t.Fatalf("UpdateReimbursementRequest failed: %v", err)
	}

	found, _ := repo.FindReimbursementRequestByID(context.Background(), created.ID)
	if found.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", found.Title)
	}
	if found.Status != ReimbStatusSubmitted {
		t.Errorf("expected status SUBMITTED, got '%s'", found.Status)
	}
}

func TestRepo_DeleteReimbursementRequest_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	if err := repo.DeleteReimbursementRequest(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteReimbursementRequest failed: %v", err)
	}

	_, err := repo.FindReimbursementRequestByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting request")
	}
}

// =========================================================================
// Reimbursement Item Repository Tests
// =========================================================================

func TestRepo_CreateReimbursementItem_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	item := &ReimbursementItem{
		ReimbursementRequestID: rr.ID,
		ExpenseDate:            "2026-07-15",
		ExpenseType:            "MEDICAL",
		Description:            "Doctor visit",
		Amount:                 250000,
	}

	if err := repo.CreateReimbursementItem(context.Background(), item); err != nil {
		t.Fatalf("CreateReimbursementItem failed: %v", err)
	}

	if item.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindReimbursementItemByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	created := createTestReimbursementItem(repo, rr.ID)

	found, err := repo.FindReimbursementItemByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindReimbursementItemByID failed: %v", err)
	}

	if found.Amount != 250000 {
		t.Errorf("expected amount 250000, got %.2f", found.Amount)
	}
}

func TestRepo_ListReimbursementItems_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	for i := 0; i < 3; i++ {
		createTestReimbursementItem(repo, rr.ID)
	}

	items, err := repo.ListReimbursementItems(context.Background(), rr.ID)
	if err != nil {
		t.Fatalf("ListReimbursementItems failed: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}

func TestRepo_UpdateReimbursementItem_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	created := createTestReimbursementItem(repo, rr.ID)

	created.Amount = 500000
	created.Description = "Updated description"

	if err := repo.UpdateReimbursementItem(context.Background(), created); err != nil {
		t.Fatalf("UpdateReimbursementItem failed: %v", err)
	}

	found, _ := repo.FindReimbursementItemByID(context.Background(), created.ID)
	if found.Amount != 500000 {
		t.Errorf("expected amount 500000, got %.2f", found.Amount)
	}
	if found.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got '%s'", found.Description)
	}
}

func TestRepo_DeleteReimbursementItem_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	created := createTestReimbursementItem(repo, rr.ID)

	if err := repo.DeleteReimbursementItem(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteReimbursementItem failed: %v", err)
	}

	_, err := repo.FindReimbursementItemByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting item")
	}
}

// =========================================================================
// Aggregation Tests
// =========================================================================

func TestRepo_SumReimbursementItems(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	// Create 3 items with different amounts
	items := []struct {
		amount float64
		desc   string
	}{
		{100000, "Item 1"},
		{200000, "Item 2"},
		{300000, "Item 3"},
	}
	for _, it := range items {
		item := &ReimbursementItem{
			ReimbursementRequestID: rr.ID,
			ExpenseDate:            "2026-07-15",
			ExpenseType:            "MEDICAL",
			Description:            it.desc,
			Amount:                 it.amount,
		}
		if err := repo.CreateReimbursementItem(context.Background(), item); err != nil {
			t.Fatalf("failed to create item: %v", err)
		}
	}

	sum, err := repo.SumReimbursementItems(context.Background(), rr.ID)
	if err != nil {
		t.Fatalf("SumReimbursementItems failed: %v", err)
	}

	if sum != 600000 {
		t.Errorf("expected sum 600000, got %.2f", sum)
	}
}
