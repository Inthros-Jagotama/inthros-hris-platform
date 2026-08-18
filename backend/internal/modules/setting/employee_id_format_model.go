package setting

import "time"

const (
	EmployeeIDGenerationModeManual = "MANUAL"
	EmployeeIDGenerationModeHybrid = "HYBRID"
	EmployeeIDGenerationModeAuto   = "AUTO"
)

// EmployeeIDFormatSetting stores the tenant-wide configuration for how
// employee_id values are generated on employee create. Singleton table:
// exactly one row per tenant, auto-created with defaults on first read
// (see Repository.GetOrCreate) — mirrors the attendance_company_settings
// singleton-row convention.
type EmployeeIDFormatSetting struct {
	ID             string    `gorm:"column:id;primaryKey" json:"id"`
	GenerationMode string    `gorm:"column:generation_mode" json:"generation_mode"`
	FormatTemplate string    `gorm:"column:format_template" json:"format_template"`
	ResetPeriod    string    `gorm:"column:reset_period" json:"reset_period"`
	LastSequence   int       `gorm:"column:last_sequence" json:"last_sequence"`
	LastResetKey   string    `gorm:"column:last_reset_key" json:"last_reset_key"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (EmployeeIDFormatSetting) TableName() string {
	return "employee_id_format_settings"
}
