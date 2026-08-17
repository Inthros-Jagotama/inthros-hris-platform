package authz

import "testing"

// =============================================================================
// Regression: permission submenu training.settings.* (admin-config training:
// categories, courses, providers, trainers) harus benar-benar ter-grant ke
// company_admin.
//
// Endpoint master training digate strict oleh requireTrainingSettings (module
// training) via permissions CLAIM di JWT. Claim platform user dibangun dari
// rbac_permissions × rbac_role_permissions
// (useraccount.FindPlatformPermissionsForRole), jadi wildcard pada policy
// hierarchy Enforcer TIDAK cukup — barisnya harus ada di database (pola sama
// dengan rbac_attendance_settings_test.go / rbac_leave_settings_test.go).
// =============================================================================

// trainingStrictPerms adalah permission submenu yang harus dimiliki
// super_admin dan company_admin secara default.
var trainingStrictPerms = []struct{ resource, action string }{
	{"training", "settings.view"},
	{"training", "settings.create"},
	{"training", "settings.update"},
	{"training", "settings.delete"},
}

func TestSeedDefaults_GrantsTrainingSettingsPermissionsToCompanyAdmin(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("NewEnforcerFromDB() error = %v", err)
	}

	for _, role := range []string{string(RoleSuperAdmin), string(RoleCompanyAdmin)} {
		claims := loadRoleClaims(t, db, role)
		for _, p := range trainingStrictPerms {
			want := p.resource + "." + p.action
			if !containsString(claims, want) {
				t.Errorf("role %s is missing permission claim %q", role, want)
			}
		}
	}
}

// Role tenant-level manager/employee TIDAK boleh ikut kebagian — company_admin
// adalah satu-satunya platform role yang mendapat akses penuh tenant; manager/
// employee di platform hanya view/create biasa (module-level training.view/create).
func TestSeedDefaults_DoesNotGrantTrainingSettingsPermissionsToManagerOrEmployee(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("NewEnforcerFromDB() error = %v", err)
	}

	for _, role := range []string{string(RoleManager), string(RoleEmployee)} {
		claims := loadRoleClaims(t, db, role)
		for _, p := range trainingStrictPerms {
			name := p.resource + "." + p.action
			if containsString(claims, name) {
				t.Errorf("role %s unexpectedly has permission claim %q", role, name)
			}
		}
	}
}

// upsertMissingPermissions dijalankan pada database platform yang SUDAH terisi.
// Tanpa grant ke company_admin di jalur ini, instalasi lama tidak pernah
// menerima permission training.settings.* dan company_admin tetap tidak bisa
// mengelola master training.
func TestUpsertMissingPermissions_BackfillsCompanyAdminTrainingSettings(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("initial seed error = %v", err)
	}

	// Simulasikan database lama: hapus permission training.settings.* beserta grant-nya.
	for _, p := range trainingStrictPerms {
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
	for _, p := range trainingStrictPerms {
		want := p.resource + "." + p.action
		if !containsString(claims, want) {
			t.Errorf("after upsert, company_admin is missing permission claim %q", want)
		}
	}
}
