package modulemgmt

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlatformModule merepresentasikan modul yang terdaftar di platform.
type PlatformModule struct {
	ID          string  `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string  `gorm:"type:varchar(100);unique;not null" json:"slug"`
	Version     string  `gorm:"type:varchar(20);not null" json:"version"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	ModuleType  string  `gorm:"type:varchar(20);default:tenant;not null" json:"module_type"`
	// int8, bukan bool: kolom is_core SMALLINT (004_create_modules.sql), pgx
	// menolak encode Go bool ke smallint ("cannot find encode plan").
	IsCore    int8           `gorm:"default:0" json:"is_core"`
	DependsOn *string        `gorm:"type:text" json:"depends_on,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PlatformModule) TableName() string {
	return "modules"
}

func (m *PlatformModule) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// CompanyModule merepresentasikan relasi company-module (module yang diaktifkan untuk company).
type CompanyModule struct {
	CompanyID uuid.UUID `gorm:"type:char(36);primaryKey" json:"company_id"`
	ModuleID  uuid.UUID `gorm:"type:char(36);primaryKey" json:"module_id"`
	// int8, bukan bool: kolom enabled SMALLINT (005_create_company_modules.sql),
	// pgx menolak encode Go bool ke smallint ("cannot find encode plan").
	Enabled     int8       `gorm:"default:1" json:"enabled"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
}

func (CompanyModule) TableName() string {
	return "company_modules"
}

// boolToInt8 mengonversi bool ke int8 untuk kolom smallint (lihat komentar
// IsCore/Enabled di atas soal kenapa bool langsung tidak bisa dipakai).
func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
