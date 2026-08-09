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

// tenantRBACResource mendefinisikan satu resource tenant + action default-nya.
// Nama permission mengikuti format Spatie: "resource.action" (contoh: "organization.view").
type tenantRBACResource struct {
	resource string
	actions  []string
}

// tenantRBACResources adalah daftar resource tenant + action yang di-seed.
func tenantRBACResources() []tenantRBACResource {
	return []tenantRBACResource{
		{"organization", []string{"view", "create", "update", "delete"}},
		{"employee", []string{"view", "create", "update", "delete"}},
		{"jobmanagement", []string{"view", "create", "update", "delete"}},
		{"competency", []string{"view", "create", "update", "delete"}},
		{"employeemovement", []string{"view", "create", "update", "delete"}},
		{"useraccount", []string{"view", "create", "update", "delete"}},
		{"attendance", []string{"view", "create", "update", "delete"}},
		{"approval", []string{"view", "create", "update", "delete"}},
		{"payroll", []string{"view", "create", "update", "delete"}},
		{"leave", []string{"view", "create", "update", "delete"}},
		{"performance", []string{"view", "create", "update", "delete"}},
		{"recruitment", []string{"view", "create", "update", "delete"}},
		{"reimbursement", []string{"view", "create", "update", "delete"}},
		{"training", []string{"view", "create", "update", "delete"}},
		{"workforceintelligence", []string{"view", "create", "update", "delete"}},
		{"careerintelligence", []string{"view", "create", "update", "delete"}},
		{"setting", []string{"view", "create", "update", "delete"}},
		{"rbac", []string{"view", "create", "update", "delete"}},
		{"notification", []string{"view", "manage"}},
	}
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
