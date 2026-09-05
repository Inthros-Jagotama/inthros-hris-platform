package modulemgmt

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlatformModule merepresentasikan modul yang terdaftar di platform.
type PlatformModule struct {
	ID          string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(100);unique;not null" json:"slug"`
	Version     string         `gorm:"type:varchar(20);not null" json:"version"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	ModuleType  string         `gorm:"type:varchar(20);default:tenant;not null" json:"module_type"`
	IsCore      bool           `gorm:"default:false" json:"is_core"`
	DependsOn   *string        `gorm:"type:text" json:"depends_on,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
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
	CompanyID   uuid.UUID  `gorm:"type:char(36);primaryKey" json:"company_id"`
	ModuleID    uuid.UUID  `gorm:"type:char(36);primaryKey" json:"module_id"`
	Enabled     bool       `gorm:"default:true" json:"enabled"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
}

func (CompanyModule) TableName() string {
	return "company_modules"
}
