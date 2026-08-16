package tenantseed

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// =============================================================================
// Seeder untuk Tenant RBAC (Level 2 — Users, Roles, Permissions)
// =============================================================================
//
// Melengkapi checklist Phase 1 (lihat docs/platform-architecture-design.md):
//   - Auto-seed master data permissions bawaan sistem
//   - Auto-seed default roles tenant (Admin, Employee)
//
// Tabel target (migration 011_settings.sql, pola Spatie):
//   - permissions            (id, name, guard_name, created_at, updated_at)
//   - roles                  (id, name, guard_name, description, is_default, ...)
//   - role_has_permissions   (permission_id, role_id) — composite PK, tanpa kolom id

// tenantRBACSubmenu mendefinisikan satu submenu tenant + action default-nya.
// Nama permission submenu mengikuti format "resource.submenu.action"
// (contoh: "competency.events.view"). Submenu kosong (nil) berarti resource
// hanya punya permission level-module ("resource.action").
type tenantRBACSubmenu struct {
	submenu string
	actions []string
}

// tenantRBACResource mendefinisikan satu resource tenant + action default-nya.
// Nama permission mengikuti format Spatie: "resource.action" (contoh: "organization.view").
type tenantRBACResource struct {
	resource string
	actions  []string
	submenus []tenantRBACSubmenu
}

// tenantRBACResources adalah daftar resource tenant + action yang di-seed.
// Submenus mengikuti struktur menu module (module.go: Menus[].Children) +
// halaman Competency 360 — tiap submenu di-seed dengan action standar
// view/create/update/delete sehingga permission bisa diatur sampai level
// submenu ("resource.submenu.action").
func tenantRBACResources() []tenantRBACResource {
	std := []string{"view", "create", "update", "delete"}
	return []tenantRBACResource{
		{"organization", std, []tenantRBACSubmenu{
			{"tree", std},
			{"zones", std},
			{"job-families", std},
			{"positions", std},
		}},
		{"employee", std, []tenantRBACSubmenu{
			{"list", std},
			{"create", std},
		}},
		{"jobmanagement", std, []tenantRBACSubmenu{
			{"titles", std},
			{"values", std},
			{"objectives", std},
			{"identifications", std},
			{"responsibilities", std},
			{"authorities", std},
			{"working-conditions", std},
			{"scores", std},
			{"competencies", std},
		}},
		{"competency", std, []tenantRBACSubmenu{
			{"competencies", std},
			{"values", std},
			{"indicators", std},
			{"events", std},
			{"scores", std},
			{"templates", std},
			{"raters", std},
			{"my-assessments", std},
			{"manager-assessments", std},
			{"results", std},
			{"reports", std},
		}},
		{"employeemovement", std, []tenantRBACSubmenu{
			{"movements", std},
			{"contracts", std},
		}},
		{"useraccount", std, nil},
		{"attendance", std, []tenantRBACSubmenu{
			{"dashboard", std},
			{"shifts", std},
			{"schedules", std},
			{"events", std},
			{"overtime", std},
			{"locations", std},
			{"business-travel", std},
			{"settings", std},
		}},
		{"approval", std, []tenantRBACSubmenu{
			{"tasks", std},
			{"flows", std},
			{"instances", std},
		}},
		{"payroll", std, []tenantRBACSubmenu{
			{"runs", std},
			{"periods", std},
			{"salary-components", std},
			{"profiles", std},
			{"bpjs-settings", std},
			{"pph21-settings", std},
			{"ptkp-rates", std},
			{"tax-brackets", std},
		}},
		{"leave", std, []tenantRBACSubmenu{
			{"dashboard", std},
			{"requests", std},
			{"types", std},
			{"balances", std},
			{"settings", std},
		}},
		{"performance", std, []tenantRBACSubmenu{
			{"evaluations", std},
			{"templates", std},
			{"indicators", std},
			{"periods", std},
			{"perspectives", std},
		}},
		{"recruitment", std, []tenantRBACSubmenu{
			{"requisitions", std},
			{"candidates", std},
			{"applications", std},
			{"interviews", std},
			{"onboarding", std},
		}},
		{"reimbursement", append(std, "approve"), []tenantRBACSubmenu{
			{"requests", std},
			{"types", std},
			{"reports", std},
		}},
		{"training", []string{"view", "create", "update", "delete", "enroll",
			"course.manage", "session.manage", "participant.manage", "attendance.manage",
			"assessment.manage", "evaluation.manage", "certificate.manage", "plan.manage",
			"request.create", "request.approve", "report.view"}, []tenantRBACSubmenu{
			{"courses", std},
			{"categories", std},
			{"providers", std},
			{"trainers", std},
			{"sessions", std},
			{"participants", std},
			{"planning", std},
			{"requests", std},
			{"needs", std},
			{"certificates", std},
			{"history", std},
			{"reports", std},
		}},
		{"workforceintelligence", std, nil},
		{"careerintelligence", std, nil},
		{"setting", std, []tenantRBACSubmenu{
			{"zones", std},
			{"provinces", std},
			{"regencies", std},
			{"districts", std},
			{"villages", std},
			{"educations", std},
			{"education-majors", std},
			{"religions", std},
			{"marital-statuses", std},
			{"relationship-types", std},
			{"banks", std},
			{"employment-statuses", std},
			{"nationalities", std},
			{"job-families", std},
			{"gradings", std},
			{"salary-grades", std},
			{"ters", std},
			{"ptkps", std},
			{"insurances", std},
			{"company-holidays", std},
			{"competencies", std},
			{"document-templates", std},
		}},
		{"rbac", std, []tenantRBACSubmenu{
			{"roles", std},
		}},
		{"notification", []string{"view", "manage"}, nil},
	}
}

// tenantRBACSubmenuNamesByResource mengembalikan daftar nama submenu per resource
// (hanya resource yang punya submenus). Dipakai untuk men-generate migration
// RBAC submenu dan test deterministik.
func tenantRBACSubmenuNamesByResource() map[string][]string {
	out := map[string][]string{}
	for _, r := range tenantRBACResources() {
		for _, sm := range r.submenus {
			out[r.resource] = append(out[r.resource], sm.submenu)
		}
	}
	return out
}

// SeedTenantRBAC melakukan seeding default roles tenant (Admin, Employee) dan
// permissions bawaan sistem ke tenant database.
//
// Idempotent: semua record menggunakan ID deterministik dari codeToUUID, sehingga
// batchInsert me-skip data yang sudah ada — aman dijalankan berulang kali
// (saat provisioning tenant baru maupun re-seed tenant existing).
func SeedTenantRBAC(db *gorm.DB, l *zap.Logger) error {
	l.Info("Seeding tenant RBAC (roles & permissions)...")

	// ── 1. Seed permissions ──
	var permRecords []map[string]interface{}
	permIDs := make(map[string]string) // "resource.action" -> permission id

	for _, r := range tenantRBACResources() {
		// Permission level-module: "resource.action"
		for _, action := range r.actions {
			name := r.resource + "." + action
			permIDs[name] = codeToUUID("permission", name)
			permRecords = append(permRecords, map[string]interface{}{
				"id":         permIDs[name],
				"name":       name,
				"guard_name": "web",
				"created_at": time.Now(),
				"updated_at": time.Now(),
			})
		}
		// Permission level-submenu: "resource.submenu.action"
		for _, sm := range r.submenus {
			for _, action := range sm.actions {
				name := r.resource + "." + sm.submenu + "." + action
				permIDs[name] = codeToUUID("permission", name)
				permRecords = append(permRecords, map[string]interface{}{
					"id":         permIDs[name],
					"name":       name,
					"guard_name": "web",
					"created_at": time.Now(),
					"updated_at": time.Now(),
				})
			}
		}
	}

	permInserted, permSkipped, err := batchInsert(db, "permissions", permRecords, 50)
	if err != nil {
		return fmt.Errorf("seed tenant permissions failed: %w", err)
	}
	l.Info("  permissions", zap.Int("inserted", permInserted), zap.Int("skipped", permSkipped))

	// ── 2. Seed roles ──
	adminID := codeToUUID("role", "ADMIN")
	employeeID := codeToUUID("role", "EMPLOYEE")

	roleRecords := []map[string]interface{}{
		{
			"id":          adminID,
			"name":        "Admin",
			"guard_name":  "web",
			"description": "Full access to all tenant modules",
			"is_default":  0,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
		{
			"id":          employeeID,
			"name":        "Employee",
			"guard_name":  "web",
			"description": "View-only access to tenant modules",
			"is_default":  1,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		},
	}

	roleInserted, roleSkipped, err := batchInsert(db, "roles", roleRecords, 50)
	if err != nil {
		return fmt.Errorf("seed tenant roles failed: %w", err)
	}
	l.Info("  roles", zap.Int("inserted", roleInserted), zap.Int("skipped", roleSkipped))

	// ── 3. Seed role_has_permissions ──
	// Admin   → semua permission (semua resource, semua action)
	// Employee → hanya action "view"
	linkInserted := 0
	linkSkipped := 0

	for _, r := range tenantRBACResources() {
		for _, action := range r.actions {
			permID := permIDs[r.resource+"."+action]

			if err := upsertRolePermission(db, adminID, permID, &linkInserted, &linkSkipped); err != nil {
				return fmt.Errorf("assign permission %s.%s to Admin failed: %w", r.resource, action, err)
			}

			if action == "view" {
				if err := upsertRolePermission(db, employeeID, permID, &linkInserted, &linkSkipped); err != nil {
					return fmt.Errorf("assign permission %s.%s to Employee failed: %w", r.resource, action, err)
				}
			}
		}
		// Submenu: Admin dapat semua action, Employee hanya "view".
		for _, sm := range r.submenus {
			for _, action := range sm.actions {
				permID := permIDs[r.resource+"."+sm.submenu+"."+action]

				if err := upsertRolePermission(db, adminID, permID, &linkInserted, &linkSkipped); err != nil {
					return fmt.Errorf("assign permission %s.%s.%s to Admin failed: %w", r.resource, sm.submenu, action, err)
				}

				if action == "view" {
					if err := upsertRolePermission(db, employeeID, permID, &linkInserted, &linkSkipped); err != nil {
						return fmt.Errorf("assign permission %s.%s.%s to Employee failed: %w", r.resource, sm.submenu, action, err)
					}
				}
			}
		}
	}
	l.Info("  role_has_permissions", zap.Int("inserted", linkInserted), zap.Int("skipped", linkSkipped))

	return nil
}

// upsertRolePermission menyisipkan satu baris ke role_has_permissions jika belum ada.
// Tabel ini punya composite PK (permission_id, role_id) tanpa kolom id,
// jadi tidak bisa pakai batchInsert() yang berbasis kolom id.
func upsertRolePermission(db *gorm.DB, roleID, permID string, inserted, skipped *int) error {
	var count int64
	if err := db.Table("role_has_permissions").
		Where("permission_id = ? AND role_id = ?", permID, roleID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		*skipped++
		return nil
	}

	if err := db.Table("role_has_permissions").Create(map[string]interface{}{
		"permission_id": permID,
		"role_id":       roleID,
	}).Error; err != nil {
		return err
	}
	*inserted++
	return nil
}
