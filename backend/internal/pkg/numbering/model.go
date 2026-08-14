package numbering

import "time"

const (
	DocumentTypeEmployeeMovement = "employee_movement"
	DocumentTypeEmployeeContract = "employee_contract"
)

const (
	ResetPeriodYearly  = "yearly"
	ResetPeriodMonthly = "monthly"
	ResetPeriodNever   = "never"
)

// DocumentNumberingSetting stores the numbering format and running sequence
// for one document type (employee_movement / employee_contract).
type DocumentNumberingSetting struct {
	ID             string    `gorm:"column:id;primaryKey"`
	DocumentType   string    `gorm:"column:document_type"`
	FormatTemplate string    `gorm:"column:format_template"`
	ResetPeriod    string    `gorm:"column:reset_period"`
	LastSequence   int       `gorm:"column:last_sequence"`
	LastResetKey   string    `gorm:"column:last_reset_key"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (DocumentNumberingSetting) TableName() string {
	return "document_numbering_settings"
}
