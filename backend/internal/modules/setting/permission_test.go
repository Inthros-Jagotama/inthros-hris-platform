package setting

import (
	"slices"
	"testing"
)

// TestModulePermissionsDeclareDocumentTemplate verifies the RBAC UI permission
// declarations for the document template feature are present in the setting
// module Info(). Runtime enforcement is handled by the authz middleware
// (resource derived from /settings/... path), but the granular declarations
// must exist for the RBAC UI to render/assign them.
func TestModulePermissionsDeclareDocumentTemplate(t *testing.T) {
	perms := (&settingModule{}).Info().Permissions

	expected := []string{
		"setting.document_template.view",
		"setting.document_template.create",
		"setting.document_template.update",
		"setting.document_template.delete",
		"setting.document_template.activate",
		"setting.document_template.deactivate",
		"setting.document_template.version",
	}

	for _, want := range expected {
		if !slices.Contains(perms, want) {
			t.Errorf("missing permission %q in setting module Info().Permissions", want)
		}
	}
}

// TestModuleMenusIncludeDocumentTemplates ensures the Settings menu declares
// the Document Templates entry so it shows up in navigation.
func TestModuleMenusIncludeDocumentTemplates(t *testing.T) {
	menus := (&settingModule{}).Info().Menus

	for _, menu := range menus {
		if menu.Name != "Settings" {
			continue
		}
		for _, child := range menu.Children {
			if child.Name == "Document Templates" {
				if child.Route != "/admin/settings/document-templates" {
					t.Errorf("unexpected Document Templates route: %s", child.Route)
				}
				return
			}
		}
	}
	t.Error("Document Templates menu not found under Settings menu")
}
