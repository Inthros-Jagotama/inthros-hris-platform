package setting

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── Zone ──
type Zone struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_zone_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Region    string         `gorm:"type:varchar(100)" json:"region,omitempty"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Zone) TableName() string { return "zones" }
func (z *Zone) BeforeCreate(tx *gorm.DB) error {
	if z.ID == uuid.Nil { z.ID = uuid.New() }
	return nil
}

// ── Province ──
type Province struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(10);not null;uniqueIndex:idx_province_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Province) TableName() string { return "provinces" }
func (p *Province) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil { p.ID = uuid.New() }
	return nil
}

// ── Regency ──
type Regency struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code       string         `gorm:"type:varchar(10);not null;uniqueIndex:idx_regency_code" json:"code"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	ProvinceID uuid.UUID      `gorm:"type:char(36);not null;index:idx_regency_province" json:"province_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Regency) TableName() string { return "regencies" }
func (r *Regency) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil { r.ID = uuid.New() }
	return nil
}

// ── District ──
type District struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(10);not null;uniqueIndex:idx_district_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	RegencyID uuid.UUID      `gorm:"type:char(36);not null;index:idx_district_regency" json:"regency_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (District) TableName() string { return "districts" }
func (d *District) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil { d.ID = uuid.New() }
	return nil
}

// ── Village ──
type Village struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code       string         `gorm:"type:varchar(15);not null;uniqueIndex:idx_village_code" json:"code"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	DistrictID uuid.UUID      `gorm:"type:char(36);not null;index:idx_village_district" json:"district_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Village) TableName() string { return "villages" }
func (v *Village) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil { v.ID = uuid.New() }
	return nil
}

// ── Education ──
type Education struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_education_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Education) TableName() string { return "educations" }
func (e *Education) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil { e.ID = uuid.New() }
	return nil
}

// ── Religion ──
type Religion struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_religion_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Religion) TableName() string { return "religions" }
func (r *Religion) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil { r.ID = uuid.New() }
	return nil
}

// ── MaritalStatus ──
type MaritalStatus struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_marital_status_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (MaritalStatus) TableName() string { return "marital_statuses" }
func (m *MaritalStatus) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil { m.ID = uuid.New() }
	return nil
}

// ── RelationshipType ──
type RelationshipType struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_relationship_type_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (RelationshipType) TableName() string { return "relationship_types" }
func (rt *RelationshipType) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil { rt.ID = uuid.New() }
	return nil
}

// ── EmploymentStatus ──
type EmploymentStatus struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_employment_status_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (EmploymentStatus) TableName() string { return "employment_statuses" }
func (es *EmploymentStatus) BeforeCreate(tx *gorm.DB) error {
	if es.ID == uuid.Nil { es.ID = uuid.New() }
	return nil
}

// ── Bank ──
type Bank struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_bank_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Bank) TableName() string { return "banks" }
func (b *Bank) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil { b.ID = uuid.New() }
	return nil
}

// ── Nationality ──
type Nationality struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_nationality_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Nationality) TableName() string { return "nationalities" }
func (n *Nationality) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil { n.ID = uuid.New() }
	return nil
}

// ── JobFamily ──
type JobFamily struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_job_family_code" json:"code"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (JobFamily) TableName() string { return "job_families" }
func (jf *JobFamily) BeforeCreate(tx *gorm.DB) error {
	if jf.ID == uuid.Nil { jf.ID = uuid.New() }
	return nil
}

// ── SalaryGrade ──
type SalaryGrade struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_salary_grade_code" json:"code"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	MinAmount   float64        `gorm:"type:decimal(18,2);default:0" json:"min_amount,omitempty"`
	MaxAmount   float64        `gorm:"type:decimal(18,2);default:0" json:"max_amount,omitempty"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (SalaryGrade) TableName() string { return "salary_grades" }
func (sg *SalaryGrade) BeforeCreate(tx *gorm.DB) error {
	if sg.ID == uuid.Nil { sg.ID = uuid.New() }
	return nil
}
