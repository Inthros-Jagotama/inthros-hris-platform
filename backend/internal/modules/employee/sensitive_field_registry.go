package employee

// SensitiveFieldDef mendefinisikan satu field sensitif yang bisa
// di-toggle enkripsinya (lewat sensitive_field_settings) dan dibatasi
// akses lihatnya (lewat permission RBAC "Resource.Action").
type SensitiveFieldDef struct {
	// Key adalah field_key di tabel sensitive_field_settings,
	// format "<model>.<field>", contoh "employee.nik".
	Key string
	// Resource dan Action membentuk permission RBAC "Resource.Action"
	// yang mengontrol apakah caller boleh melihat nilai asli (bukan masked).
	Resource string
	Action   string
}

// SensitiveFieldRegistry adalah daftar tetap field sensitif yang didukung
// sistem. Admin hanya bisa toggle enkripsi field yang ada di daftar ini —
// daftar ini tidak bisa diubah lewat UI/API, hanya lewat kode.
var SensitiveFieldRegistry = []SensitiveFieldDef{
	{Key: "employee.nik", Resource: "employee", Action: "view_nik"},
	{Key: "employee.passport", Resource: "employee", Action: "view_passport"},
	{Key: "employee.phone_number", Resource: "employee", Action: "view_phone_number"},
	{Key: "employee.email", Resource: "employee", Action: "view_email"},
	{Key: "employee_family.nik", Resource: "employee_family", Action: "view_nik"},
	{Key: "employee_bank_account.account_number", Resource: "employee_bank_account", Action: "view_account_number"},
	{Key: "employee_bank_account.account_name", Resource: "employee_bank_account", Action: "view_account_name"},
	{Key: "emergency_contact.phone_number", Resource: "emergency_contact", Action: "view_phone_number"},
}

// FieldDef mencari definisi field sensitif berdasarkan field_key.
func FieldDef(key string) (SensitiveFieldDef, bool) {
	for _, d := range SensitiveFieldRegistry {
		if d.Key == key {
			return d, true
		}
	}
	return SensitiveFieldDef{}, false
}
