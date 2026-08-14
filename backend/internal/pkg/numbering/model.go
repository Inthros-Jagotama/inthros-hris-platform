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
	ID             string    `gorm:"column:id;primaryKey" json:"id"`
	DocumentType   string    `gorm:"column:document_type" json:"document_type"`
	FormatTemplate string    `gorm:"column:format_template" json:"format_template"`
	ResetPeriod    string    `gorm:"column:reset_period" json:"reset_period"`
	LastSequence   int       `gorm:"column:last_sequence" json:"last_sequence"`
	LastResetKey   string    `gorm:"column:last_reset_key" json:"last_reset_key"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (DocumentNumberingSetting) TableName() string {
	return "document_numbering_settings"
}
