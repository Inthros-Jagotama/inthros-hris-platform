package documenttemplate

import "testing"

func TestVariableRegistryHasExpectedCategories(t *testing.T) {
	reg := VariableRegistry()
	wantCategories := []string{"employee", "contract", "movement", "company"}
	for _, cat := range wantCategories {
		found := false
		for _, group := range reg {
			if group.Category == cat {
				found = true
				if len(group.Variables) == 0 {
					t.Errorf("category %q has no variables", cat)
				}
				break
			}
		}
		if !found {
			t.Errorf("expected category %q in registry, not found", cat)
		}
	}
}

func TestVariableRegistryKeysAreDotted(t *testing.T) {
	reg := VariableRegistry()
	for _, group := range reg {
		for _, v := range group.Variables {
			if v.Key == "" || v.Label == "" {
				t.Errorf("variable in category %q has empty key or label: %+v", group.Category, v)
			}
		}
	}
}
