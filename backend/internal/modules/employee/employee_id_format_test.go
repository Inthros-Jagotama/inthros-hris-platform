package employee

import (
	"context"
	"errors"
	"testing"
)

// fakeEmployeeIDFormatProvider adalah implementasi test dari
// EmployeeIDFormatProvider.
type fakeEmployeeIDFormatProvider struct {
	mode      string
	generated string
	genErr    error
	genCalls  int
}

func (f *fakeEmployeeIDFormatProvider) GenerationMode(ctx context.Context) (string, error) {
	return f.mode, nil
}

func (f *fakeEmployeeIDFormatProvider) Generate(ctx context.Context) (string, error) {
	f.genCalls++
	if f.genErr != nil {
		return "", f.genErr
	}
	return f.generated, nil
}

// TestService_Create_NoProvider_FallsBackToManual — tanpa provider (belum
// diwire), perilaku historis dipertahankan: employee_id wajib diisi manual.
func TestService_Create_NoProvider_FallsBackToManual(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	if _, err := svc.Create(ctx, CreateEmployeeRequest{Name: "No Provider"}); !errors.Is(err, ErrEmployeeIDRequired) {
		t.Fatalf("expected ErrEmployeeIDRequired, got %v", err)
	}

	resp, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "MANUAL-001", Name: "No Provider 2"})
	if err != nil {
		t.Fatalf("create with employee_id failed: %v", err)
	}
	if resp.EmployeeID != "MANUAL-001" {
		t.Fatalf("expected employee_id MANUAL-001, got %q", resp.EmployeeID)
	}
}

// TestService_Create_ManualMode_RequiresEmployeeID — mode MANUAL menolak
// request tanpa employee_id, dan memakai nilai yang diberikan bila ada.
func TestService_Create_ManualMode_RequiresEmployeeID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	provider := &fakeEmployeeIDFormatProvider{mode: EmployeeIDGenerationModeManual, generated: "SHOULD-NOT-BE-USED"}
	svc.SetEmployeeIDFormatProvider(provider)
	ctx := context.Background()

	if _, err := svc.Create(ctx, CreateEmployeeRequest{Name: "Manual No ID"}); !errors.Is(err, ErrEmployeeIDRequired) {
		t.Fatalf("expected ErrEmployeeIDRequired, got %v", err)
	}
	if provider.genCalls != 0 {
		t.Fatalf("Generate should not be called in MANUAL mode, got %d calls", provider.genCalls)
	}

	resp, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "MAN-001", Name: "Manual With ID"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if resp.EmployeeID != "MAN-001" {
		t.Fatalf("expected employee_id MAN-001, got %q", resp.EmployeeID)
	}
}

// TestService_Create_AutoMode_IgnoresClientEmployeeID — mode AUTO selalu
// men-generate employee_id, bahkan bila klien mengirim employee_id.
func TestService_Create_AutoMode_IgnoresClientEmployeeID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	provider := &fakeEmployeeIDFormatProvider{mode: EmployeeIDGenerationModeAuto, generated: "AUTO-GEN-001"}
	svc.SetEmployeeIDFormatProvider(provider)
	ctx := context.Background()

	resp, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "CLIENT-SUPPLIED", Name: "Auto Mode"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if resp.EmployeeID != "AUTO-GEN-001" {
		t.Fatalf("expected generated employee_id AUTO-GEN-001, got %q", resp.EmployeeID)
	}
	if provider.genCalls != 1 {
		t.Fatalf("expected Generate called once, got %d", provider.genCalls)
	}

	// Even with no client-supplied employee_id, AUTO still generates.
	provider.generated = "AUTO-GEN-002"
	resp2, err := svc.Create(ctx, CreateEmployeeRequest{Name: "Auto Mode No ID"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if resp2.EmployeeID != "AUTO-GEN-002" {
		t.Fatalf("expected generated employee_id AUTO-GEN-002, got %q", resp2.EmployeeID)
	}
}

// TestService_Create_HybridMode_UsesClientIDIfGiven — mode HYBRID memakai
// employee_id klien jika diisi, atau generate bila kosong.
func TestService_Create_HybridMode_UsesClientIDIfGiven(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	provider := &fakeEmployeeIDFormatProvider{mode: EmployeeIDGenerationModeHybrid, generated: "HYBRID-GEN-001"}
	svc.SetEmployeeIDFormatProvider(provider)
	ctx := context.Background()

	resp, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "HYB-CLIENT", Name: "Hybrid With ID"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if resp.EmployeeID != "HYB-CLIENT" {
		t.Fatalf("expected employee_id HYB-CLIENT, got %q", resp.EmployeeID)
	}
	if provider.genCalls != 0 {
		t.Fatalf("Generate should not be called when client supplies employee_id, got %d calls", provider.genCalls)
	}

	resp2, err := svc.Create(ctx, CreateEmployeeRequest{Name: "Hybrid Without ID"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if resp2.EmployeeID != "HYBRID-GEN-001" {
		t.Fatalf("expected generated employee_id HYBRID-GEN-001, got %q", resp2.EmployeeID)
	}
	if provider.genCalls != 1 {
		t.Fatalf("expected Generate called once, got %d", provider.genCalls)
	}
}
