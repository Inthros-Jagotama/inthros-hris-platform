package employee

import (
	"context"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
	"github.com/inthros/hris-platform/internal/pkg/mask"
)

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

func maskEmployeeResponse(ctx context.Context, r *EmployeeResponse) {
	if r == nil {
		return
	}
	maskField(ctx, "employee", "view_nik", &r.NIK)
	maskField(ctx, "employee", "view_passport", &r.Passport)
	maskField(ctx, "employee", "view_phone_number", &r.PhoneNumber)
	maskField(ctx, "employee", "view_email", &r.Email)

	for i := range r.Families {
		maskFamilyResponse(ctx, &r.Families[i])
	}
	for i := range r.Banks {
		maskBankResponse(ctx, &r.Banks[i])
	}
	for i := range r.EmergencyContacts {
		maskEmergencyContactResponse(ctx, &r.EmergencyContacts[i])
	}
}

func maskFamilyResponse(ctx context.Context, r *FamilyResponse) {
	if r == nil {
		return
	}
	maskField(ctx, "employee_family", "view_nik", &r.NIK)
}

func maskBankResponse(ctx context.Context, r *BankResponse) {
	if r == nil {
		return
	}
	maskField(ctx, "employee_bank_account", "view_account_number", &r.AccountNumber)
	maskField(ctx, "employee_bank_account", "view_account_name", &r.AccountName)
}

func maskEmergencyContactResponse(ctx context.Context, r *EmergencyContactResponse) {
	if r == nil {
		return
	}
	maskField(ctx, "emergency_contact", "view_phone_number", &r.PhoneNumber)
}
