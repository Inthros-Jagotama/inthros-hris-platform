package rbac

import "testing"

// TestSplitPermissionName memverifikasi pemecahan nama permission menjadi
// resource, submenu, dan action untuk format level-module ("resource.action")
// dan level-submenu ("resource.submenu.action").
func TestSplitPermissionName(t *testing.T) {
	cases := []struct {
		name     string
		resource string
		submenu  string
		action   string
	}{
		// Level-module: "resource.action"
		{"employee.view", "employee", "", "view"},
		{"organization.create", "organization", "", "create"},
		{"rbac.view", "rbac", "", "view"},
		// Level-submenu: "resource.submenu.action"
		{"competency.events.view", "competency", "events", "view"},
		{"competency.templates.update", "competency", "templates", "update"},
		{"attendance.business-travel.view", "attendance", "business-travel", "view"},
		{"setting.zones.delete", "setting", "zones", "delete"},
		// Submenu dengan titik di bagian tengah (khusus training legacy)
		{"training.course.manage", "training", "course", "manage"},
		{"training.request.approve", "training", "request", "approve"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resource, submenu, action := splitPermissionName(c.name)
			if resource != c.resource {
				t.Errorf("resource = %q, want %q", resource, c.resource)
			}
			if submenu != c.submenu {
				t.Errorf("submenu = %q, want %q", submenu, c.submenu)
			}
			if action != c.action {
				t.Errorf("action = %q, want %q", action, c.action)
			}
		})
	}
}

// TestSplitPermissionNameEdgeCases menangani nama tanpa titik atau nama kosong.
func TestSplitPermissionNameEdgeCases(t *testing.T) {
	resource, submenu, action := splitPermissionName("")
	if resource != "" || submenu != "" || action != "" {
		t.Errorf("empty name → (%q,%q,%q), want all empty", resource, submenu, action)
	}

	resource, submenu, action = splitPermissionName("permission")
	if resource != "permission" || submenu != "" || action != "" {
		t.Errorf("single-part name → (%q,%q,%q), want (permission,,)", resource, submenu, action)
	}
}
