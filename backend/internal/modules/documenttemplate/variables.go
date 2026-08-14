package documenttemplate

type TemplateVariable struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type VariableGroup struct {
	Category  string             `json:"category"`
	Variables []TemplateVariable `json:"variables"`
}

// VariableRegistry is the static list of placeholders available to the
// template editor's "Insert Variable" picker (spec §7). Backend-owned so the
// frontend never has to keep this list in sync by hand.
func VariableRegistry() []VariableGroup {
	return []VariableGroup{
		{
			Category: "employee",
			// Field mengikuti kolom tabel `employees` (+ relasi religion &
			// marital_status, join_date dari employment pertama) — resolusi
			// dilakukan di employeemovement.GetEmployeeProfile.
			Variables: []TemplateVariable{
				{Key: "employee.employee_id", Label: "Employee Number"},
				{Key: "employee.name", Label: "Name"},
				{Key: "employee.nik", Label: "NIK"},
				{Key: "employee.family_id", Label: "Family Card Number"},
				{Key: "employee.mother_name", Label: "Mother's Name"},
				{Key: "employee.gender", Label: "Gender"},
				{Key: "employee.dob", Label: "Date of Birth"},
				{Key: "employee.pob", Label: "Place of Birth"},
				{Key: "employee.nationality_type", Label: "Nationality Type"},
				{Key: "employee.nationality_id", Label: "Nationality Code"},
				{Key: "employee.passport", Label: "Passport Number"},
				{Key: "employee.phone_number", Label: "Phone Number"},
				{Key: "employee.email", Label: "Email"},
				{Key: "employee.linkedin", Label: "LinkedIn"},
				{Key: "employee.instagram", Label: "Instagram"},
				{Key: "employee.religion", Label: "Religion"},
				{Key: "employee.marital_status", Label: "Marital Status"},
				{Key: "employee.status", Label: "Employee Status"},
				{Key: "employee.join_date", Label: "Join Date"},
				{Key: "employee.position", Label: "Position"},
				{Key: "employee.organization", Label: "Organization"},
			},
		},
		{
			Category: "contract",
			Variables: []TemplateVariable{
				{Key: "contract.number", Label: "Contract Number"},
				{Key: "contract.start_date", Label: "Start Date"},
				{Key: "contract.end_date", Label: "End Date"},
			},
		},
		{
			Category: "movement",
			Variables: []TemplateVariable{
				{Key: "movement.number", Label: "Movement Number"},
				{Key: "movement.effective_date", Label: "Effective Date"},
				{Key: "movement.previous_position", Label: "Previous Position"},
				{Key: "movement.new_position", Label: "New Position"},
			},
		},
		{
			Category: "company",
			Variables: []TemplateVariable{
				{Key: "company.name", Label: "Name"},
				{Key: "company.address", Label: "Address"},
			},
		},
	}
}
