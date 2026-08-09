package attendance

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Overtime Two-Flow (§32b) Service Tests
// =========================================================================

// TestService_AssignOvertimeRequest_Success — alur ASSIGNED: tanpa approval
// penugasan, langsung WAITING_ACTUAL, flow_type ASSIGNED.
func TestService_AssignOvertimeRequest_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := AssignOvertimeRequest{
		AssignedEmployeeID: uuidStr(),
		WorkDate:           "2026-01-15",
		StartTimeLocal:     "2026-01-15T18:00:00+07:00",
		EndTimeLocal:       "2026-01-15T20:00:00+07:00",
		RequestedMinutes:   120,
		Reason:             "Client demo tomorrow",
	}

	resp, err := svc.AssignOvertimeRequest(ctx(), req)
	if err != nil {
		t.Fatalf("AssignOvertimeRequest failed: %v", err)
	}

	if resp.Status != "WAITING_ACTUAL" {
		t.Errorf("expected status WAITING_ACTUAL, got '%s'", resp.Status)
	}
	if resp.FlowType != "ASSIGNED" {
		t.Errorf("expected flow_type ASSIGNED, got '%s'", resp.FlowType)
	}
	if resp.AssignedAt == nil {
		t.Error("expected assigned_at to be set")
	}
	if resp.RequestedMinutes != 120 {
		t.Errorf("expected 120 requested minutes, got %d", resp.RequestedMinutes)
	}
}

func TestService_AssignOvertimeRequest_InvalidEmployeeID(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := AssignOvertimeRequest{
		AssignedEmployeeID: "not-a-uuid",
		WorkDate:           "2026-01-15",
		StartTimeLocal:     "2026-01-15T18:00:00+07:00",
		EndTimeLocal:       "2026-01-15T20:00:00+07:00",
		RequestedMinutes:   120,
	}

	if _, err := svc.AssignOvertimeRequest(ctx(), req); err == nil {
		t.Fatal("expected error for invalid assigned_employee_id")
	}
}

// TestService_SubmitActualOvertime_Success — alur SELF yang sudah di-approve
// instance #1 (WAITING_ACTUAL): isian aktual tersimpan, status ACTUAL_SUBMITTED,
// dan instance approval #2 dibuat (karena flow aktif tersedia).
func TestService_SubmitActualOvertime_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{activeFlowID: uuid.New().String()}
	svc.SetApprovalEngine(fake)

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimeWaitingActual
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	resp, err := svc.SubmitActualOvertime(ctx(), o.ID, SubmitOvertimeActualRequest{
		ActualStartTimeLocal: "2026-01-15T18:00:00+07:00",
		ActualEndTimeLocal:   "2026-01-15T20:30:00+07:00",
		ActualNote:           "Selesai lebih malam dari rencana",
	})
	if err != nil {
		t.Fatalf("SubmitActualOvertime failed: %v", err)
	}

	if resp.Status != "ACTUAL_SUBMITTED" {
		t.Errorf("expected status ACTUAL_SUBMITTED, got '%s'", resp.Status)
	}
	if resp.ActualMinutes == nil || *resp.ActualMinutes != 150 {
		t.Errorf("expected actual_minutes 150, got %v", resp.ActualMinutes)
	}
	if resp.ActualApprovalInstanceID == nil {
		t.Error("expected actual_approval_instance_id to be set (instance #2)")
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].documentID != o.ID.String() {
		t.Errorf("expected documentID %s, got %s", o.ID.String(), fake.createCalls[0].documentID)
	}
}

// TestService_SubmitActualOvertime_WrongState — hanya WAITING_ACTUAL yang boleh
// diisi aktualnya.
func TestService_SubmitActualOvertime_WrongState(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New()) // status SUBMITTED

	_, err := svc.SubmitActualOvertime(ctx(), o.ID, SubmitOvertimeActualRequest{
		ActualStartTimeLocal: "2026-01-15T18:00:00+07:00",
		ActualEndTimeLocal:   "2026-01-15T20:00:00+07:00",
	})
	if !errors.Is(err, ErrOvertimeInvalidState) {
		t.Errorf("expected ErrOvertimeInvalidState, got %v", err)
	}
}

// TestService_SubmitActualOvertime_InvalidRange — jam selesai harus setelah
// jam mulai.
func TestService_SubmitActualOvertime_InvalidRange(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimeWaitingActual
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	_, err := svc.SubmitActualOvertime(ctx(), o.ID, SubmitOvertimeActualRequest{
		ActualStartTimeLocal: "2026-01-15T20:00:00+07:00",
		ActualEndTimeLocal:   "2026-01-15T18:00:00+07:00",
	})
	if !errors.Is(err, ErrOvertimeInvalidActualRange) {
		t.Errorf("expected ErrOvertimeInvalidActualRange, got %v", err)
	}
}

// TestService_SubmitActualOvertime_NoActiveFlow_RevertsToWaitingActual —
// tanpa flow aktif, instance #2 tidak dibuat dan status kembali WAITING_ACTUAL
// (data aktual tetap tersimpan) agar bisa submit ulang.
func TestService_SubmitActualOvertime_NoActiveFlow_RevertsToWaitingActual(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{flowResolveErr: errors.New("no active flow for module")}
	svc.SetApprovalEngine(fake)

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimeWaitingActual
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	resp, err := svc.SubmitActualOvertime(ctx(), o.ID, SubmitOvertimeActualRequest{
		ActualStartTimeLocal: "2026-01-15T18:00:00+07:00",
		ActualEndTimeLocal:   "2026-01-15T20:00:00+07:00",
	})
	if err != nil {
		t.Fatalf("SubmitActualOvertime failed: %v", err)
	}

	if resp.Status != "WAITING_ACTUAL" {
		t.Errorf("expected status reverted to WAITING_ACTUAL, got '%s'", resp.Status)
	}
	if resp.ActualMinutes == nil || *resp.ActualMinutes != 120 {
		t.Errorf("expected actual data preserved (120m), got %v", resp.ActualMinutes)
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}
}

// TestService_CancelOvertimeRequest_PendingApproval — batal sebelum isian
// aktual (alur SELF, request belum di-approve) → CANCELLED + instance aktif
// dibatalkan.
func TestService_CancelOvertimeRequest_PendingApproval(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimePendingApproval
	instanceID := uuid.New()
	o.ApprovalInstanceID = &instanceID
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.CancelOvertimeRequest(ctx(), o.ID); err != nil {
		t.Fatalf("CancelOvertimeRequest failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.Status != "CANCELLED" {
		t.Errorf("expected status CANCELLED, got '%s'", updated.Status)
	}
	if len(fake.cancelCalls) != 1 || fake.cancelCalls[0] != instanceID.String() {
		t.Errorf("expected CancelApprovalInstance called with %s, got %v", instanceID.String(), fake.cancelCalls)
	}
}

// TestService_CancelOvertimeRequest_WaitingActual — batal di alur ASSIGNED
// (atau SELF setelah approve #1, sebelum isian aktual) juga diperbolehkan.
func TestService_CancelOvertimeRequest_WaitingActual(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimeWaitingActual
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.CancelOvertimeRequest(ctx(), o.ID); err != nil {
		t.Fatalf("CancelOvertimeRequest failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.Status != "CANCELLED" {
		t.Errorf("expected status CANCELLED, got '%s'", updated.Status)
	}
}

// TestService_CancelOvertimeRequest_InvalidState — tidak bisa batal setelah
// isian aktual disubmit (ACTUAL_SUBMITTED) atau sudah final.
func TestService_CancelOvertimeRequest_InvalidState(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimeActualSubmitted
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	err := svc.CancelOvertimeRequest(ctx(), o.ID)
	if !errors.Is(err, ErrOvertimeInvalidState) {
		t.Errorf("expected ErrOvertimeInvalidState, got %v", err)
	}
}

// TestService_SubmitActualOvertime_NotOwner — hanya employee pemilik request
// yang boleh mengisi aktual. Tanpa user login (authctx kosong) ownership
// guard di-skip, jadi test ini hanya memastikan jalur dengan user lain ditolak
// ketika repository mampu me-resolve user → employee.
func TestService_SubmitActualOvertime_NotOwner(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	ownerEmpID := uuid.New()
	seedEmployeeOrg(db, ownerEmpID, "Owner", uuid.New(), "Org A")

	o := createTestOvertimeRequest(repo, ownerEmpID)
	o.Status = OvertimeWaitingActual
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	// Login sebagai employee lain (bukan pemilik).
	otherUserID := uuid.New()
	otherEmpID := uuid.New()
	seedUserEmployee(db, otherUserID, otherEmpID)

	cctx := ctxWithUser(otherUserID)
	_, err := svc.SubmitActualOvertime(cctx, o.ID, SubmitOvertimeActualRequest{
		ActualStartTimeLocal: "2026-01-15T18:00:00+07:00",
		ActualEndTimeLocal:   "2026-01-15T20:00:00+07:00",
	})
	if !errors.Is(err, ErrOvertimeNotOwner) {
		t.Errorf("expected ErrOvertimeNotOwner, got %v", err)
	}
}

// =========================================================================
// Assignable Employees (§32b alur ASSIGNED) Service Tests
// =========================================================================

// TestService_ListAssignableEmployees_DirectChildren — user dengan anak org
// langsung yang terisi: hanya karyawan anak langsung yang muncul.
func TestService_ListAssignableEmployees_DirectChildren(t *testing.T) {
	svc, _, db, cleanup := newTestService()
	defer cleanup()

	userID := uuid.New()
	myEmpID := uuid.New()
	myOrgID := uuid.New()
	childOrgID := uuid.New()
	seedOrg(db, myOrgID, nil, "Root")
	seedOrg(db, childOrgID, &myOrgID, "Child")
	seedUserEmployee(db, userID, myEmpID)
	seedEmployment(db, myEmpID, myOrgID)
	childEmpID := uuid.New()
	seedEmployeeRow(db, childEmpID, "Child Employee")
	seedEmployment(db, childEmpID, childOrgID)

	resp, err := svc.ListAssignableEmployees(ctxWithUser(userID))
	if err != nil {
		t.Fatalf("ListAssignableEmployees failed: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 assignable employee, got %d: %+v", len(resp), resp)
	}
	if resp[0].EmployeeID != childEmpID.String() {
		t.Errorf("expected employee %s, got %s", childEmpID.String(), resp[0].EmployeeID)
	}
	if resp[0].EmployeeCode == "" {
		t.Errorf("expected employee_code to be populated, got empty")
	}
}

// TestService_ListAssignableEmployees_SkipVacantLevels — anak langsung kosong:
// turun satu level lagi sampai menemukan org terisi ("bawahan langsung, atau
// cari ke level bawah berikutnya sampai ada karyawannya").
func TestService_ListAssignableEmployees_SkipVacantLevels(t *testing.T) {
	svc, _, db, cleanup := newTestService()
	defer cleanup()

	userID := uuid.New()
	myEmpID := uuid.New()
	myOrgID := uuid.New()
	vacantChild := uuid.New()
	occupiedGrandchild := uuid.New()
	seedOrg(db, myOrgID, nil, "Root")
	seedOrg(db, vacantChild, &myOrgID, "Vacant Child")
	seedOrg(db, occupiedGrandchild, &vacantChild, "Occupied Grandchild")
	seedUserEmployee(db, userID, myEmpID)
	seedEmployment(db, myEmpID, myOrgID)
	grandchildEmpID := uuid.New()
	seedEmployeeRow(db, grandchildEmpID, "Grandchild Employee")
	seedEmployment(db, grandchildEmpID, occupiedGrandchild)

	resp, err := svc.ListAssignableEmployees(ctxWithUser(userID))
	if err != nil {
		t.Fatalf("ListAssignableEmployees failed: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 assignable employee (from grandchild), got %d: %+v", len(resp), resp)
	}
	if resp[0].EmployeeID != grandchildEmpID.String() {
		t.Errorf("expected grandchild employee %s, got %s", grandchildEmpID.String(), resp[0].EmployeeID)
	}
	if resp[0].EmployeeCode == "" {
		t.Errorf("expected employee_code to be populated, got empty")
	}
}

// TestService_ListAssignableEmployees_NoActiveEmployment — user tanpa
// employment aktif: daftar kosong (bukan error).
func TestService_ListAssignableEmployees_NoActiveEmployment(t *testing.T) {
	svc, _, db, cleanup := newTestService()
	defer cleanup()

	userID := uuid.New()
	myEmpID := uuid.New()
	seedUserEmployee(db, userID, myEmpID)

	resp, err := svc.ListAssignableEmployees(ctxWithUser(userID))
	if err != nil {
		t.Fatalf("ListAssignableEmployees failed: %v", err)
	}
	if len(resp) != 0 {
		t.Errorf("expected empty list, got %+v", resp)
	}
}

// TestService_AssignOvertimeRequest_NotSubordinate_Rejected — server-side
// guard: employee di luar bawahan efektif tidak bisa ditugaskan.
func TestService_AssignOvertimeRequest_NotSubordinate_Rejected(t *testing.T) {
	svc, _, db, cleanup := newTestService()
	defer cleanup()

	userID := uuid.New()
	myEmpID := uuid.New()
	myOrgID := uuid.New()
	seedOrg(db, myOrgID, nil, "Root")
	seedUserEmployee(db, userID, myEmpID)
	seedEmployment(db, myEmpID, myOrgID)
	outsiderEmpID := uuid.New()
	outsiderOrgID := uuid.New()
	seedOrg(db, outsiderOrgID, nil, "Sibling Root")
	seedEmployment(db, outsiderEmpID, outsiderOrgID)

	_, err := svc.AssignOvertimeRequest(ctxWithUser(userID), AssignOvertimeRequest{
		AssignedEmployeeID: outsiderEmpID.String(),
		WorkDate:           "2026-01-15",
		StartTimeLocal:     "2026-01-15T18:00:00+07:00",
		EndTimeLocal:       "2026-01-15T20:00:00+07:00",
		RequestedMinutes:   120,
	})
	if !errors.Is(err, ErrOvertimeNotAssignable) {
		t.Fatalf("expected ErrOvertimeNotAssignable, got %v", err)
	}
}

// TestService_AssignOvertimeRequest_Subordinate_OK — employee di org bawahan
// efektif boleh ditugaskan.
func TestService_AssignOvertimeRequest_Subordinate_OK(t *testing.T) {
	svc, _, db, cleanup := newTestService()
	defer cleanup()

	userID := uuid.New()
	myEmpID := uuid.New()
	myOrgID := uuid.New()
	childOrgID := uuid.New()
	seedOrg(db, myOrgID, nil, "Root")
	seedOrg(db, childOrgID, &myOrgID, "Child")
	seedUserEmployee(db, userID, myEmpID)
	seedEmployment(db, myEmpID, myOrgID)
	childEmpID := uuid.New()
	seedEmployeeRow(db, childEmpID, "Child Employee")
	seedEmployment(db, childEmpID, childOrgID)

	resp, err := svc.AssignOvertimeRequest(ctxWithUser(userID), AssignOvertimeRequest{
		AssignedEmployeeID: childEmpID.String(),
		WorkDate:           "2026-01-15",
		StartTimeLocal:     "2026-01-15T18:00:00+07:00",
		EndTimeLocal:       "2026-01-15T20:00:00+07:00",
		RequestedMinutes:   120,
	})
	if err != nil {
		t.Fatalf("AssignOvertimeRequest failed: %v", err)
	}
	if resp.Status != "WAITING_ACTUAL" {
		t.Errorf("expected status WAITING_ACTUAL, got '%s'", resp.Status)
	}
	if resp.EmployeeID != childEmpID.String() {
		t.Errorf("expected assigned employee %s, got %s", childEmpID.String(), resp.EmployeeID)
	}
}
