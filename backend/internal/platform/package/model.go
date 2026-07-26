package pkgmgr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PackageStatus enum
type PackageStatus string

const (
	PackageDraft     PackageStatus = "draft"
	PackagePublished PackageStatus = "published"
	PackageArchived  PackageStatus = "archived"
)

// Package merepresentasikan bundling modul tenant dengan harga.
type Package struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Price       float64        `gorm:"type:decimal(15,2);default:0" json:"price"`
	Status      string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	IsPublic    bool           `gorm:"default:false" json:"is_public"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Modules     []PackageModule `gorm:"foreignKey:PackageID" json:"modules,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Package) TableName() string {
	return "packages"
}

func (p *Package) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// PackageModule merepresentasikan relasi package-module.
type PackageModule struct {
	PackageID   uuid.UUID `gorm:"type:char(36);primaryKey" json:"package_id"`
	ModuleID    uuid.UUID `gorm:"type:char(36);primaryKey" json:"module_id"`
	IsMandatory bool      `gorm:"default:false" json:"is_mandatory"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	ModuleName string `gorm:"type:varchar(255);not null" json:"module_name,omitempty"`
	ModuleSlug string `gorm:"type:varchar(100);not null" json:"module_slug,omitempty"`
}

func (PackageModule) TableName() string {
	return "package_modules"
}
