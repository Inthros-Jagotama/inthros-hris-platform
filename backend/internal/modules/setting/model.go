package setting

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── Zone ──
type Zone struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_zone_code" json:"code"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Zone        string         `gorm:"column:zone;type:varchar(200);not null" json:"-"`
	Region      string         `gorm:"type:varchar(100)" json:"region,omitempty"`
	Description string         `gorm:"column:description;type:varchar(255);not null" json:"description,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Zone) TableName() string { return "zones" }
func (z *Zone) BeforeCreate(tx *gorm.DB) error {
	if z.ID == uuid.Nil { z.ID = uuid.New() }
	if z.Zone == "" { z.Zone = z.Name }

	return nil
}

// ── Province (ID from BPS code, char(2)) ──
type Province struct {
	ID        string         `gorm:"type:char(2);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(10);not null;uniqueIndex:idx_province_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Province) TableName() string { return "provinces" }

// ── Regency (ID from BPS code, char(4)) ──
type Regency struct {
	ID         string         `gorm:"type:char(4);primaryKey" json:"id"`
	Code       string         `gorm:"type:varchar(10);not null;uniqueIndex:idx_regency_code" json:"code"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	ProvinceID string         `gorm:"type:char(2);not null;index:idx_regency_province" json:"province_id"`
	Province   *Province      `gorm:"foreignKey:ProvinceID;references:ID" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Regency) TableName() string { return "regencies" }

// ── District (ID from BPS code, char(6)) ──
type District struct {
	ID        string         `gorm:"type:char(6);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(15);not null;uniqueIndex:idx_district_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	RegencyID string         `gorm:"type:char(4);not null;index:idx_district_regency" json:"regency_id"`
	Regency   *Regency       `gorm:"foreignKey:RegencyID;references:ID" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (District) TableName() string { return "districts" }

// ── Village (ID from BPS code, char(10)) ──
type Village struct {
	ID         string         `gorm:"type:char(10);primaryKey" json:"id"`
	Code       string         `gorm:"type:varchar(15);not null;uniqueIndex:idx_village_code" json:"code"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	DistrictID string         `gorm:"type:char(6);not null;index:idx_village_district" json:"district_id"`
	District   *District      `gorm:"foreignKey:DistrictID;references:ID" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Village) TableName() string { return "villages" }

// ── Education ──
type Education struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_education_code" json:"code"`
	Name      string         `gorm:"type:varchar(200);not null" json:"name"`
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

// ── EducationMajor (jurusan pendidikan) ──
type EducationMajor struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_education_major_code" json:"code"`
	Name      string         `gorm:"type:varchar(200);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (EducationMajor) TableName() string { return "education_majors" }
func (em *EducationMajor) BeforeCreate(tx *gorm.DB) error {
	if em.ID == uuid.Nil { em.ID = uuid.New() }
	return nil
}

// ── Religion ──
type Religion struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_religion_code" json:"code"`
	Name      string         `gorm:"type:varchar(200);not null" json:"name"`
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
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
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

// ── Grading ──
type Grading struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_grading_code" json:"code"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Grading) TableName() string { return "gradings" }
func (g *Grading) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil { g.ID = uuid.New() }
	return nil
}

// ── Insurance ──
type Insurance struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_insurance_code" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Insurance) TableName() string { return "insurances" }
func (ins *Insurance) BeforeCreate(tx *gorm.DB) error {
	if ins.ID == uuid.Nil { ins.ID = uuid.New() }
	return nil
}

// ── CompanyHoliday (Hari Libur Perusahaan) ──
// Tabel company_holidays hanya memiliki kolom: id, holiday_date, name,
// description, is_active, created_at (tanpa updated_at / deleted_at).
// Catatan: IsActive sengaja TANPA tag `default:true` (berbeda dari Zone) karena
// GORM akan mengubah false → true untuk field zero-value bertag default.
type CompanyHoliday struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	HolidayDate string    `gorm:"type:date;not null;uniqueIndex:uk_company_holiday" json:"holiday_date"`
	Name        string    `gorm:"type:varchar(200);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

func (CompanyHoliday) TableName() string { return "company_holidays" }
func (ch *CompanyHoliday) BeforeCreate(tx *gorm.DB) error {
	if ch.ID == uuid.Nil { ch.ID = uuid.New() }
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

// ── TER (Tarif Efektif Rata-rata) ──
type TER struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Group     string         `gorm:"column:group;type:char(1);not null;index:idx_ter_group" json:"group"`
	BrutoMin  *int64         `gorm:"column:bruto_min" json:"bruto_min,omitempty"`
	BrutoMax  *int64         `gorm:"column:bruto_max" json:"bruto_max,omitempty"`
	Rate      float64        `gorm:"type:decimal(10,2);not null" json:"rate"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (TER) TableName() string { return "ters" }
func (t *TER) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil { t.ID = uuid.New() }
	return nil
}

// ── Competency (kompetensi) ──
// Tabel competencies (berbagi dengan module competency): id, name, field,
// cluster, definition, created_by, updated_by, timestamps.
// TIDAK ada kolom code / sort_order / deleted_at → validasi duplikat pakai name,
// dan Delete adalah hard delete.
type Competency struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name       string     `gorm:"type:varchar(255);not null;index" json:"name"`
	Field      *string    `gorm:"type:varchar(255)" json:"field,omitempty"`
	Cluster    *string    `gorm:"type:varchar(255)" json:"cluster,omitempty"`
	Definition *string    `gorm:"type:text" json:"definition,omitempty"`
	CreatedBy  *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy  *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (Competency) TableName() string { return "competencies" }
func (c *Competency) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil { c.ID = uuid.New() }
	return nil
}

// ── PTKP (Penghasilan Tidak Kena Pajak) ──
type PTKP struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Code      string         `gorm:"column:code;type:varchar(20);not null;default:'';index:idx_ptkp_code" json:"code"`
	PTKP      int64          `gorm:"column:ptkp;not null" json:"ptkp"`
	Group     string         `gorm:"column:group;type:char(1);not null;index:idx_ptkp_group" json:"group"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (PTKP) TableName() string { return "ptkps" }
func (p *PTKP) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil { p.ID = uuid.New() }
	return nil
}
