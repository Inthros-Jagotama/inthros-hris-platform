package reimbursement

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func ctx() context.Context {
	return context.Background()
}

// =========================================================================
// Reimbursement Type Service Tests
// =========================================================================

func TestService_CreateReimbursementType_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateReimbursementTypeRequest{
		Code:        "MED",
		Name:        "Medical",
		Description: "Medical reimbursement",
		IsActive:    boolPtr(true),
	}

	resp, err := svc.CreateReimbursementType(ctx(), req)
	if err != nil {
		t.Fatalf("CreateReimbursementType failed: %v", err)
	}

	if resp.Name != "Medical" {
		t.Errorf("expected name 'Medical', got '%s'", resp.Name)
	}
	if !resp.IsActive {
		t.Error("expected IsActive = true")
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
}

// TestService_CreateReimbursementType_AutoGeneratesCode guards the
// auto-generated-code behavior: when the form doesn't send a code, the
// service derives one from the display name (UPPER_SNAKE_CASE) and appends
// a numeric suffix when the base code already exists.
func TestService_CreateReimbursementType_AutoGeneratesCode(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	first, err := svc.CreateReimbursementType(ctx(), CreateReimbursementTypeRequest{
		Name:        "Transport & Travel",
		Description: "Transport reimbursement",
	})
	if err != nil {
		t.Fatalf("CreateReimbursementType failed: %v", err)
	}
	if first.Code != "TRANSPORT_TRAVEL" {
		t.Errorf("expected code 'TRANSPORT_TRAVEL', got '%s'", first.Code)
	}

	second, err := svc.CreateReimbursementType(ctx(), CreateReimbursementTypeRequest{
		Name: "Transport & Travel",
	})
	if err != nil {
		t.Fatalf("CreateReimbursementType (duplicate) failed: %v", err)
	}
	if second.Code != "TRANSPORT_TRAVEL_2" {
		t.Errorf("expected code 'TRANSPORT_TRAVEL_2', got '%s'", second.Code)
	}
}

func TestService_GetReimbursementTypeByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestReimbursementType(repo)

	found, err := svc.GetReimbursementTypeByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetReimbursementTypeByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_ListReimbursementTypes_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestReimbursementType(repo)
	}

	resp, err := svc.ListReimbursementTypes(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListReimbursementTypes failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_UpdateReimbursementType_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestReimbursementType(repo)
	newName := "Updated Type"
	req := UpdateReimbursementTypeRequest{Name: &newName}

	updated, err := svc.UpdateReimbursementType(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateReimbursementType failed: %v", err)
	}

	if updated.Name != "Updated Type" {
		t.Errorf("expected name 'Updated Type', got '%s'", updated.Name)
	}
}

func TestService_DeleteReimbursementType_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestReimbursementType(repo)

	if err := svc.DeleteReimbursementType(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteReimbursementType failed: %v", err)
	}

	_, err := svc.GetReimbursementTypeByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting reimbursement type")
	}
}

// =========================================================================
// Reimbursement Request Service Tests
// =========================================================================

func TestService_CreateReimbursementRequest_Success(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	userID := uuid.New()
	empID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)

	req := CreateReimbursementRequest{
		RequestTypeID: rType.ID.String(),
		Title:         "Medical Reimbursement July",
		Description:   "Biaya berobat di klinik",
		Currency:      "IDR",
	}

	resp, err := svc.CreateReimbursementRequest(context.WithValue(ctx(), "user_id", userID.String()), req)
	if err != nil {
		t.Fatalf("CreateReimbursementRequest failed: %v", err)
	}

	if resp.Title != "Medical Reimbursement July" {
		t.Errorf("expected title 'Medical Reimbursement July', got '%s'", resp.Title)
	}
	if resp.Status != "DRAFT" {
		t.Errorf("expected status DRAFT, got '%s'", resp.Status)
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
	// employee_id harus employee asli hasil resolve user_id, bukan user_id itu sendiri.
	if resp.EmployeeID != empID.String() {
		t.Errorf("expected employee_id %s, got %s", empID.String(), resp.EmployeeID)
	}
}

// TestService_CreateReimbursementRequest_UserWithoutEmployeeLink fails loudly
// when the logged-in user has no employee_accounts row — the JWT user_id is
// the account UUID, not the employee UUID, so attribution would be wrong.
func TestService_CreateReimbursementRequest_UserWithoutEmployeeLink(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)

	req := CreateReimbursementRequest{
		RequestTypeID: rType.ID.String(),
		Title:         "Test",
	}

	_, err := svc.CreateReimbursementRequest(context.WithValue(ctx(), "user_id", uuid.New().String()), req)
	if err == nil {
		t.Fatal("expected error when user has no linked employee account")
	}
}

func TestService_CreateReimbursementRequest_MissingUserContext(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)

	req := CreateReimbursementRequest{
		RequestTypeID: rType.ID.String(),
		Title:         "Test",
	}

	_, err := svc.CreateReimbursementRequest(ctx(), req)
	if err == nil {
		t.Fatal("expected error when user context is missing")
	}
}

func TestService_GetReimbursementRequestByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	found, err := svc.GetReimbursementRequestByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetReimbursementRequestByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_ListReimbursementRequests_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	createTestReimbursementRequest(repo, empID, rType.ID)
	createTestReimbursementRequest(repo, empID, rType.ID)

	resp, err := svc.ListReimbursementRequests(ctx(), nil, nil, 0, 0)
	if err != nil {
		t.Fatalf("ListReimbursementRequests failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestService_UpdateReimbursementRequest_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	newTitle := "Updated Title"
	req := UpdateReimbursementRequest{Title: &newTitle}

	updated, err := svc.UpdateReimbursementRequest(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateReimbursementRequest failed: %v", err)
	}

	if updated.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", updated.Title)
	}
}

func TestService_UpdateReimbursementRequest_NonDraftFails(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	// Submit first
	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("UpdateReimbursementRequestStatus failed: %v", err)
	}

	newTitle := "Should Fail"
	req := UpdateReimbursementRequest{Title: &newTitle}

	_, err = svc.UpdateReimbursementRequest(ctx(), created.ID.String(), req)
	if err == nil {
		t.Fatal("expected error when updating non-draft request")
	}
}

// =========================================================================
// Status Flow Tests
// =========================================================================

func TestService_SubmitReimbursementRequest_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	// Add an item first
	createTestReimbursementItem(repo, created.ID)

	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("Submit reimbursement request failed: %v", err)
	}

	if updated.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", updated.Status)
	}
	if updated.SubmittedAt == nil {
		t.Error("expected SubmittedAt to be set")
	}
	if updated.TotalAmount != 250000 {
		t.Errorf("expected TotalAmount 250000, got %.2f", updated.TotalAmount)
	}
}

func TestService_ApproveReimbursementRequest_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	// Submit first
	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	// Then approve
	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "APPROVED", "Approved by HR", nil, nil)
	if err != nil {
		t.Fatalf("Approve reimbursement request failed: %v", err)
	}

	if updated.Status != "APPROVED" {
		t.Errorf("expected status APPROVED, got '%s'", updated.Status)
	}
	if updated.ApprovedAt == nil {
		t.Error("expected ApprovedAt to be set")
	}
}

func TestService_RejectReimbursementRequest_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	// Submit first
	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	// Then reject
	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "REJECTED", "Receipt not valid", nil, nil)
	if err != nil {
		t.Fatalf("Reject reimbursement request failed: %v", err)
	}

	if updated.Status != "REJECTED" {
		t.Errorf("expected status REJECTED, got '%s'", updated.Status)
	}
	if updated.RejectedAt == nil {
		t.Error("expected RejectedAt to be set")
	}
}

func TestService_PayReimbursementRequest_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	// Submit → approve → pay
	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}
	_, err = svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "APPROVED", "Approved", nil, nil)
	if err != nil {
		t.Fatalf("Approve failed: %v", err)
	}

	paidAmount := 250000.00
	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "PAID", "Paid via transfer", float64Ptr(paidAmount), nil)
	if err != nil {
		t.Fatalf("Pay reimbursement request failed: %v", err)
	}

	if updated.Status != "PAID" {
		t.Errorf("expected status PAID, got '%s'", updated.Status)
	}
	if updated.PaidAt == nil {
		t.Error("expected PaidAt to be set")
	}
	if updated.PaidAmount == nil || *updated.PaidAmount != paidAmount {
		t.Errorf("expected PaidAmount %.2f, got %.2f", paidAmount, *updated.PaidAmount)
	}
}

func TestService_PayReimbursementRequest_RecordsPaymentDetails(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}
	_, err = svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "APPROVED", "Approved", nil, nil)
	if err != nil {
		t.Fatalf("Approve failed: %v", err)
	}

	method := "BANK_TRANSFER"
	reference := "TRX-2026-0816-001"
	note := "Transfer BCA a.n. karyawan"
	updated, err := svc.UpdateReimbursementRequestStatus(
		ctx(), created.ID.String(), "PAID", "", nil, nil,
		&method, &reference, &note,
	)
	if err != nil {
		t.Fatalf("Pay reimbursement request failed: %v", err)
	}

	if updated.Status != "PAID" {
		t.Errorf("expected status PAID, got '%s'", updated.Status)
	}
	if updated.PaymentMethod == nil || *updated.PaymentMethod != method {
		t.Errorf("expected payment_method %q, got %v", method, updated.PaymentMethod)
	}
	if updated.PaymentReference == nil || *updated.PaymentReference != reference {
		t.Errorf("expected payment_reference %q, got %v", reference, updated.PaymentReference)
	}
	if updated.PaymentNote == nil || *updated.PaymentNote != note {
		t.Errorf("expected payment_note %q, got %v", note, updated.PaymentNote)
	}

	// Reload from repo to confirm persistence.
	reloaded, err := repo.FindReimbursementRequestByID(ctx(), created.ID)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if reloaded.PaymentMethod == nil || *reloaded.PaymentMethod != method {
		t.Errorf("persisted payment_method mismatch: %v", reloaded.PaymentMethod)
	}
	if reloaded.PaymentReference == nil || *reloaded.PaymentReference != reference {
		t.Errorf("persisted payment_reference mismatch: %v", reloaded.PaymentReference)
	}
	if reloaded.PaymentNote == nil || *reloaded.PaymentNote != note {
		t.Errorf("persisted payment_note mismatch: %v", reloaded.PaymentNote)
	}
}

func TestService_CancelReimbursementRequest_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "CANCELLED", "No longer needed", nil, nil)
	if err != nil {
		t.Fatalf("Cancel reimbursement request failed: %v", err)
	}

	if updated.Status != "CANCELLED" {
		t.Errorf("expected status CANCELLED, got '%s'", updated.Status)
	}
	if updated.CancelledAt == nil {
		t.Error("expected CancelledAt to be set")
	}
}

func TestService_ApproveWithoutSubmit_Fails(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	// Try to approve directly without submitting
	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "APPROVED", "", nil, nil)
	if err == nil {
		t.Fatal("expected error when approving without prior submission")
	}
}

func TestService_PayWithoutApprove_Fails(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	// Submit but don't approve, try to pay directly
	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	_, err = svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "PAID", "", nil, nil)
	if err == nil {
		t.Fatal("expected error when paying without prior approval")
	}
}

func TestService_InvalidStatusTransition_Fails(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "INVALID_STATUS", "", nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

// =========================================================================
// Reimbursement Item Service Tests
// =========================================================================

func TestService_CreateReimbursementItem_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	req := CreateReimbursementItemRequest{
		ExpenseDate: "2026-07-15",
		ExpenseType: "MEDICAL",
		Description: "Doctor consultation",
		Amount:      150000,
	}

	resp, err := svc.CreateReimbursementItem(ctx(), rr.ID.String(), req)
	if err != nil {
		t.Fatalf("CreateReimbursementItem failed: %v", err)
	}

	if resp.Amount != 150000 {
		t.Errorf("expected amount 150000, got %.2f", resp.Amount)
	}
	if resp.ReimbursementRequestID != rr.ID.String() {
		t.Errorf("expected request ID '%s', got '%s'", rr.ID.String(), resp.ReimbursementRequestID)
	}
}

func TestService_CreateReimbursementItem_NonDraftFails(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	// Submit first
	_, err := svc.UpdateReimbursementRequestStatus(ctx(), rr.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	req := CreateReimbursementItemRequest{
		ExpenseDate: "2026-07-15",
		ExpenseType: "MEDICAL",
		Amount:      100000,
	}

	_, err = svc.CreateReimbursementItem(ctx(), rr.ID.String(), req)
	if err == nil {
		t.Fatal("expected error when adding items to submitted request")
	}
}

func TestService_ListReimbursementItems_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	for i := 0; i < 3; i++ {
		createTestReimbursementItem(repo, rr.ID)
	}

	items, err := svc.ListReimbursementItems(ctx(), rr.ID.String())
	if err != nil {
		t.Fatalf("ListReimbursementItems failed: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}

func TestService_UpdateReimbursementItem_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	item := createTestReimbursementItem(repo, rr.ID)

	newAmount := 500000.0
	req := UpdateReimbursementItemRequest{Amount: &newAmount}

	updated, err := svc.UpdateReimbursementItem(ctx(), rr.ID.String(), item.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateReimbursementItem failed: %v", err)
	}

	if updated.Amount != 500000 {
		t.Errorf("expected amount 500000, got %.2f", updated.Amount)
	}
}

func TestService_DeleteReimbursementItem_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	item := createTestReimbursementItem(repo, rr.ID)

	if err := svc.DeleteReimbursementItem(ctx(), rr.ID.String(), item.ID.String()); err != nil {
		t.Fatalf("DeleteReimbursementItem failed: %v", err)
	}

	items, _ := svc.ListReimbursementItems(ctx(), rr.ID.String())
	if len(items) != 0 {
		t.Errorf("expected 0 items after delete, got %d", len(items))
	}
}
