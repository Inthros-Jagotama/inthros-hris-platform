package authz

import (
	"testing"

	"gorm.io/gorm"
)

// =============================================================================
// Regression: permission view_* per-field (Sensitive Data Masking, spec §3)
// harus benar-benar ter-grant ke company_admin.
//
// Masking dievaluasi lewat permissions CLAIM di JWT (authctx.HasPermission).
// Claim platform user dibangun dari rbac_permissions × rbac_role_permissions
// (useraccount.FindPlatformPermissionsForRole), jadi wildcard pada policy
// hierarchy Enforcer TIDAK cukup — barisnya harus ada di database.
// =============================================================================

// sensitiveViewPerms adalah permission per-field yang di-spec §3 harus dimiliki
// super_admin dan company_admin secara default.
var sensitiveViewPerms = []struct{ resource, action string }{
	{"employee", "view_nik"},
	{"employee", "view_passport"},
	{"employee", "view_phone_number"},
	{"employee", "view_email"},
	{"employee_family", "view_nik"},
	{"employee_bank_account", "view_account_number"},
	{"employee_bank_account", "view_account_name"},
	{"emergency_contact", "view_phone_number"},
}

// loadRoleClaims meniru useraccount.FindPlatformPermissionsForRole: daftar
// "resource.action" hasil join role → role_permissions → permissions, yaitu
// persis isi permissions claim di JWT platform user.
func loadRoleClaims(t *testing.T, db *gorm.DB, roleSlug string) []string {
	t.Helper()
	var rows []struct {
		Resource string
		Action   string
	}
	if err := db.Table("rbac_permissions p").
		Select("p.resource, p.action").
		Joins("JOIN rbac_role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN rbac_roles r ON r.id = rp.role_id").
		Where("r.slug = ?", roleSlug).
		Scan(&rows).Error; err != nil {
		t.Fatalf("load claims for %s: %v", roleSlug, err)
	}
	claims := make([]string, 0, len(rows))
	for _, row := range rows {
		claims = append(claims, row.Resource+"."+row.Action)
	}
	return claims
}

func containsString(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}

func TestSeedDefaults_GrantsSensitiveViewPermissionsToCompanyAdmin(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("NewEnforcerFromDB() error = %v", err)
	}

	for _, role := range []string{string(RoleSuperAdmin), string(RoleCompanyAdmin)} {
		claims := loadRoleClaims(t, db, role)
		for _, p := range sensitiveViewPerms {
			want := p.resource + "." + p.action
			if !containsString(claims, want) {
				t.Errorf("role %s is missing permission claim %q", role, want)
			}
		}
	}
}

// Role tenant-level manager/employee TIDAK boleh ikut kebagian (spec §3 hanya
// menyebut super_admin & company_admin).
func TestSeedDefaults_DoesNotGrantSensitiveViewPermissionsToManagerOrEmployee(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("NewEnforcerFromDB() error = %v", err)
	}

	for _, role := range []string{string(RoleManager), string(RoleEmployee)} {
		claims := loadRoleClaims(t, db, role)
		for _, p := range sensitiveViewPerms {
			name := p.resource + "." + p.action
			if containsString(claims, name) {
				t.Errorf("role %s unexpectedly has permission claim %q", role, name)
			}
		}
	}
}

// upsertMissingPermissions dijalankan pada database platform yang SUDAH terisi.
// Tanpa grant ke company_admin di jalur ini, instalasi lama tidak pernah
// menerima permission view_* dan company_admin tetap melihat data ter-mask.
func TestUpsertMissingPermissions_BackfillsCompanyAdmin(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("initial seed error = %v", err)
	}

	// Simulasikan database lama: hapus permission view_* beserta grant-nya.
	for _, p := range sensitiveViewPerms {
		var perm RbacPermission
		if err := db.Where("resource = ? AND action = ?", p.resource, p.action).First(&perm).Error; err != nil {
			t.Fatalf("permission %s.%s not seeded: %v", p.resource, p.action, err)
		}
		if err := db.Where("permission_id = ?", perm.ID).Delete(&RbacRolePermission{}).Error; err != nil {
			t.Fatalf("delete role_permissions: %v", err)
		}
		if err := db.Delete(&perm).Error; err != nil {
			t.Fatalf("delete permission: %v", err)
		}
	}

	// Startup berikutnya: tabel tidak kosong → jalur upsertMissingPermissions.
	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("re-seed error = %v", err)
	}

	claims := loadRoleClaims(t, db, string(RoleCompanyAdmin))
	for _, p := range sensitiveViewPerms {
		want := p.resource + "." + p.action
		if !containsString(claims, want) {
			t.Errorf("after upsert, company_admin is missing permission claim %q", want)
		}
	}
}

// Regression untuk bug akumulasi action di loadFromDB: pengecekan duplikat
// berbasis substring membuat "view" tertelan oleh "view_nik", sehingga
// Check(company_admin, employee, view) berubah jadi deny.
func TestContainsAction_ExactTokenMatch(t *testing.T) {
	if containsAction("view_nik,view_email", "view") {
		t.Error("containsAction must not match a substring of another action")
	}
	if !containsAction("view_nik,view", "view") {
		t.Error("containsAction must match an exact token")
	}
	if !containsAction("view", "view") {
		t.Error("containsAction must match a single-token list")
	}
	if containsAction("view", "delete") {
		t.Error("containsAction must not match an absent action")
	}
}
