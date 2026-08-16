package employee

import "testing"

func TestSensitiveFieldRegistry_HasEightEntries(t *testing.T) {
	if len(SensitiveFieldRegistry) != 8 {
		t.Fatalf("len(SensitiveFieldRegistry) = %d, want 8", len(SensitiveFieldRegistry))
	}
}

func TestFieldDef_Found(t *testing.T) {
	def, ok := FieldDef("employee.nik")
	if !ok {
		t.Fatal("FieldDef(employee.nik) not found")
	}
	if def.Resource != "employee" || def.Action != "view_nik" {
		t.Errorf("FieldDef(employee.nik) = %+v, want Resource=employee Action=view_nik", def)
	}
}

func TestFieldDef_NotFound(t *testing.T) {
	_, ok := FieldDef("employee.does_not_exist")
	if ok {
		t.Error("FieldDef(employee.does_not_exist) should not be found")
	}
}
