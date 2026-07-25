package organization

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization merepresentasikan unit organisasi dalam struktur tree.
// Setiap organization bisa memiliki parent (parent_id) untuk membentuk
// hierarki organisasi.
type Organization struct {
	ID                   uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationSummaryID *uuid.UUID     `gorm:"type:char(36)" json:"organization_summary_id,omitempty"`
	Code                 string          `gorm:"type:varchar(10);not null" json:"code"`
	FullCode             string          `gorm:"type:varchar(50);not null;index" json:"full_code"`
	Nomenclature         string          `gorm:"type:varchar(255);not null" json:"nomenclature"`
	ParentID             *uuid.UUID      `gorm:"type:char(36)" json:"parent_id,omitempty"`
	ZoneID               *uuid.UUID      `gorm:"type:char(36)" json:"zone_id,omitempty"`
	JobFamilyID          *uuid.UUID      `gorm:"type:char(36)" json:"job_family_id,omitempty"`
	GradingID            *uuid.UUID      `gorm:"type:char(36)" json:"grading_id,omitempty"`
	Level                int             `gorm:"default:0" json:"level"`
	SortOrder            int             `gorm:"default:0" json:"sort_order"`
	CreatedBy            *uuid.UUID      `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID      `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	DeletedAt            gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	Parent     *Organization  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children   []Organization `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OrganizationHistory (Audit log perubahan organisasi)
// =========================================================================

// ActionType mendefinisikan tipe aksi pada history.
type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

// OrganizationHistory mencatat setiap perubahan yang terjadi pada entitas Organization.
// OldValues dan NewValues disimpan sebagai JSON string untuk fleksibilitas.
type OrganizationHistory struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID uuid.UUID  `gorm:"type:char(36);not null;index:idx_org_hitory_org" json:"organization_id"`
	Action         ActionType `gorm:"type:varchar(20);not null;index:idx_org_history_action" json:"action"`
	OldValues      *string    `gorm:"type:json" json:"old_values,omitempty"`
	NewValues      *string    `gorm:"type:json" json:"new_values,omitempty"`
	ChangedBy      *uuid.UUID `gorm:"type:char(36)" json:"changed_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (OrganizationHistory) TableName() string {
	return "organization_histories"
}

func (h *OrganizationHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OrganizationVersion (Snapshot struktur organisasi)
// =========================================================================

// VersionStatus mendefinisikan status versi.
type VersionStatus string

const (
	VersionStatusDraft   VersionStatus = "DRAFT"
	VersionStatusActive  VersionStatus = "ACTIVE"
	VersionStatusArchived VersionStatus = "ARCHIVED"
)

// OrganizationVersion menyimpan snapshot lengkap pohon organisasi pada suatu waktu.
// Snapshot berisi JSON array dari seluruh node organisasi dengan parent-child relationship.
type OrganizationVersion struct {
	ID              uuid.UUID    `gorm:"type:char(36);primaryKey" json:"id"`
	VersionName     string       `gorm:"type:varchar(100);not null" json:"version_name"`
	Description     *string      `gorm:"type:text" json:"description,omitempty"`
	Snapshot        string       `gorm:"type:longtext;not null" json:"snapshot"`
	Status          VersionStatus `gorm:"type:varchar(20);default:ACTIVE" json:"status"`
	CreatedBy       *uuid.UUID   `gorm:"type:char(36)" json:"created_by,omitempty"`
	ParentVersionID *uuid.UUID   `gorm:"type:char(36)" json:"parent_version_id,omitempty"`
	NodeCount       int          `gorm:"default:0" json:"node_count"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func (OrganizationVersion) TableName() string {
	return "organization_versions"
}

func (v *OrganizationVersion) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}