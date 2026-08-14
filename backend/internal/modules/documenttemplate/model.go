package documenttemplate

import (
	"encoding/json"
	"time"
)

const (
	DocumentTypeContractAgreement = "CONTRACT_AGREEMENT"
	DocumentTypeMovementSK        = "MOVEMENT_SK"
)

const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
)

// ValidDocumentTypes is the closed set accepted by this phase; extend when
// spec §2.2's future document types are added.
var ValidDocumentTypes = map[string]bool{
	DocumentTypeContractAgreement: true,
	DocumentTypeMovementSK:        true,
}

type DocumentTemplate struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	Name            string     `gorm:"column:name" json:"name"`
	Code            string     `gorm:"column:code" json:"code"`
	DocumentType    string     `gorm:"column:type" json:"document_type"`
	Description     *string    `gorm:"column:description" json:"description,omitempty"`
	Content         *string    `gorm:"column:content" json:"content,omitempty"`
	ActiveVersionID *string    `gorm:"column:active_version_id" json:"active_version_id,omitempty"`
	Status          string     `gorm:"column:status" json:"status"`
	IsActive        bool       `gorm:"column:is_active" json:"is_active"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (DocumentTemplate) TableName() string { return "document_templates" }

type DocumentTemplateVersion struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	TemplateID   string    `gorm:"column:template_id" json:"template_id"`
	Version      int       `gorm:"column:version" json:"version"`
	Content      string    `gorm:"column:content" json:"content"`
	FileName     *string   `gorm:"column:file_name" json:"file_name,omitempty"`
	FileURL      string    `gorm:"-" json:"file_url,omitempty"`
	PaperSize    string    `gorm:"column:paper_size" json:"paper_size"`
	Orientation  string    `gorm:"column:orientation" json:"orientation"`
	MarginTop    int       `gorm:"column:margin_top" json:"margin_top"`
	MarginRight  int       `gorm:"column:margin_right" json:"margin_right"`
	MarginBottom int       `gorm:"column:margin_bottom" json:"margin_bottom"`
	MarginLeft   int       `gorm:"column:margin_left" json:"margin_left"`
	CreatedBy    *string   `gorm:"column:created_by" json:"created_by,omitempty"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (DocumentTemplateVersion) TableName() string { return "document_template_versions" }

type DocumentTemplateAudit struct {
	ID         string          `gorm:"column:id;primaryKey" json:"id"`
	TemplateID string          `gorm:"column:template_id" json:"template_id"`
	VersionID  *string         `gorm:"column:version_id" json:"version_id,omitempty"`
	Action     string          `gorm:"column:action" json:"action"`
	ActorID    *string         `gorm:"column:actor_id" json:"actor_id,omitempty"`
	Payload    json.RawMessage `gorm:"column:payload" json:"payload,omitempty"`
	CreatedAt  time.Time       `gorm:"column:created_at" json:"created_at"`
}

func (DocumentTemplateAudit) TableName() string { return "document_template_audits" }

type GeneratedDocument struct {
	ID                string    `gorm:"column:id;primaryKey" json:"id"`
	TemplateID        string    `gorm:"column:template_id" json:"template_id"`
	TemplateVersionID string    `gorm:"column:template_version_id" json:"template_version_id"`
	DocumentType      string    `gorm:"column:document_type" json:"document_type"`
	ReferenceType     string    `gorm:"column:reference_type" json:"reference_type"`
	ReferenceID       string    `gorm:"column:reference_id" json:"reference_id"`
	FileName          string    `gorm:"column:file_name" json:"file_name"`
	FilePath          string    `gorm:"column:file_path" json:"file_path"`
	MimeType          string    `gorm:"column:mime_type" json:"mime_type"`
	GeneratedBy       *string   `gorm:"column:generated_by" json:"generated_by,omitempty"`
	GeneratedAt       time.Time `gorm:"column:generated_at" json:"generated_at"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
	FileURL           string    `gorm:"-" json:"file_url,omitempty"`
}

func (GeneratedDocument) TableName() string { return "generated_documents" }
