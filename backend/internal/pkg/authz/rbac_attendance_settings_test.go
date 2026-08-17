package authz

import "testing"

// =============================================================================
// Regression: permission submenu attendance yang digate strict (admin-config
// absensi: settings.* untuk company settings/shifts/employee-shifts/locations/
// exempt-positions; report.view untuk laporan) harus benar-benar ter-grant ke
// company_admin.
//
// Endpoint absensi yang digate (requireAttendanceSettings /
// requireAttendanceReport, module attendance) dicek via permissions CLAIM di
// JWT. Claim platform user dibangun dari rbac_permissions ×
// rbac_role_permissions (useraccount.FindPlatformPermissionsForRole), jadi
// wildcard pada policy hierarchy Enforcer TIDAK cukup — barisnya harus ada di
// database (pola sama dengan rbac_leave_settings_test.go /
// rbac_sensitive_fields_test.go).
// =============================================================================

// attendanceStrictPerms adalah permission submenu yang harus dimiliki
// super_admin dan company_admin secara default.
var attendanceStrictPerms = []struct{ resource, action string }{
	{"attendance", "settings.view"},
	{"attendance", "settings.create"},
	{"attendance", "settings.update"},
	{"attendance", "settings.delete"},
	{"attendance", "report.view"},
}

func TestSeedDefaults_GrantsAttendanceSettingsPermissionsToCompanyAdmin(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("NewEnforcerFromDB() error = %v", err)
	}

	for _, role := range []string{string(RoleSuperAdmin), string(RoleCompanyAdmin)} {
		claims := loadRoleClaims(t, db, role)
		for _, p := range attendanceStrictPerms {
			want := p.resource + "." + p.action
			if !containsString(claims, want) {
				t.Errorf("role %s is missing permission claim %q", role, want)
			}
		}
	}
}

// Role tenant-level manager/employee TIDAK boleh ikut kebagian — company_admin
// adalah satu-satunya platform role yang mendapat akses penuh tenant; manager/
// employee di platform hanya view/create biasa (module-level attendance.view/create).
func TestSeedDefaults_DoesNotGrantAttendanceSettingsPermissionsToManagerOrEmployee(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("NewEnforcerFromDB() error = %v", err)
	}

	for _, role := range []string{string(RoleManager), string(RoleEmployee)} {
		claims := loadRoleClaims(t, db, role)
		for _, p := range attendanceStrictPerms {
			name := p.resource + "." + p.action
			if containsString(claims, name) {
				t.Errorf("role %s unexpectedly has permission claim %q", role, name)
			}
		}
	}
}

// upsertMissingPermissions dijalankan pada database platform yang SUDAH terisi.
// Tanpa grant ke company_admin di jalur ini, instalasi lama tidak pernah
// menerima permission attendance.settings.* dan company_admin tetap tidak bisa
// mengelola setting absensi.
func TestUpsertMissingPermissions_BackfillsCompanyAdminAttendanceSettings(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	if _, err := NewEnforcerFromDB(db); err != nil {
		t.Fatalf("initial seed error = %v", err)
	}

	// Simulasikan database lama: hapus permission attendance strict beserta grant-nya.
	for _, p := range attendanceStrictPerms {
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
	for _, p := range attendanceStrictPerms {
		want := p.resource + "." + p.action
		if !containsString(claims, want) {
			t.Errorf("after upsert, company_admin is missing permission claim %q", want)
		}
	}
}
