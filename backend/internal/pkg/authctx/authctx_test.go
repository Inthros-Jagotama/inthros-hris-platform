package authctx

import (
	"context"
	"testing"
)

func TestGetPermissions(t *testing.T) {
	ctx := context.WithValue(context.Background(), "permissions", []string{"employee.view", "employee.view_nik"})
	got := GetPermissions(ctx)
	if len(got) != 2 || got[0] != "employee.view" {
		t.Fatalf("GetPermissions() = %v, want [employee.view employee.view_nik]", got)
	}
}

func TestGetPermissions_Missing(t *testing.T) {
	got := GetPermissions(context.Background())
	if len(got) != 0 {
		t.Fatalf("GetPermissions() = %v, want empty", got)
	}
}

func TestHasPermission(t *testing.T) {
	ctx := context.WithValue(context.Background(), "permissions", []string{"employee.view_nik"})
	if !HasPermission(ctx, "employee", "view_nik") {
		t.Error("HasPermission(employee, view_nik) = false, want true")
	}
	if HasPermission(ctx, "employee", "view_account_number") {
		t.Error("HasPermission(employee, view_account_number) = true, want false")
	}
}

func TestHasPermission_Wildcard(t *testing.T) {
	ctx := context.WithValue(context.Background(), "permissions", []string{"*"})
	if !HasPermission(ctx, "employee", "view_nik") {
		t.Error("HasPermission with wildcard = false, want true")
	}
}
