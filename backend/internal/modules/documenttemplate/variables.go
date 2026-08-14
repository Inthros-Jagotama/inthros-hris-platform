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
			Variables: []TemplateVariable{
				{Key: "employee.name", Label: "Name"},
				{Key: "employee.nik", Label: "NIK"},
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
