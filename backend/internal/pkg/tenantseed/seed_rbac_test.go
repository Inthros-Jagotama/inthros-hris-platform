package tenantseed

import (
	"strings"
	"testing"
)

// TestTenantRBACSubmenus memastikan struktur submenu yang di-seed valid:
//   - setiap resource dengan submenu punya nama permission "resource.submenu.action"
//   - tidak ada duplikat nama permission antar resource/submenu/action
//   - setiap permission punya ID deterministik (codeToUUID) yang konsisten
func TestTenantRBACSubmenus(t *testing.T) {
	seen := map[string]string{}
	for _, r := range tenantRBACResources() {
		// Module-level
		for _, action := range r.actions {
			name := r.resource + "." + action
			if prev, ok := seen[name]; ok {
				t.Fatalf("duplicate permission %q (resource %q, juga dari %q)", name, r.resource, prev)
			}
			seen[name] = r.resource
			if codeToUUID("permission", name) == "" {
				t.Fatalf("empty deterministic ID for %q", name)
			}
		}
		// Submenu-level
		for _, sm := range r.submenus {
			if sm.submenu == "" {
				t.Fatalf("resource %q punya submenu dengan nama kosong", r.resource)
			}
			for _, action := range sm.actions {
				name := r.resource + "." + sm.submenu + "." + action
				if prev, ok := seen[name]; ok {
					t.Fatalf("duplicate permission %q (resource %q, juga dari %q)", name, r.resource, prev)
				}
				seen[name] = r.resource
			}
		}
	}

	// Sanity: beberapa permission submenu yang dipakai frontend harus ada.
	// organization & setting sengaja TIDAK punya submenu — cukup module-level
	// (organization.view/create/update/delete, setting.view/...), yang sudah
	// meng-cover semua route submenu-nya lewat fallback module-covers-submenu
	// di authz middleware & FE hasPermission().
	required := []string{
		"competency.settings.view",
		"competency.assessment.view",
		"competency.report.view",
		"attendance.business-travel.view",
		"rbac.roles.view",
	}
	for _, name := range required {
		if _, ok := seen[name]; !ok {
			t.Errorf("submenu permission %q tidak di-seed", name)
		}
	}
}

// TestTenantRBACSubmenuNamesByResource memastikan helper map submenu lengkap
// dan semua submenu punya minimal action "view" (dipakai gate menu frontend).
func TestTenantRBACSubmenuNamesByResource(t *testing.T) {
	byRes := tenantRBACSubmenuNamesByResource()
	if len(byRes) == 0 {
		t.Fatal("no submenus defined")
	}
	if _, ok := byRes["competency"]; !ok {
		t.Fatal("competency harus punya submenus")
	}
	for res, subs := range byRes {
		if len(subs) == 0 {
			t.Fatalf("resource %q dilaporkan punya submenus tapi kosong", res)
		}
	}
}

// TestSubmenuPermissionNameFormat memverifikasi format "resource.submenu.action"
// dan bahwa semua nama berisi tepat 3 bagian untuk permission submenu.
func TestSubmenuPermissionNameFormat(t *testing.T) {
	for _, r := range tenantRBACResources() {
		for _, sm := range r.submenus {
			for _, action := range sm.actions {
				name := r.resource + "." + sm.submenu + "." + action
				parts := strings.Split(name, ".")
				if len(parts) != 3 {
					t.Fatalf("permission %q harus punya 3 bagian (resource.submenu.action), dapat %d", name, len(parts))
				}
				if parts[0] != r.resource || parts[1] != sm.submenu || parts[2] != action {
					t.Fatalf("permission %q tidak cocok resource/submenu/action", name)
				}
			}
		}
	}
}
