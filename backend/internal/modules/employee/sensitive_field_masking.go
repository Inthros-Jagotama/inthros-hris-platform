package employee

import (
	"context"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
	"github.com/inthros/hris-platform/internal/pkg/mask"
)

// isSelfEmployee mengecek apakah employeeID adalah employee yang terhubung
// ke user yang sedang login (lihat Repository.FindEmployeeIDByUserID) —
// dipakai untuk melewati data masking sensitive field pada data milik sendiri.
func (s *Service) isSelfEmployee(ctx context.Context, employeeID uuid.UUID) bool {
	userID := authctx.GetUserID(ctx)
	if userID == nil {
		return false
	}
	ownID, err := s.repo.FindEmployeeIDByUserID(ctx, *userID)
	if err != nil || ownID == nil {
		return false
	}
	return *ownID == employeeID
}

// isSelfEmployeeStr adalah varian isSelfEmployee yang menerima employeeID
// sebagai string (dipakai sub-resource handler yang employeeID-nya masih
// string mentah dari path param) — parse gagal dianggap bukan self.
func (s *Service) isSelfEmployeeStr(ctx context.Context, employeeID string) bool {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return false
	}
	return s.isSelfEmployee(ctx, uid)
}

// maskField menyamarkan value jika caller tidak punya permission
// "resource.action" untuk field terkait.
func maskField(ctx context.Context, resource, action string, value *string) {
	if *value == "" {
		return
	}
	if authctx.HasPermission(ctx, resource, action) {
		return
	}
	*value = mask.PartialMask(*value)
}

// maskEmployeeResponse menyamarkan field sensitif sesuai permission viewer,
// KECUALI isSelf true (viewer sedang melihat datanya sendiri) — data milik
// sendiri tidak pernah disamarkan, terlepas dari permission view_* yang dimiliki.
func maskEmployeeResponse(ctx context.Context, r *EmployeeResponse, isSelf bool) {
	if r == nil || isSelf {
		return
	}
	maskField(ctx, "employee", "view_nik", &r.NIK)
	maskField(ctx, "employee", "view_passport", &r.Passport)
	maskField(ctx, "employee", "view_phone_number", &r.PhoneNumber)
	maskField(ctx, "employee", "view_email", &r.Email)

	for i := range r.Families {
		maskFamilyResponse(ctx, &r.Families[i], false)
	}
	for i := range r.Banks {
		maskBankResponse(ctx, &r.Banks[i], false)
	}
	for i := range r.EmergencyContacts {
		maskEmergencyContactResponse(ctx, &r.EmergencyContacts[i], false)
	}
}

func maskFamilyResponse(ctx context.Context, r *FamilyResponse, isSelf bool) {
	if r == nil || isSelf {
		return
	}
	maskField(ctx, "employee_family", "view_nik", &r.NIK)
}

func maskBankResponse(ctx context.Context, r *BankResponse, isSelf bool) {
	if r == nil || isSelf {
		return
	}
	maskField(ctx, "employee_bank_account", "view_account_number", &r.AccountNumber)
	maskField(ctx, "employee_bank_account", "view_account_name", &r.AccountName)
}

func maskEmergencyContactResponse(ctx context.Context, r *EmergencyContactResponse, isSelf bool) {
	if r == nil || isSelf {
		return
	}
	maskField(ctx, "emergency_contact", "view_phone_number", &r.PhoneNumber)
}
